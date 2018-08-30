package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type ErrorList []HttpResponseError

func (errorList ErrorList) WriteFile() error {
	errorData, err := json.MarshalIndent(errorList, "", "\t")
	if err != nil {
		log.Fatalf("marshalling errorList: %s", err)
	}
	err = ioutil.WriteFile("errors.json", errorData, 0666)
	return err
}

func (errorList ErrorList) Add(e HttpResponseError) {
	errorLock.Lock()
	defer errorLock.Unlock()
	errorList = append(errorList, e)
}

var errorLock sync.Mutex

type HttpResponseError struct {
	Code  int
	Body  string
	Error string
}
