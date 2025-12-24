// Package main реализует CLI приложение для управления задачами (todo list)
// с поддержкой добавления, редактирования, удаления и отслеживания статуса задач.
package main

import (
	"encoding/json" // Для работы с JSON файлами
	"flag"          // Для парсинга аргументов командной строки
	"fmt"           // Для форматированного вывода
	"log"           // Для логирования
	"os"            // Для работы с файловой системой
	"strconv"       // Для конвертации строк в числа
	"strings"       // Для работы со строками
	"time"          // Для работы с временем
)

// Task представляет собой отдельную задачу в списке дел.
// Содержит всю необходимую информацию о задаче, включая статус и временные метки.
type Task struct {
	Id          int    `json:"id"`                     // Уникальный идентификатор задачи
	Content     string `json:"content"`                // Текстовое описание задачи
	Done        bool   `json:"done"`                   // Статус выполнения задачи
	CreatedAt   string `json:"created_at"`             // Время создания задачи (формат: "YYYY-MM-DD HH:MM:SS")
	CompletedAt string `json:"completed_at,omitempty"` // Время выполнения задачи (пусто, если задача не выполнена)
}

// TodoList представляет собой полный список задач.
// Используется для хранения всех задач и управления их идентификаторами.
type TodoList struct {
	Tasks  []Task `json:"tasks"`   // Список всех задач
	NextId int    `json:"next_id"` // Следующий доступный ID для новой задачи (автоинкремент)
}

// Константы приложения
const maxTaskLength = 200      // Максимально допустимая длина текста задачи в символах
const tasksPath = "tasks.json" // Путь к файлу для хранения задач в формате JSON

// Глобальные переменные
var logger *log.Logger // Логгер для записи событий и ошибок в файл app.log

