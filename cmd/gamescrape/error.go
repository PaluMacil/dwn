package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
)

type ErrorList []HttpResponseError

func (el ErrorList) WriteFile() error {
	errorData, err := json.MarshalIndent(el, "", "\t")
	if err != nil {
		log.Fatalf("marshalling errorList: %s", err)
	}
	err = ioutil.WriteFile(filepath.Join("data", "errors.json"), errorData, 0666)
	return err
}

func (el ErrorList) FromFile() ErrorList {
	var e ErrorList
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatalf("File error: %v\n", err)
	}
	err = json.Unmarshal(file, &e)
	if err != nil {
		log.Fatalf("unmarshalling error file: %v\n", err)
	}
	return e
}

func (el *ErrorList) Add(e HttpResponseError) {
	errorLock.Lock()
	defer errorLock.Unlock()
	*el = append(*el, e)
}

var errorLock sync.Mutex

type HttpResponseError struct {
	Code  int
	Page  int
	Body  string
	Error string
}
