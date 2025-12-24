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

- **toggle** - Поменять статус задачи (выполнено/не выполнено)
  ```bash
  ./golang-todo-cli toggle 1
  ```

- **delete** - Удалить задачу
  ```bash
  ./golang-todo-cli delete 1
  ```

## Пример работы

```bash
# Добавляем задачи
./golang-todo-cli add "Сделать домашнее задание"
./golang-todo-cli add "Позвонить маме"

# Просматриваем список
./golang-todo-cli list

# Отмечаем задачу как выполненную
./golang-todo-cli toggle 1

# Удаляем задачу
./golang-todo-cli delete 2
```

## Хранение данных

Все задачи сохраняются в файле `tasks.json` в формате JSON.

## Логирование

Приложение ведет лог в файле `app.log` с уровнем детализации INFO.

## Лицензия

MIT License - см. файл [LICENSE](LICENSE) для деталей.