// loadTasks загружает список задач из JSON файла.
// Если файл не существует, создается новый пустой список задач.
// Возвращает указатель на TodoList и ошибку (если возникла).
func loadTasks() (*TodoList, error) {
	logger.Println("Получение всех задач")

	data, err := os.ReadFile(tasksPath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Println("Задачи не найдены")

			// Если файл не существует, возвращаем новый пустой список с NextId = 1
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

// saveTask сохраняет текущий список задач в JSON файл.
// Использует форматированный вывод для удобства чтения.
// Возвращает ошибку, если сохранение не удалось.
func saveTask(t1 *TodoList) error {
	logger.Println("Сохранение задачи")

	// Маршалинг с отступами для удобочитаемости
	data, err := json.MarshalIndent(t1, "", "  ")
	if err != nil {
		logger.Print(err.Error())

		return err
	}

	// Сохранение в файл с правами 0644 (чтение/запись для владельца, только чтение для остальных)
	return os.WriteFile(tasksPath, data, 0644)
}

// listTasks отображает список всех задач в консоли.
// Для каждой задачи показывает ID, статус, текст и временные метки.
// Если список пуст, выводит соответствующее сообщение.
func listTasks(t1 *TodoList) {
	if len(t1.Tasks) == 0 {
		logger.Println("Список задач пуст")

		fmt.Println("Список задач пуст")
		return
	}
	logger.Println("Выведен список всех задач")

	fmt.Println("Список задач:")
	for _, task := range t1.Tasks {
		// Определяем символ статуса: "x" для выполненных, " " для невыполненных
		status := " "
		if task.Done {
			status = "x"
		}
		// Выводим информацию о задаче с временем создания
		fmt.Printf("%d [%s], %s (создана: %s)", task.Id, status, task.Content, task.CreatedAt)
		// Если задача выполнена, добавляем время выполнения
		if task.Done && task.CompletedAt != "" {
			fmt.Printf(", выполнена: %s", task.CompletedAt)
		}
		fmt.Println()
	}
}

// parseTaskId парсит строковый ID задачи в целое число.
// Возвращает ID и булево значение, указывающее на успех парсинга.
// Используется для валидации пользовательского ввода.
func parseTaskId(strId string) (int, bool) {
	id, err := strconv.Atoi(strId)
	if err != nil {
		logger.Println("Не верный id")
		fmt.Println("Ошибка: не верный id")
		return 0, false
	}
	return id, true
}

// findTaskIndex ищет задачу по ID и возвращает её индекс в списке.
// Если задача не найдена, возвращает -1.
// Используется во многих функциях для поиска задачи перед её обработкой.
func findTaskIndex(t1 *TodoList, id int) int {
	for i := range t1.Tasks {
		if t1.Tasks[i].Id == id {
			return i
		}
	}
	return -1
}

// isDuplicateTask проверяет, существует ли уже задача с таким же текстом.
// Сравнение нечувствительно к регистру (использует strings.EqualFold).
// Возвращает true, если дубликат найден, false - если нет.
// Используется при добавлении и редактировании задач для предотвращения дубликатов.
func isDuplicateTask(t1 *TodoList, content string) bool {
	for _, task := range t1.Tasks {
		if strings.EqualFold(task.Content, content) {
			return true
		}
	}
	return false
}

// addTask добавляет новую задачу в список.
// Выполняет валидацию: проверка на дубликаты и максимальную длину.
// Устанавливает текущее время в CreatedAt и статус Done = false.
// После добавления увеличивает NextId для следующей задачи.
func addTask(t1 *TodoList, content string) {
	// Проверка на дубликаты
	if isDuplicateTask(t1, content) {
		fmt.Println("Ошибка: задача с таким текстом уже существует")
		logger.Printf("Ошибка: попытка добавить дубликат задачи: %s\n", content)
		return
	}

	// Проверка максимальной длины
	if len(content) > maxTaskLength {
		fmt.Printf("Ошибка: текст задачи не должен превышать %d символов\n", maxTaskLength)
		logger.Printf("Ошибка: текст задачи превышает лимит в %d символов\n", maxTaskLength)
		return
	}

	// Создание новой задачи
	task := Task{
		Id:        t1.NextId,
		Content:   content,
		Done:      false,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"), // Текущее время
	}
	t1.Tasks = append(t1.Tasks, task)
	t1.NextId++ // Увеличение счетчика для следующей задачи

	logger.Printf("Добавлена задача %d: %s\n", task.Id, content)
	fmt.Printf("Добавлена задача %d: %s\n", task.Id, content)
}

// toggleTask переключает статус выполнения задачи (выполнено/не выполнено).
// Если задача отмечается как выполненная, устанавливается текущее время в CompletedAt.
// Если задача отмечается как невыполненная, CompletedAt очищается.
// Использует parseTaskId и findTaskIndex для валидации и поиска задачи.
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

	// Переключение статуса
	t1.Tasks[index].Done = !t1.Tasks[index].Done
	t1.Tasks[index].CompletedAt = ""
	status := "не выполнено"
	if t1.Tasks[index].Done {
		status = "выполнено"
		t1.Tasks[index].CompletedAt = time.Now().Format("2006-01-02 15:04:05") // Устанавливаем время выполнения
	}

	logger.Printf("Задача #%d отмечена как %s\n", id, status)
	fmt.Printf("Задача #%d отмечена как %s\n", id, status)
}

// editTask редактирует текст существующей задачи.
// Выполняет валидацию: проверка на пустой текст и дубликаты.
// Сохраняет оригинальные даты создания и выполнения.
// Использует parseTaskId и findTaskIndex для валидации и поиска задачи.
func editTask(t1 *TodoList, strId string, newContent string) {
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

	// Проверка на пустой текст
	if newContent == "" {
		fmt.Println("Ошибка: новый текст задачи не может быть пустым")
		logger.Printf("Ошибка: новый текст задачи #%d пустой\n", id)
		return
	}

	// Проверка на дубликаты
	if isDuplicateTask(t1, newContent) {
		fmt.Println("Ошибка: задача с таким текстом уже существует")
		logger.Printf("Ошибка: попытка добавить дубликат задачи: %s\n", newContent)
		return
	}

	// Редактирование задачи (даты сохраняются автоматически)
	t1.Tasks[index].Content = newContent

	logger.Printf("Задача #%d отредактирована: %s\n", id, newContent)
	fmt.Printf("Задача #%d отредактирована: %s\n", id, newContent)
}

// deleteTask удаляет задачу из списка по её ID.
// Использует parseTaskId и findTaskIndex для валидации и поиска задачи.
// Удаление выполняется с помощью среза (append), что эффективно для небольших списков.
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
	// Удаление элемента из среза: [элементы до индекса] + [элементы после индекса]
	t1.Tasks = append(t1.Tasks[:index], t1.Tasks[index+1:]...)
	fmt.Printf("Задача #%d была удалена\n", id)
}

// clearAllTasks очищает все задачи и сбрасывает счетчик ID.
// Устанавливает пустой срез задач и NextId = 1 для начала с чистого листа.
// Полезно для полного сброса списка задач.
func clearAllTasks(t1 *TodoList) {
	t1.Tasks = []Task{} // Очистка списка задач
	t1.NextId = 1       // Сброс счетчика ID
	logger.Println("Все задачи очищены")
	fmt.Println("Все задачи очищены")
}

// completeAllTasks отмечает все невыполненные задачи как выполненные.
// Устанавливает текущее время в CompletedAt для всех задач, которые еще не выполнены.
// Полезно для массового завершения задач.
func completeAllTasks(t1 *TodoList) {
	currentTime := time.Now().Format("2006-01-02 15:04:05") // Текущее время для всех задач
	for i := range t1.Tasks {
		if !t1.Tasks[i].Done {
			t1.Tasks[i].Done = true
			t1.Tasks[i].CompletedAt = currentTime // Устанавливаем время выполнения
		}
	}
	logger.Println("Все задачи отмечены как выполненные")
	fmt.Println("Все задачи отмечены как выполненные")
}

