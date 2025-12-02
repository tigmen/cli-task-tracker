package main

import (
	"fmt"
	"log"
)

func main() {
	var s Storage = FileStorage{path: "storage.txt"}
	id, err := s.Update(2, Task{Desctiption: "abc"})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(id)
}
