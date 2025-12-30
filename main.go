package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

// Task представляет собой отдельную задачу
type Task struct {
	Id          int    `json:"id"`                     // Уникальный идентификатор задачи
	Content     string `json:"content"`                // Текст задачи
	Done        bool   `json:"done"`                   // Статус выполнения
	CreatedAt   string `json:"created_at"`             // Дата и время создания
	CompletedAt string `json:"completed_at,omitempty"` // Дата и время завершения (если выполнена)
}

// TodoList содержит список всех задач и информацию о следующем доступном ID
type TodoList struct {
	Tasks  []Task `json:"tasks"`   // Список задач
	NextId int    `json:"next_id"` // Следующий доступный ID для новой задачи
}

const maxTaskLength = 200      // Максимальная длина текста задачи в символах
const tasksPath = "tasks.json" // Путь к файлу для хранения задач

// loadTasks загружает список задач из файла
// Если файл не существует, создается новый пустой список
func loadTasks() (*TodoList, error) {
	data, err := os.ReadFile(tasksPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &TodoList{NextId: 1}, nil
		}

		return nil, err
	}

	var tl TodoList
	if err := json.Unmarshal(data, &tl); err != nil {
		return nil, err
	}

	return &tl, nil
}

// saveTask сохраняет текущий список задач в файл
func saveTask(tl *TodoList) error {
	data, err := json.MarshalIndent(tl, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(tasksPath, data, 0644)
}

// parseTaskId преобразует строковый ID в числовой и проверяет его корректность
func parseTaskId(strId string) (int, bool) {
	id, err := strconv.Atoi(strId)
	if err != nil {
		fmt.Println("Ошибка: не верный id")
		return 0, false
	}

	return id, true
}

// findTaskIndex находит индекс задачи по её ID
// Возвращает -1, если задача не найдена
func findTaskIndex(tl *TodoList, id int) int {
	for i := range tl.Tasks {
		if tl.Tasks[i].Id == id {
			return i
		}
	}

	return -1
}

// validateTask проверяет корректность задачи перед добавлением или редактированием
func validateTask(tl *TodoList, task Task) error {
	if len(task.Content) > maxTaskLength {
		return fmt.Errorf("Ошибка: текст задачи не должен превышать %d символов\n", maxTaskLength)
	}

	if strings.TrimSpace(task.Content) == "" {
		return fmt.Errorf("Ошибка: новый текст задачи не может быть пустым")
	}

	for _, t := range tl.Tasks {
		if strings.EqualFold(t.Content, task.Content) {
			return fmt.Errorf("Ошибка: задача с таким заголовком уже существует")
		}
	}

	return nil
}

// listTasks выводит список всех задач с их статусами
func listTasks(tl *TodoList) {
	if len(tl.Tasks) == 0 {
		fmt.Println("Список задач пуст")
		return
	}

	fmt.Println("Список задач:")
	for _, task := range tl.Tasks {
		status := " "
		if task.Done {
			status = "x"
		}

		fmt.Printf("%d [%s], %s (создана: %s)", task.Id, status, task.Content, task.CreatedAt)
		if task.Done && task.CompletedAt != "" {
			fmt.Printf(", выполнена: %s", task.CompletedAt)
		}

		fmt.Println()
	}
}

// addTask добавляет новую задачу в список
func addTask(tl *TodoList, content string) {
	task := Task{
		Id:        tl.NextId,
		Content:   content,
		Done:      false,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	err := validateTask(tl, task)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tl.Tasks = append(tl.Tasks, task)
	tl.NextId++
	fmt.Printf("Добавлена задача %d: %s\n", task.Id, content)
}

// toggleTask изменяет статус выполнения задачи (выполнено/не выполнено)
func toggleTask(tl *TodoList, strId string) {
	id, ok := parseTaskId(strId)
	if !ok {
		return
	}

	index := findTaskIndex(tl, id)
	if index == -1 {
		fmt.Println("Задача не найдена")
		return
	}

	tl.Tasks[index].Done = !tl.Tasks[index].Done
	tl.Tasks[index].CompletedAt = ""
	status := "не выполнено"
	if tl.Tasks[index].Done {
		status = "выполнено"
		tl.Tasks[index].CompletedAt = time.Now().Format("2006-01-02 15:04:05")
	}

	fmt.Printf("Задача #%d отмечена как %s\n", id, status)
}

// deleteTask удаляет задачу из списка по её ID
func deleteTask(tl *TodoList, strId string) {
	id, ok := parseTaskId(strId)
	if !ok {
		return
	}

	index := findTaskIndex(tl, id)
	if index == -1 {
		fmt.Println("Задача не найдена")
		return
	}

	tl.Tasks = append(tl.Tasks[:index], tl.Tasks[index+1:]...)
	fmt.Printf("Задача #%d была удалена\n", id)
}

// clearAllTasks удаляет все задачи и сбрасывает счётчик ID
func clearAllTasks(tl *TodoList) {
	tl.Tasks = []Task{}
	tl.NextId = 1
	fmt.Println("Все задачи очищены")
}

// completeAllTasks отмечает все задачи как выполненные
func completeAllTasks(tl *TodoList) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	for i := range tl.Tasks {
		if !tl.Tasks[i].Done {
			tl.Tasks[i].Done = true
			tl.Tasks[i].CompletedAt = currentTime
		}
	}

	fmt.Println("Все задачи отмечены как выполненные")
}

func main() {
	listFlag := flag.Bool("list", false, "List all tasks")
	addFlag := flag.String("add", "", "Add a new task")
	toggleFlag := flag.String("toggle", "", "Toggle task status (provide task ID)")
	deleteFlag := flag.String("delete", "", "Delete a task (provide task ID)")
	clearFlag := flag.Bool("clear", false, "Clear all tasks")
	completeAllFlag := flag.Bool("complete-all", false, "Mark all tasks as complete")

	flag.Parse()

	tl, err := loadTasks()
	if err != nil {
		fmt.Printf("Ошибка загрузки задач: %v\n", err)
		return
	}

	if *listFlag {
		listTasks(tl)
		return
	}

	if *addFlag != "" {
		content := *addFlag
		addTask(tl, content)
		if err := saveTask(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err.Error())
			return
		}
		return
	}

	if *toggleFlag != "" {
		id := *toggleFlag
		toggleTask(tl, id)
		if err := saveTask(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err.Error())
			return
		}
		return
	}

	if *deleteFlag != "" {
		id := *deleteFlag
		deleteTask(tl, id)
		if err := saveTask(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err)
			return
		}
		return
	}

	if *clearFlag {
		clearAllTasks(tl)
		if err := saveTask(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err)
			return
		}
		return
	}

	if *completeAllFlag {
		completeAllTasks(tl)
		if err := saveTask(tl); err != nil {
			fmt.Printf("Ошибка сохранения задач: %v\n", err)
			return
		}
		return
	}

	flag.Usage()
}