// initLogger инициализирует логгер для записи событий в файл.
// Создает или открывает файл app.log с флагами:
// - O_CREATE: создать файл, если не существует
// - O_WRONLY: открыть только для записи
// - O_APPEND: добавлять записи в конец файла
// Устанавливает префикс "INFO: " и формат с датой, временем и именем файла.
func initLogger() *os.File {
	const logPath = "app.log" // Путь к лог-файлу
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Ошибка: OpenFile: %v\n", err)
	}

	// Настройка логгера с префиксом и флагами форматирования
	logger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	return file
}

// main - основная функция приложения, точка входа.
// Инициализирует логгер, парсит аргументы командной строки и выполняет соответствующие действия.
// Поддерживает следующие команды: list, add, toggle, delete, edit, clear-all, complete-all.
// Использует паттерн "flag" для парсинга аргументов и обработки ошибок.
func main() {
	// Инициализация логгера (открытие файла и настройка)
	file := initLogger()
	// defer обеспечивает закрытие файла при выходе из функции (даже при панике)
	defer file.Close()

	// Инициализация флагов для парсинга аргументов команд
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	toggleCmd := flag.NewFlagSet("toggle", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)

	// Проверка минимального количества аргументов
	if len(os.Args) < 2 {
		// Вывод справки при недостаточном количестве аргументов
		fmt.Println("Использование: todo [list|add|toggle|delete|edit|clear-all|complete-all] [аргументы]")
		fmt.Println("  todo list         - показать все задачи")
		fmt.Println("  todo add \"текст\"  - добавить задачу")
		fmt.Println("  todo toggle 5     - поменять статус задачи #5")
		fmt.Println("  todo delete 5     - удалить задачу #5")
		fmt.Println("  todo edit 5 \"новый текст\" - редактировать задачу #5")
		fmt.Println("  todo clear-all    - очистить все задачи")
		fmt.Println("  todo complete-all - отметить все задачи как выполненные")
		return
	}

	// Основной switch по командам
	switch os.Args[1] {
	case "list":
		// Загрузка и отображение списка задач
		t1, err := loadTasks()
		if err != nil {
			fmt.Printf("Ошибка: loadTasks: %v\n", err.Error())
			logger.Printf("Ошибка: loadTasks: %v\n", err.Error())
			return
		}
		listTasks(t1)

	case "add":
		// Парсинг аргументов и добавление новой задачи
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
		// Парсинг аргументов и переключение статуса задачи
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
		// Парсинг аргументов и удаление задачи
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

	case "clear-all":
		// Очистка всех задач
		t1, err := loadTasks()
		if err != nil {
			fmt.Printf("Ошибка: loadTasks: %v\n", err)
			logger.Printf("Ошибка: loadTasks: %v\n", err)
			return
		}
		clearAllTasks(t1)
		if err := saveTask(t1); err != nil {
			fmt.Printf("Ошибка: saveTask: %v\n", err)
			logger.Printf("Ошибка: saveTask: %v\n", err)
			return
		}

	case "edit":
		// Парсинг аргументов и редактирование задачи
		editCmd.Parse(os.Args[2:])
		if len(editCmd.Args()) < 2 {
			fmt.Println("Укажите id задачи и новый текст")
			logger.Println("Ошибка: не указан id задачи или новый текст")
			return
		}
		t1, err := loadTasks()
		if err != nil {
			fmt.Printf("Ошибка: loadTasks: %v\n", err)
			logger.Printf("Ошибка: loadTasks: %v\n", err)
			return
		}
		newContent := strings.Join(editCmd.Args()[1:], " ")
		editTask(t1, editCmd.Args()[0], newContent)
		if err := saveTask(t1); err != nil {
			fmt.Printf("Ошибка: saveTask: %v\n", err)
			logger.Printf("Ошибка: saveTask: %v\n", err)
			return
		}

	case "complete-all":
		// Отметка всех задач как выполненных
		t1, err := loadTasks()
		if err != nil {
			fmt.Printf("Ошибка: loadTasks: %v\n", err)
			logger.Printf("Ошибка: loadTasks: %v\n", err)
			return
		}
		completeAllTasks(t1)
		if err := saveTask(t1); err != nil {
			fmt.Printf("Ошибка: saveTask: %v\n", err)
			logger.Printf("Ошибка: saveTask: %v\n", err)
			return
		}

	default:
		// Обработка неизвестных команд
		fmt.Println("Неизвестная команда")
		logger.Println("Введена не известная команда -", os.Args[1])
	}
}
