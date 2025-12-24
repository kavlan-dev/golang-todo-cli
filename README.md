# golang-todo-cli

Простое CLI приложение для управления задачами, написанное на Go.

## Установка

```bash
git clone https://github.com/kavlan-dev/golang-todo-cli.git
cd golang-todo-cli
go build
```

## Использование

```bash
./golang-todo-cli [команда] [аргументы]
```

### Доступные команды:

- **list** - Показать все задачи
  ```bash
  ./golang-todo-cli list
  ```

- **add** - Добавить новую задачу
  ```bash
  ./golang-todo-cli add "Купить молоко"
  ```

- **edit** - Редактировать задачу
  ```bash
  ./golang-todo-cli edit 1 "Купить хлеб и молоко"
  ```

- **toggle** - Поменять статус задачи (выполнено/не выполнено)
  ```bash
  ./golang-todo-cli toggle 1
  ```

- **delete** - Удалить задачу
  ```bash
  ./golang-todo-cli delete 1
  ```

- **clear-all** - Очистить все задачи
  ```bash
  ./golang-todo-cli clear-all
  ```

- **complete-all** - Отметить все задачи как выполненные
  ```bash
  ./golang-todo-cli complete-all
  ```

## Пример работы

```bash
# Добавляем задачи
./golang-todo-cli add "Сделать домашнее задание"
./golang-todo-cli add "Позвонить маме"

# Просматриваем список
./golang-todo-cli list

# Редактируем задачу
./golang-todo-cli edit 1 "Сделать домашнее задание по математике"

# Отмечаем задачу как выполненную
./golang-todo-cli toggle 1

# Удаляем задачу
./golang-todo-cli delete 2

# Очищаем все задачи
./golang-todo-cli clear-all

# Отмечаем все задачи как выполненные
./golang-todo-cli complete-all
```

## Хранение данных

Все задачи сохраняются в файле `tasks.json` в формате JSON.

## Ограничения

- Максимальная длина задачи: 200 символов
- Задачи с одинаковым текстом (без учета регистра) считаются дубликатами и не могут быть добавлены

## Логирование

Приложение ведет лог в файле `app.log` с уровнем детализации INFO.

## Лицензия

MIT License - см. файл [LICENSE](LICENSE) для деталей.
