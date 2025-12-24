package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Task struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

type TodoList struct {
	Tasks  []Task `json:"tasks"`
	NextId int    `json:"next_id"`
}

const tasksPath = "tasks.json"

var logger *log.Logger

func loadTasks() (*TodoList, error) {
	logger.Println("Получение всех задач")

	data, err := os.ReadFile(tasksPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Println("Задачи не найдены")

			return &TodoList{NextId: 1}, nil
		}
		logger.Println(err.Error())

		return nil, err
	}

	var t1 TodoList
	if err := json.Unmarshal(data, &t1); err != nil {
		logger.Println(err.Error())

		return nil, err
	}

	return &t1, nil
}

func saveTask(t1 *TodoList) error {
	logger.Println("Сохранение задачи")

	data, err := json.MarshalIndent(t1, "", "  ")
	if err != nil {
		logger.Print(err.Error())

		return err
	}

	return os.WriteFile(tasksPath, data, 0644)
}

func listTasks(t1 *TodoList) {
	if len(t1.Tasks) == 0 {
		logger.Println("Список задач пуст")

		fmt.Println("Список задач пуст")
		return
	}
	logger.Println("Выведен список всех задач")

	fmt.Println("Список задач:")
	for _, task := range t1.Tasks {
		status := " "
		if task.Done {
			status = "x"
		}
		fmt.Printf("%d [%s], %s\n", task.Id, status, task.Content)
	}
}

func parseTaskId(strId string) (int, bool) {
	id, err := strconv.Atoi(strId)
	if err != nil {
		logger.Println("Не верный id")
		fmt.Println("Ошибка: не верный id")
		return 0, false
	}
	return id, true
}

func findTaskIndex(t1 *TodoList, id int) int {
	for i := range t1.Tasks {
		if t1.Tasks[i].Id == id {
			return i
		}
	}
	return -1
}

func addTask(t1 *TodoList, content string) {
	task := Task{
		Id:      t1.NextId,
		Content: content,
		Done:    false,
	}
	t1.Tasks = append(t1.Tasks, task)
	t1.NextId++
	logger.Printf("Добавлена задача %d: %s\n", task.Id, content)

	fmt.Printf("Добавлена задача %d: %s\n", task.Id, content)
}

func toggleTask(t1 *TodoList, strId string) {
	id, ok := parseTaskId(strId)
	if !ok {
		return
	}

	index := findTaskIndex(t1, id)
	if index == -1 {
		logger.Printf("Задача #%d не найдена\n", id)
		fmt.Println("Задача не найдена")
		return
	}

	t1.Tasks[index].Done = !t1.Tasks[index].Done
	status := "выполнено"
	if !t1.Tasks[index].Done {
		status = "не выполнено"
	}
	logger.Printf("Задача #%d отмечена как %s\n", id, status)
	fmt.Printf("Задача #%d отмечена как %s\n", id, status)
}

func deleteTask(t1 *TodoList, strId string) {
	id, ok := parseTaskId(strId)
	if !ok {
		return
	}

	index := findTaskIndex(t1, id)
	if index == -1 {
		logger.Printf("Задача #%d не найдена\n", id)
		fmt.Println("Задача не найдена")
		return
	}

	logger.Printf("Задача #%d была удалена\n", id)
	t1.Tasks = append(t1.Tasks[:index], t1.Tasks[index+1:]...)
	fmt.Printf("Задача #%d была удалена\n", id)
}

func initLogger() *os.File {
	const logPath = "app.log"
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка: OpenFile: %v\n", err)
	}

	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return file
}

func main() {
	file := initLogger()
	defer file.Close()

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	toggleCmd := flag.NewFlagSet("toggle", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Использование: todo [list|add|toggle|delete] [аргументы]")
		fmt.Println("  todo list     - показать все задачи")
		fmt.Println("  todo add \"текст\" - добавить задачу")
		fmt.Println("  todo toggle 5 - поменять статус задачи #5")
		fmt.Println("  todo delete 5 - удалить задачу #5")
		return
	}

	switch os.Args[1] {
	case "list":
		t1, err := loadTasks()
		if err != nil {
			fmt.Printf("Ошибка: loadTasks: %v\n", err.Error())
			logger.Printf("Ошибка: loadTasks: %v\n", err.Error())
			return
		}
		listTasks(t1)
	case "add":
		addCmd.Parse(os.Args[2:])
		content := strings.Join(addCmd.Args(), " ")
		if content == "" {
			fmt.Println("Ошибка: укажите текст задачи")
			logger.Printf("Ошибка: не указан текст задачи")
			return
		}
		t1, err := loadTasks()
		if err != nil {
			fmt.Printf("Ошибка: loadTasks: %v\n", err)
			logger.Printf("Ошибка: loadTasks: %v\n", err)
			return
		}
		addTask(t1, content)
		if err := saveTask(t1); err != nil {
			fmt.Printf("Ошибка: saveTask: %v\n", err)
			logger.Printf("Ошибка: saveTask: %v\n", err)
			return
		}
	case "toggle":
		toggleCmd.Parse(os.Args[2:])
		if len(toggleCmd.Args()) == 0 {
			fmt.Println("Укажите id задачи")
			logger.Println("Ошибка: не указан id задачи")
			return
		}
		t1, err := loadTasks()
		if err != nil {
			fmt.Printf("Ошибка: loadTasks: %v\n", err)
			logger.Printf("Ошибка: loadTasks: %v\n", err)
			return
		}
		toggleTask(t1, toggleCmd.Args()[0])
		if err := saveTask(t1); err != nil {
			fmt.Printf("Ошибка: saveTask: %v\n", err)
			logger.Printf("Ошибка: saveTask: %v\n", err)
			return
		}
	case "delete":
		deleteCmd.Parse(os.Args[2:])
		if len(deleteCmd.Args()) == 0 {
			fmt.Println("Укажите id задачи")
			logger.Println("Ошибка: не указан id задачи")
			return
		}
		t1, err := loadTasks()
		if err != nil {
			fmt.Printf("Ошибка: loadTasks: %v\n", err)
			logger.Printf("Ошибка: loadTasks: %v\n", err)
			return
		}
		deleteTask(t1, deleteCmd.Args()[0])
		if err := saveTask(t1); err != nil {
			fmt.Printf("Ошибка: saveTask: %v\n", err)
			logger.Printf("Ошибка: saveTask: %v\n", err)
			return
		}
	default:
		fmt.Println("Неизвестная команда")
		logger.Println("Введена не известная команда -", os.Args[1])
	}
}
