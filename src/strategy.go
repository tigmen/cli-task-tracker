package main

const (
	TODO = iota
	INPROGRESS
	DONE
)

type Task struct {
	Id          uint
	Desctiption string
	Status      uint
	CreatedAt   int32
	UpdatedAt   int32
}

type Context struct {
	Task
	Command string
}

type Stratrgy interface {
	Execute(context Context)
}
