package filewriter

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

type FileWriter struct {
	filename string
}

func NewFileWriter(filename string) *FileWriter {
	return &FileWriter{filename: filename}
}

func (f *FileWriter) Write(data map[string]string) error {
	err := ioutil.WriteFile(fmt.Sprintf(f.filename, time.Now()), []byte(fmt.Sprint(data)), 0644)
	if err != nil {
		log.Println(err)
	}
	return err
}
