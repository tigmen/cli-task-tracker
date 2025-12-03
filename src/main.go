package main

import (
	"fmt"
	"log"
	"os"
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

func main() {
	if len(os.Args) < 2 {
		log.Fatal("You must pass a sub-command")
	}

	cmds := []Strategy{
		NewCommandAdd(),
	}

	subcommand := os.Args[1]
	for _, cmd := range cmds {
		if cmd.Name() == subcommand {
			err := cmd.SetStorage(FileStorage{path: "storage.txt"})
			if err != nil {
				log.Fatal(err)
			}

			err = cmd.Init(os.Args[2:])
			if err != nil {
				log.Fatal(err)
			}

			out, err := cmd.Execute()
			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(out)
			os.Exit(0)
		}
	}
	log.Fatal("Unknown sub-command")

}
