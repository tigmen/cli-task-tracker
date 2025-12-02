package main

import (
	"fmt"
	"log"
)

func main() {
	var description string
	fmt.Scan(&description)
	var cmd Stratrgy = Command_add{storage: FileStorage{path: "storage.txt"}, description: description}

	out, err := cmd.Execute(Context{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print(out)
}
