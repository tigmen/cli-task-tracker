package main

import (
	"flag"
	"fmt"
	"strconv"
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
	err := c.flagset.Parse(args)
	if err != nil {
		return fmt.Errorf("Failed parse flagset: %w", err)
	}

	args = c.flagset.Args()
	if len(args) < 1 {
		return fmt.Errorf("Need 1 args: description")
	}

	c.description = args[0]
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
	description string
	id          uint
}

func NewCommandUpdate() *Command_update {
	cmd := &Command_update{
		Command: Command{
			flagset: flag.NewFlagSet("update", flag.ContinueOnError),
		},
	}

	return cmd
}

func (c *Command_update) Init(args []string) error {
	err := c.flagset.Parse(args)
	if err != nil {
		return fmt.Errorf("Failed parse flagset: %w", err)
	}

	args = c.flagset.Args()
	if len(args) < 2 {
		return fmt.Errorf("Need 2 args: id and new description")
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("Failed parse uint id from %s: %w", args[0], err)
	}

	c.id = uint(id)
	c.description = args[1]

	return nil
}

func (c Command_update) Execute() (string, error) {
	_, err := c.storage.Update(c.id, Task{Desctiption: c.description})
	if err != nil {
		return "", fmt.Errorf("Failed to execute command_update: %w", err)
	}

	return "", nil
}

type Command_delete struct {
	Command
	id uint
}

func NewCommandDelete() *Command_delete {
	cmd := &Command_delete{
		Command: Command{
			flagset: flag.NewFlagSet("delete", flag.ContinueOnError),
		},
	}

	return cmd
}

func (c *Command_delete) Init(args []string) error {
	err := c.flagset.Parse(args)
	if err != nil {
		return fmt.Errorf("Failed parse flagset: %w", err)
	}

	args = c.flagset.Args()
	if len(args) < 1 {
		return fmt.Errorf("Need 1 args: id")
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("Failed parse uint id from %s: %w", args[0], err)
	}

	c.id = uint(id)

	return nil
}

func (c Command_delete) Execute() (string, error) {
	_, err := c.storage.Delete(c.id)
	if err != nil {
		return "", fmt.Errorf("Failed to execute command_delete: %w", err)
	}

	return "", nil
}

type Command_mark struct {
	Command
	id   uint
	mark uint
}

func NewCommandMark() *Command_mark {
	cmd := &Command_mark{
		Command: Command{
			flagset: flag.NewFlagSet("mark", flag.ContinueOnError),
		},
	}

	return cmd
}

func (c *Command_mark) Init(args []string) error {
	c.flagset.Parse(args)
	if len(args) < 2 {
		return fmt.Errorf("Need 2 args: id, mark")
	}

	id, err := strconv.ParseUint(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("Failed parse uint id from %s: %w", args[0], err)
	}

	c.id = uint(id)

	switch args[1] {
	case "todo":
		c.mark = TODO
	case "inprogress":
		c.mark = INPROGRESS
	case "done":
		c.mark = DONE
	default:
		return fmt.Errorf("No such mark: %s", args[1])
	}

	return nil
}

func (c Command_mark) Execute() (string, error) {
	_, err := c.storage.Update(c.id, Task{Status: c.mark})
	if err != nil {
		return "", fmt.Errorf("Failed to execute command_mark: %w", err)
	}

	return "", nil
}

type Command_list struct {
	Command
	mark uint
}

func NewCommandList() *Command_list {
	cmd := &Command_list{
		Command: Command{
			flagset: flag.NewFlagSet("list", flag.ContinueOnError),
		},
	}

	return cmd
}

func (c *Command_list) Init(args []string) error {
	mark := c.flagset.String("status", "all", "Show tasks with mark")
	err := c.flagset.Parse(args)
	if err != nil {
		return fmt.Errorf("Failed parse flagset: %w", err)
	}

	switch *mark {
	case "all":
		c.mark = ALL
	case "todo":
		c.mark = TODO
	case "inprogress":
		c.mark = INPROGRESS
	case "done":
		c.mark = DONE
	default:
		return fmt.Errorf("Failed to parse flagse: unknown mark")
	}

	return nil
}

func status(status uint) (string, error) {
	switch status {
	case TODO:
		return "todo", nil
	case INPROGRESS:
		return "inprogress", nil
	case DONE:
		return "done", nil
	}

	return "", fmt.Errorf("Unknown status")
}

func (c Command_list) Execute() (string, error) {
	tasks, err := c.storage.GetAll()
	if err != nil {
		return "", fmt.Errorf("Failed Command_list: %w", err)
	}

	var out string

	for _, val := range tasks {
		if c.mark == ALL || c.mark == val.Status {
			status, err := status(val.Status)
			if err != nil {
				return "", fmt.Errorf("Failed parse status: %w", err)
			}

			out += fmt.Sprintf("#Id: %d\t%s\nDescription: %s\nCreated: %s\nUpdated: %s\n",
				val.Id, status, val.Desctiption,
				time.Unix(val.CreatedAt, 0).Format(time.RFC850), time.Unix(val.UpdatedAt, 0).Format(time.RFC850))
		}
	}

	return out, nil
}
