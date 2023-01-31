# Project_Tasks

На данном репозитории представлен небольшой pet-project на языке Golang.

Программа задумывалась как некая напоминалка.  
В неё можно добавлять и удалять события, а также производить поиск по различным параметрам.
Событие записывается в базу данных в виде названия, даты, времени, продолжительности, описания.

Также есть возможность создавать списки с заметками (к примеру, список покупок).

**ВАЖНО!!!** На данный момент есть нерешённая проблема:
при смене языка во время ввода появляются лишние символы, что может привести к неправильному распознанию команд и вводов.
Поэтому язык следует менять до начала ввода, а лучше пока обойтись только английской раскладкой.

Работа с программой ведётся через терминал.  
Ниже приведено описание команд:

```evt add``` — добавить событие.
После ввода команды будет предложено ввести данные о событии: 
- название (минимум один любой символ)
- дата (в формате DD.MM.YYYY)
- время (в формате hh:mm)
- продолжительность (в формате hh:mm)
- описание.

Каждому событию присваивается уникальный id.

```evt dlt``` — удалить событие.
После ввода команды будет предложено ввести id события (натуральное число) и подтвердить удаление.

```evt shw [argument]``` — показать события.
В зависимости от аргумента далее будет предложено ввести различные данные для поиска:
- ```all``` — вывод всех событий.
- ```date``` — будет предложено ввести дату события (в формате DD.MM.YYYY). Вывод событий на конкретную дату.
- ```desc``` — будет предложено ввести описание события. Вывод событий по описанию.
- ```dur``` — будет предложено ввести продолжительность события (в формате hh:mm). Вывод событий по продолжительности.
- ```id``` — будет предложено ввести id события (натуральное число). Вывод события по id.
- ```intv``` — будет предложено ввести интервал (даты начала и конца интервала в формате DD.MM.YYYY).
Вывод событий в заданном интервале (включая крайние даты).
- ```name``` — будет предложено ввести название события (минимум один любой символ). Вывод событий по имени.
- ```time``` — будет предложено ввести дату и время события (в формате DD:MM:YYYY и hh:mm соответственно).
Вывод событий по дате и времени (выводятся события, которые активны, то есть учитывает продолжительность).

```ext``` — выход из программы.

```help``` — помощь.  
При вводе 'help' будет выведена инструкция по всем командам.  
Если нужна инструкция по конкретной команде, введите 'help [команда]'.

```lst add``` — добавить список.
После ввода будет предложено ввести название нового списка (минимум один любой символ).

```lst dlt``` — удалить список.
После ввода команды будет предложено ввести название списка (минимум один любой символ) и подтвердить удаление.

```lst shw``` — показать список.
После ввода команды будет предложено ввести название списка (минимум один любой символ). Вывод всего списка.

```lsts shw``` — показать все списки.
Вывод всех списков.

```prch add``` — добавить заметку.
После ввода команды будет предложено ввести данные о заметке:
- название (минимум один любой символ)
- описание.

Каждой заметке присваивается уникальный id.

```prch dlt``` — удалить заметку.
После ввода команды будет предложено ввести данные о заметке и подтвердить удаление:
- название списка (минимум один любой символ)
- id заметки (натуральное число).