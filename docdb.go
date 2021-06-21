package docdb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var queue sync.Map

func Start() {
	go func() {
		fmt.Println("DOCDB START SCHEDULER")
		for {
			queue.Range(func(key, value interface{}) bool {
				start := time.Now()
				compressedData := compress(value.([]byte))
				ioutil.WriteFile(key.(string), compressedData, 0777)
				fmt.Printf("Store %v\n", key)
				queue.Delete(key)
				fmt.Printf("%v\n", time.Since(start))
				return true
			})
			time.Sleep(time.Second * 1)
		}
	}()
}

func Get(scope string, recordName string, out interface{}) interface{} {
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
			fmt.Println(err)
		}
		defer dataFile.Close()
		data, _ := ioutil.ReadAll(dataFile)
		data = decompress(data)

		json.Unmarshal(data, out)
	}

	return out
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
