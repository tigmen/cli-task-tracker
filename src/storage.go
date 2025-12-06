package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Storage interface {
	Add(task Task) (uint, error)
	Update(id uint, task Task) (uint, error)
	Delete(id uint) (uint, error)
	Get(id uint) (Task, error)
	GetAll() ([]Task, error)
}

type FileStorage struct {
	path string
}

func (fs FileStorage) Add(task Task) (uint, error) {
	file, err := os.OpenFile(fs.path, os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return 0, fmt.Errorf("Failed to open storage-file: %w", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal("Failed to close storage-file: ", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	var id uint

	for scanner.Scan() {
		line := scanner.Bytes()

		var tmp Task

		err = json.Unmarshal(line, &tmp)
		if err != nil {
			return 0, fmt.Errorf("Failed read json: %w", err)
		}

		id = tmp.Id + 1
	}

	if err = scanner.Err(); err != nil {
		return 0, fmt.Errorf("Scanner error: %w", err)
	}

	if task.Desctiption == "" {
		task.Desctiption = "(null)"
	}

	if task.CreatedAt == 0 {
		task.CreatedAt = time.Now().Unix()
	}

	if task.UpdatedAt == 0 {
		task.UpdatedAt = time.Now().Unix()
	}

	if task.Id == 0 {
		task.Id = id
	}

	out, err := json.Marshal(task)
	if err != nil {
		return 0, fmt.Errorf("Failed json: %w", err)
	}

	_, err = fmt.Fprintln(file, string(out))
	if err != nil {
		return 0, fmt.Errorf("Error writing new task: %w", err)
	}

	return id, nil
}

func (fs FileStorage) rewrite(tasks []Task) error {
	file, err := os.OpenFile(fs.path, os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to open storage-file: %w", err)
	}

	defer func() {
		if err != nil {
			log.Fatal("Failed to close storage-file: ", err)
		}
	}()

	writer := bufio.NewWriter(file)

	for _, val := range tasks {
		out, err := json.Marshal(val)
		if err != nil {
			return fmt.Errorf("Failed json: %w", err)
		}

		_, err = fmt.Fprintln(writer, string(out))
		if err != nil {
			return fmt.Errorf("Error rewriting task: %w", err)
		}
	}

	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Failed write to storage-file: %w", err)
	}

	return nil
}

func (fs FileStorage) Update(id uint, task Task) (uint, error) {
	tasks, err := fs.GetAll()
	if err != nil {
		return 0, err
	}

	exist := false

	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == id {
			tasks[i].Id = id
			if task.Desctiption != "" {
				tasks[i].Desctiption = task.Desctiption
			}
			tasks[i].Status = task.Status
			tasks[i].UpdatedAt = time.Now().Unix()
			exist = true
		}
	}

	if exist {
		err = fs.rewrite(tasks)
		if err != nil {
			return 0, err
		}

		return id, nil
	}

	id, err = fs.Add(task)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (fs FileStorage) Delete(id uint) (uint, error) {
	tasks, err := fs.GetAll()
	if err != nil {
		return 0, err
	}

	new_tasks := make([]Task, len(tasks))
	j := 0

	for i := 0; i < len(tasks); i++ {
		if tasks[i].Id == id {
			continue
		}
		new_tasks[j] = tasks[i]
		j++
	}

	err = fs.rewrite(new_tasks[:j])
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (fs FileStorage) GetAll() ([]Task, error) {
	file, err := os.OpenFile(fs.path, os.O_CREATE|os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("Failed to open storage-file: %w", err)
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal("Failed to close storage-file: ", err)
		}
	}()

	scanner := bufio.NewScanner(file)
	tasks := make([]Task, 0)

	for scanner.Scan() {
		line := scanner.Bytes()

		var tmp Task

		err = json.Unmarshal(line, &tmp)
		if err != nil {
			return nil, fmt.Errorf("Failed read json: %w", err)
		}

		tasks = append(tasks, tmp)
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("Scanner error: %w", err)
	}

	return tasks, nil
}

func (fs FileStorage) Get(id uint) (Task, error) {
	tasks, err := fs.GetAll()
	if err != nil {
		return Task{}, err
	}

	for _, val := range tasks {
		if val.Id == id {
			return val, nil
		}
	}

	return Task{Id: 0}, nil
}
