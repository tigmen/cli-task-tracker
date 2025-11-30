package main

import (
	"fmt"
	"time"
)

const (
	TODO = iota
	INPROGRESS
	DONE
)

type Task struct {
	Id          uint
	Desctiption string
	Status      uint
	CreatedAt   int64
	UpdatedAt   int64
}

type Response struct {
	Task
	Next *Response
}

type Context struct {
	Args []string
}

type Stratrgy interface {
	Execute(context Context) (string, error)
}

type Command_add struct {
	path        string
	description string
}

func (c Command_add) Execute(context Context) (string, error) {
	var storage Storage = FileStorage{path: c.path}

	id, err := storage.Add(Task{Desctiption: c.description, Status: TODO,
		CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()})
	if err != nil {
	}

	return fmt.Sprintf("Task added successfully (ID: %d)", id), err
}

type Command_update struct {
	id uint
}

func (c Command_update) Execute(context Context) (string, error) {
	return fmt.Sprintf(""), nil
}

type Command_delete struct{}

func (c Command_delete) Execute(context Context) (string, error) {
	return fmt.Sprintf(""), nil
}

type Command_mark struct{}

func (c Command_mark) Execute(context Context) (string, error) {
	return fmt.Sprintf(""), nil
}

type Command_list struct{}

func (c Command_list) Execute(context Context) (string, error) {
	return fmt.Sprintf(""), nil
}
