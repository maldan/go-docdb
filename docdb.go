package docdb

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/maldan/go-cmhp"
)

type Config struct {
	Mode       string
	WriteDelay int64
	CacheDelay int64

	SpacesKey    string
	SpacesSecret string
	Endpoint     string
	Bucket       string
}

type WriteInfo struct {
	Path        string
	Data        []byte
	IsCompress  bool
	Visibility  string
	ContentType string
}

type CacheInfo struct {
	Path    string
	Data    []byte
	Created time.Time
}

var config Config
var writeQueue sync.Map
var cache sync.Map
var s3Client *s3.S3

func Start(args Config) {
	config = args

	// Init S3
	if config.Mode == "s3" {
		newSession, err := session.NewSession(&aws.Config{
			Credentials: credentials.NewStaticCredentials(config.SpacesKey, config.SpacesSecret, ""),
			Endpoint:    aws.String(config.Endpoint),
			Region:      aws.String("us-east-1"),
		})
		if err != nil {
			log.Fatal(err)
		}
		s3Client = s3.New(newSession)
	}

	go writeSchedule()
	go cacheSchedule()
}

func writeSchedule() {
	log.Println("DOCDB START SCHEDULER")
	log.Printf("Writing Delay %v", config.WriteDelay)

	for {
		writeQueue.Range(func(key, value interface{}) bool {
			info := value.(WriteInfo)
			data := info.Data

			if info.IsCompress {
				data, _ = cmhp.DataCompress(data)
			}

			if config.Mode == "os" {
				ioutil.WriteFile(key.(string), data, 0777)
			} else {
				object := s3.PutObjectInput{
					Bucket:      aws.String(config.Bucket),
					Key:         aws.String(key.(string)),
					Body:        bytes.NewReader(data),
					ACL:         aws.String(info.Visibility),
					ContentType: aws.String(info.ContentType),
				}
				_, err := s3Client.PutObject(&object)
				if err != nil {
					log.Println(err.Error())
					return true
				}
			}
			log.Printf("Saved [%v] -> %v\n", config.Mode, key)

			writeQueue.Delete(key)
			return true
		})

		time.Sleep(time.Second * time.Duration(config.WriteDelay))
	}
}

func cacheSchedule() {
	log.Println("DOCDB Start Cache Scheduler")
	log.Printf("Cache Delay %v", config.CacheDelay)

	for {
		cache.Range(func(key, value interface{}) bool {
			info := value.(CacheInfo)
			if time.Since(info.Created).Seconds() > float64(config.CacheDelay) {
				log.Printf("Remove from cache %v\n", key)
				cache.Delete(key)
			}

			return true
		})

		time.Sleep(time.Second * 1)
	}
}

func SaveDocument(path string, in interface{}) error {
	if config.Mode == "os" {
		os.MkdirAll(filepath.Dir(path), 0777)
	}

	// Stringify doc
	data, err := json.Marshal(in)
	if err != nil {
		return err
	}

	// Put to writing
	writeQueue.Store(path, WriteInfo{
		Path:       path,
		Data:       data,
		IsCompress: true,
	})

	// Save to cache
	cache.Store(path, CacheInfo{
		Path:    path,
		Data:    data,
		Created: time.Now(),
	})
	return nil
}

func LoadDocument(path string, out interface{}) error {
	v, ok := cache.Load(path)
	if ok {
		json.Unmarshal(v.(CacheInfo).Data, out)
	} else {
		if config.Mode == "os" {
			dataFile, err := os.OpenFile(path, os.O_RDONLY, 0777)
			if err != nil {
				return err
			}
			defer dataFile.Close()
			data, err := ioutil.ReadAll(dataFile)
			if err != nil {
				return err
			}

			data, err = cmhp.DataDecompress(data)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data, out)
			if err != nil {
				return err
			}
		} else {
			// Download file
			result, err := s3Client.GetObject(&s3.GetObjectInput{
				Bucket: aws.String(config.Bucket),
				Key:    aws.String(path),
			})
			if err != nil {
				return err
			}

			// Read and decompress
			data, err := io.ReadAll(result.Body)
			if err != nil {
				return err
			}

			data2, err := cmhp.DataDecompress(data)
			if err != nil {
				return err
			}

			err = json.Unmarshal(data2, out)
			if err != nil {
				return err
			}

			// Everything is ok
			// Save to cache
			cache.Store(path, CacheInfo{
				Path:    path,
				Data:    data2,
				Created: time.Now(),
			})
		}
	}

	return nil
}

func SaveFile(path string, data []byte, visibility string, contentType string) {
	// Put to writing
	writeQueue.Store(path, WriteInfo{
		Path:        path,
		Data:        data,
		ContentType: contentType,
		Visibility:  visibility,
	})
}

func LoadFile(path string) ([]byte, error) {
	if config.Mode == "os" {
		dataFile, err := os.OpenFile(path, os.O_RDONLY, 0777)
		if err != nil {
			return nil, err
		}
		defer dataFile.Close()

		data, err := ioutil.ReadAll(dataFile)
		if err != nil {
			return nil, err
		}
		return data, nil
	} else {
		// Download file
		result, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(config.Bucket),
			Key:    aws.String(path),
		})
		if err != nil {
			return nil, err
		}

		// Read
		data, err := io.ReadAll(result.Body)
		if err != nil {
			return nil, err
		}

		return data, nil
	}
}

/*var queue sync.Map

func Start() {
	go func() {
		log.Println("DOCDB START SCHEDULER")
		for {
			queue.Range(func(key, value interface{}) bool {
				start := time.Now()
				compressedData, err := cmhp.DataCompress(value.([]byte))
				if err != nil {
					log.Printf("Can't compress... %v\n", key)
					return true
				}

				ioutil.WriteFile(key.(string), compressedData, 0777)
				log.Printf("Store %v\n", key)
				queue.Delete(key)
				log.Printf("%v\n", time.Since(start))
				return true
			})
			time.Sleep(time.Second * 1)
		}
	}()
}

func Get(scope string, recordName string, out interface{}) error {
	// Build path
	finalPath := scope
	if scope != "" {
		finalPath += "/"
	}
	finalPath += recordName + ".data"

	v, ok := queue.Load(finalPath)
	if ok {
		json.Unmarshal(v.([]byte), out)
	} else {
		dataFile, err := os.OpenFile(finalPath, os.O_RDONLY, 0777)
		if err != nil {
			return err
		}
		defer dataFile.Close()
		data, _ := ioutil.ReadAll(dataFile)
		data, err = cmhp.DataDecompress(data)
		if err != nil {
			return err
		}

		json.Unmarshal(data, out)
	}

	return nil
}

func Save(scope string, recordName string, in interface{}) {
	// Build path
	finalPath := scope
	if scope != "" {
		finalPath += "/"
	}
	os.MkdirAll(finalPath, 0777)
	finalPath += recordName + ".data"
	data, _ := json.Marshal(in)
	queue.Store(finalPath, data)
}
*/
