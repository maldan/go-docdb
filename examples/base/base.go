package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Result struct {
	Id          string   `json:"id"`
	Created     string   `json:"created"`
	PhotoIdList []string `json:"photoIdList"`
}

type Optimization struct {
	Id                   string   `json:"id"`
	Created              string   `json:"created"`
	ResultId             string   `json:"resultId"`
	PhotoIdList          []string `json:"photoIdList"`
	OptimizedPhotoIdList []string `json:"optimizedPhotoIdList"`
}

type Preset struct {
	Id      string `json:"id"`
	Created string `json:"created"`
}

type Info struct {
	BestLabel   string          `json:"bestLabel"`
	WebEntities [][]interface{} `json:"webEntities"`
	Labels      [][]interface{} `json:"labels"`

	FullMatch    []string `json:"fullMatch"`
	PartialMatch []string `json:"partialMatch"`
	Pages        []string `json:"pages"`
}

type Photo struct {
	Id      string `json:"id"`
	Created string `json:"created"`

	IsChanged  string `json:"isChanged"`
	IsCompress string `json:"isCompress"`

	Path         string `json:"path"`
	OriginalPath string `json:"originalPath"`

	Resolution         []int `json:"resolution"`
	OriginalResolution []int `json:"originalResolution"`

	Size         int `json:"size"`
	OriginalSize int `json:"originalSize"`

	Info Info `json:"info"`
}

type XX struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Balance   int    `json:"balance"`
	Created   string `json:"created"`

	ResultList       []Result       `json:"resultList"`
	OptimizationList []Optimization `json:"optimizationList"`
	PresetList       []Preset       `json:"presetList"`
	PhotoList        []Photo        `json:"photoList"`
}

/*func sex() {
	var yy = []byte{1, 2, 3}

	for i := 0; i < 100000; i++ {
		ioutil.WriteFile(fmt.Sprintf("../xx/%v", xid.New().String()), yy, 0777)
	}
}*/

func main() {
	/*go sex()
	go sex()
	go sex()
	go sex()
	go sex()
	go sex()
	go sex()
	go sex()

	time.Sleep(time.Second * 100)*/

	/*docdb.Start()

	start := time.Now()
	var xx XX
	docdb.Get("examples/lox", "y", &xx)
	fmt.Printf("%v\n", time.Since(start))
	fmt.Println(xx.Email)*/

	/*start = time.Now()
	docdb.Save("examples/lox", "y", &xx)
	fmt.Printf("%v\n", time.Since(start))*/

	/*start = time.Now()
	docdb.Save("examples/lox", "y", &xx)
	fmt.Printf("%v\n", time.Since(start))

	time.Sleep(time.Second * 5)*/

	start := time.Now()
	dataFile, err := os.OpenFile("../xx/c36e6m44hjv37541rjp0", os.O_RDONLY, 0777)
	if err != nil {
		fmt.Println(err)
	}
	defer dataFile.Close()
	data, _ := ioutil.ReadAll(dataFile)
	fmt.Println(data)
	fmt.Printf("%v\n", time.Since(start))

	files, err := ioutil.ReadDir("../xx")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(len(files))
}
