package docdb

/*type WriteInfo struct {
	Path        string
	Data        []byte
	IsCompress  bool
	Visibility  string
	ContentType string
}

var s3queue sync.Map
var s3Client *s3.S3
var s3Bucket string

func S3Start(spacesKey string, spacesSecret string, endpoint string, bucket string) {
	log.Println("INIT S3")
	newSession, err := session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(spacesKey, spacesSecret, ""),
		Endpoint:    aws.String(endpoint),
		Region:      aws.String("us-east-1"),
	})
	if err != nil {
		log.Fatal(err)
	}
	s3Client = s3.New(newSession)
	s3Bucket = bucket

	go func() {
		log.Println("DOCDB S3 START SCHEDULER")
		for {
			s3queue.Range(func(key, value interface{}) bool {
				info := value.(WriteInfo)

				data := info.Data

				if info.IsCompress {
					data, err = cmhp.DataCompress(info.Data)
					if err != nil {
						log.Printf("Can't compress... %v\n", key)
						return true
					}
				}

				log.Printf("Saving %v\n", key)

				object := s3.PutObjectInput{
					Bucket:      aws.String(bucket),
					Key:         aws.String(key.(string)),
					Body:        bytes.NewReader(data),
					ACL:         aws.String(info.Visibility),
					ContentType: aws.String(info.ContentType),
				}
				_, err = s3Client.PutObject(&object)

				if err != nil {
					log.Println(err.Error())
				} else {
					log.Printf("Saved %v\n", key)
					s3queue.Delete(key)
				}

				return true
			})

			time.Sleep(time.Second * 5)
		}
	}()
}

func S3Get(scope string, recordName string, out interface{}) error {
	// Build path
	finalPath := scope
	if scope != "" {
		finalPath += "/"
	}
	finalPath += recordName + ".data"

	v, ok := s3queue.Load(finalPath)
	if ok {
		err := json.Unmarshal(v.([]byte), out)
		if err != nil {
			return err
		}
	} else {
		// Download file
		result, err := s3Client.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(finalPath),
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
	}

	return nil
}

func S3Save(scope string, recordName string, in interface{}) {
	// Build path
	finalPath := scope
	if scope != "" {
		finalPath += "/"
	}

	finalPath += recordName + ".data"
	data, _ := json.Marshal(in)

	writeInfo := WriteInfo{
		Path:        finalPath,
		Data:        data,
		IsCompress:  true,
		ContentType: "binary/octet-stream",
		Visibility:  "private",
	}

	s3queue.Store(finalPath, writeInfo)
}

func S3SaveFile(path string, data []byte, isCompress bool, visibility string, contentType string) {
	writeInfo := WriteInfo{
		Path:        path,
		Data:        data,
		IsCompress:  isCompress,
		ContentType: contentType,
		Visibility:  visibility,
	}

	s3queue.Store(path, writeInfo)
}*/
