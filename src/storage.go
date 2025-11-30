package main

import (
	"os"
)

type Storage interface {
	Add(task Task) (uint, error)
	Update(id uint, task Task) (Task, error)
	Delete(id uint) (uint, error)
	Get(id uint) (Task, error)
}

type FileStorage struct {
	path string
}

func (fs FileStorage) Add(task Task) (uint, error) {
	file, err := os.OpenFile(fs.path, os.O_RDWR, 0666)
	if err != nil {
	}
	defer file.Close()

	return 0, nil
}

func (fs FileStorage) Update(id uint, task Task) (Task, error) {
	return Task{}, nil
}

func (fs FileStorage) Delete(id uint) (uint, error) {
	return 0, nil
}

func (fs FileStorage) Get(id uint) (Task, error) {
	return Task{}, nil
}
