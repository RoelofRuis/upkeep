package main

import (
	"fmt"
	"os"
)

type FileIO struct {
	path string
}

func (fio FileIO) OpenForDay(day string) (*os.File, error) {
	return os.Open(fmt.Sprintf("%s/sheet_%s.csv", fio.path, day))
}

func (fio FileIO) CreateForDay(day string) (*os.File, error) {
	return os.Create(fmt.Sprintf("%s/sheet_%s.csv", fio.path, day))
}

func (fio FileIO) DeleteForDay(day string) error {
	return os.Remove(fmt.Sprintf("%s/sheet_%s.csv", fio.path, day))
}
