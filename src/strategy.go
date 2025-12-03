package main

import (
	"flag"
	"fmt"
	"log"
	"time"
)

type Strategy interface {
	Name() string
	Init(args []string) error
	SetStorage(storage Storage) error
	Execute() (string, error)
}

type Command struct {
	storage Storage
	flagset *flag.FlagSet
}

func (c *Command) SetStorage(storage Storage) error {
	c.storage = storage
	return nil
}

func (c Command) Name() string {
	return c.flagset.Name()
}

type Command_add struct {
	Command
	description string
}

func (c Command_add) Execute() (string, error) {
	id, err := c.storage.Add(Task{Desctiption: c.description, Status: TODO,
		CreatedAt: time.Now().Unix(), UpdatedAt: time.Now().Unix()})
	if err != nil {
		return "", fmt.Errorf("Failed Command_add: %w", err)
	}

	return fmt.Sprintf("Task added successfully (ID: %d)", id), err
}

func (c *Command_add) Init(args []string) error {
	c.flagset.Parse(args)
	c.description = c.flagset.Arg(0)
	log.Println(c.flagset.Args())
	return nil
}

func NewCommandAdd() *Command_add {
	cmd := &Command_add{
		Command: Command{
			flagset: flag.NewFlagSet("add", flag.ContinueOnError),
		},
	}

	return cmd
}

type Command_update struct {
	Command
	id uint
}

func (c Command_update) Execute() (string, error) {
	return fmt.Sprintf(""), nil
}

type Command_delete struct {
	Command
}

func (c Command_delete) Execute() (string, error) {
	return fmt.Sprintf(""), nil
}

type Command_mark struct {
	Command
}

func (c Command_mark) Execute() (string, error) {
	return fmt.Sprintf(""), nil
}

type Command_list struct {
	Command
}

func (c Command_list) Execute() (string, error) {

	return fmt.Sprintf(""), nil
}

func (c Command_list) Init(args []string) error {
	c.flagset.Parse(args)
	return nil
}

func NewCommandList() *Command_list {
	cmd := &Command_list{
		Command: Command{
			flagset: flag.NewFlagSet("add", flag.ContinueOnError),
		},
	}

	return cmd
}
