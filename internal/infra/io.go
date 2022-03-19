package infra

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type IO interface {
	Read(fname string, dst interface{}) error
	Write(fname string, src interface{}) error
	Delete(fname string) error
	Export(fname string, records [][]string) error
}

type InMemoryIO struct {
	Files   map[string]interface{}
	Exports map[string][][]string
}

func NewInMemoryIO() *InMemoryIO {
	return &InMemoryIO{
		Files:   make(map[string]interface{}),
		Exports: make(map[string][][]string),
	}
}

func (io *InMemoryIO) Read(fname string, std interface{}) error {
	data, has := io.Files[fname]
	if !has {
		return nil
	}
	std = data
	return nil
}

func (io *InMemoryIO) Write(fname string, src interface{}) error {
	io.Files[fname] = src
	return nil
}

func (io *InMemoryIO) Delete(fname string) error {
	delete(io.Files, fname)
	return nil
}

func (io *InMemoryIO) Export(fname string, records [][]string) error {
	io.Exports[fname] = records
	return nil
}

type FileIO struct {
	PrettyJson bool
	HomePath   string
	DataFolder string
}

func (io FileIO) Read(fname string, dst interface{}) error {
	fpath := path.Join(io.HomePath, io.DataFolder, fname)

	data, err := ioutil.ReadFile(fpath)
	if err != nil {
		switch {
		case os.IsNotExist(err):
			return nil
		default:
			return err
		}
	}

	if err := json.Unmarshal(data, dst); err != nil {
		return err
	}

	return nil
}

func (io FileIO) Write(fname string, src interface{}) error {
	var data []byte
	var err error
	if io.PrettyJson {
		data, err = json.MarshalIndent(src, "", "  ")
	} else {
		data, err = json.Marshal(src)
	}
	if err != nil {
		return nil
	}

	fpath := path.Join(io.HomePath, io.DataFolder, fname)

	if _, err := os.Stat(fpath); os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(fpath), 0o700)
		if err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(fpath, data, 0o700); err != nil {
		return err
	}

	return nil
}

func (io FileIO) Delete(fname string) error {
	fpath := path.Join(io.HomePath, io.DataFolder, fname)
	if err := os.Remove(fpath); !os.IsNotExist(err) {
		return err
	}

	return nil
}

func (io FileIO) Export(fname string, records [][]string) error {
	file, err := os.Create(fname)
	if err != nil {
		return err
	}
	defer file.Close()

	csvWriter := csv.NewWriter(file)

	err = csvWriter.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

type IOLoggerDecorator struct {
	Inner IO
}

func (io IOLoggerDecorator) Read(fname string, dst interface{}) error {
	_, _ = fmt.Fprintf(os.Stderr, "reading from [%s]\n", fname)
	return io.Inner.Read(fname, dst)
}

func (io IOLoggerDecorator) Write(fname string, src interface{}) error {
	_, _ = fmt.Fprintf(os.Stderr, "writing to [%s]\n", fname)
	return io.Inner.Write(fname, src)
}

func (io IOLoggerDecorator) Delete(fname string) error {
	_, _ = fmt.Fprintf(os.Stderr, "deleting [%s]\n", fname)
	return io.Inner.Delete(fname)
}

func (io IOLoggerDecorator) Export(fname string, records [][]string) error {
	_, _ = fmt.Fprintf(os.Stderr, "exporting to [%s]\n", fname)
	return io.Inner.Export(fname, records)
}
