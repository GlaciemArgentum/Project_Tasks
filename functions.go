package main

import (
	"bufio"
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"regexp"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

type Event struct {
	id          int
	name        string
	myTime      string
	duration    string
	description string
}

type Field struct {
	name   string
	ruName string
	param  string
	format string
	err    string
}

type (
	HelpT struct {
		evt  HelpEvent
		lst  HelpList
		prch HelpPurch
		ext  HelpExit
	}
	HelpEvent struct {
		add string
		dlt string
		shw string
	}
	HelpList struct {
		add      string
		dlt      string
		shw      string
		shwLists string
	}
	HelpPurch struct {
		add string
		dlt string
	}
	HelpExit struct {
		ext string
	}
)

var (
	fieldId          = Field{"id", "id", "string", `^[1-9][0-9]*$`, "Ошибка: id должен быть натуральным числом"}
	fieldName        = Field{"name", "Название", "string", `^.{1,}$`, "Ошибка: название должно содержать хотя бы один символ"}
	fieldDescription = Field{"description", "Описание", "string", `^.{0,}$`, ""}
	fieldDuration    = Field{"duration", "Продолжительность", "time", "15:04", "Ошибка: продолжтельность должна быть в формате hh:mm"}
	fieldOClock      = Field{"oClock", "Время", "time", "15:04", "Ошибка: время должно быть в формате hh:mm"}
	fieldDate        = Field{"date", "Дата", "time", "02.01.2006", "Ошибка: дата должна быть в формате DD:MM:YYYY"}
	fieldDecide      = Field{"decide", "Подтвердите удаление (y/n)", "decide", "y n", "Ошибка: решение должно быть y/n"}
	fieldListName    = Field{"listName", "Название списка", "string", `^.{1,}$`, "Ошибка: название списка должно содержать хотя бы один символ"}
)

var help = HelpT{
	HelpEvent{
		"evt add — добавить событие.\n" +
			"После ввода команды будет предложено ввести данные о событии: \n" +
			"\t- название (минимум один любой символ)\n" +
			"\t- дата (в формате DD.MM.YYYY)\n" +
			"\t- время (в формате hh:mm)\n" +
			"\t- продолжительность (в формате hh:mm)\n" +
			"\t- описание.\n\n" +
			"Каждому событию присваивается уникальный id.\n",

		"evt dlt — удалить событие.\n" +
			"После ввода команды будет предложено ввести id события (натуральное число) и подтвердить удаление.\n",

		"evt shw [argument] — показать события.\n" +
			"В зависимости от аргумента далее будет предложено ввести различные данные для поиска:\n" +
			"\t- all — вывод всех событий.\n" +
			"\t- date — будет предложено ввести дату события (в формате DD.MM.YYYY). Вывод событий на конкретную дату.\n" +
			"\t- desc — будет предложено ввести описание события. Вывод событий по описанию.\n" +
			"\t- dur — будет предложено ввести продолжительность события (в формате hh:mm). Вывод событий по продолжительности.\n" +
			"\t- id — будет предложено ввести id события (натуральное число). Вывод события по id.\n" +
			"\t- intv — будет предложено ввести интервал (даты начала и конца интервала в формате DD.MM.YYYY).\n" +
			"\tВывод событий в заданном интервале (включая крайние даты).\n" +
			"\t- name — будет предложено ввести название события (минимум один любой символ). Вывод событий по имени.\n" +
			"\t- time — будет предложено ввести дату и время события (в формате DD:MM:YYYY и hh:mm соответственно).\n" +
			"\tВывод событий по дате и времени (выводятся события, которые активны, то есть учитывает продолжительность).\n"},

	HelpList{
		"lst add — добавить список.\n" +
			"После ввода будет предложено ввести название нового списка (минимум один любой символ).\n",

		"lst dlt — удалить список.\n" +
			"После ввода команды будет предложено ввести название списка (минимум один любой символ) и подтвердить удаление.\n",

		"lst shw — показать список.\n" +
			"После ввода команды будет предложено ввести название списка (минимум один любой символ). Вывод всего списка.\n",

		"lsts shw — показать все списки.\n" +
			"Вывод всех списков.\n"},

	HelpPurch{
		"prch add — добавить заметку.\n" +
			"После ввода команды будет предложено ввести данные о заметке:\n" +
			"\t- название (минимум один любой символ)\n" +
			"\t- описание.\n\n" +
			"Каждой заметке присваивается уникальный id.\n",

		"prch dlt — удалить заметку.\n" +
			"После ввода команды будет предложено ввести данные о заметке и подтвердить удаление:\n" +
			"\t- название списка (минимум один любой символ)\n" +
			"\t- id заметки (натуральное число).\n"},

	HelpExit{
		"ext — выход из программы.\n"}}

func AddEvent() {
	var (
		event        Event
		flag         bool
		oClock, date string
	)

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	fmt.Println("Введите данные")

	if event.name, flag = MyScan(fieldName); flag == false {
		return
	}

	if date, flag = MyScan(fieldDate); flag == false {
		return
	}
	if oClock, flag = MyScan(fieldOClock); flag == false {
		return
	}
	event.myTime = TimeToSQL(date, oClock)

	if event.duration, flag = MyScan(fieldDuration); flag == false {
		return
	}
	event.duration += ":00"

	if event.description, flag = MyScan(fieldDescription); flag == false {
		return
	}

	_, err = db.Exec("INSERT INTO `table_of_events`(`time`, `duration`, `name`, `description`) VALUES (?, ?, ?, ?)", event.myTime, event.duration, event.name, event.description)
	if !CheckErr(err, "Ошибка добавления в БД") {
		return
	}
}

func ShowEvent(param string) {
	var (
		event Event
		value string
		flag  bool
	)

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	fmt.Println("Введите данные")

	switch param {
	case "name":
		if value, flag = MyScan(fieldName); flag == false {
			return
		}

		res, err := db.Query("SELECT * FROM `table_of_events` WHERE name = ?", value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRowsEvent(res) {
			return
		}

		return
	case "desc":
		if value, flag = MyScan(fieldDescription); flag == false {
			return
		}

		res, err := db.Query("SELECT * FROM `table_of_events` WHERE `description` = ?", value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRowsEvent(res) {
			return
		}

		return
	case "date":
		if value, flag = MyScan(fieldDate); flag == false {
			return
		}
		sliceValue := strings.Split(value, ".")
		value = sliceValue[2] + "-" + sliceValue[1] + "-" + sliceValue[0]

		res, err := db.Query("SELECT * FROM `table_of_events` WHERE DATE_FORMAT(time, '%Y-%m-%d') = ? ORDER BY `time`", value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRowsEvent(res) {
			return
		}

		return
	case "dur":
		if value, flag = MyScan(fieldDuration); flag == false {
			return
		}
		value += ":00"

		res, err := db.Query("SELECT * FROM `table_of_events` WHERE `duration` = ?", value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRowsEvent(res) {
			return
		}

		return
	case "time":
		var date, oClock string
		if date, flag = MyScan(fieldDate); flag == false {
			return
		}
		if oClock, flag = MyScan(fieldOClock); flag == false {
			return
		}
		value = TimeToSQL(date, oClock)

		res, err := db.Query("SELECT * FROM `table_of_events` WHERE `time` <= ? AND ? <= `time` + INTERVAL `duration` HOUR_SECOND ORDER BY time", value, value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRowsEvent(res) {
			return
		}

		return
	case "id":
		if value, flag = MyScan(fieldId); flag == false {
			return
		}
		event.id, _ = strconv.Atoi(value)

		res, err := db.Query("SELECT * FROM `table_of_events` WHERE `id` = ?", event.id)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRowsEvent(res) {
			return
		}

		return
	case "intv":
		var date1, date2 string

		fmt.Println("Даты начала и конца интервала:")

		if date1, flag = MyScan(fieldDate); flag == false {
			return
		}
		if date2, flag = MyScan(fieldDate); flag == false {
			return
		}

		time1, err := time.Parse("02.01.2006", date1)
		if !CheckErr(err, "Ошибка функции time.Parse:") {
			return
		}
		time2, err := time.Parse("02.01.2006", date2)
		if !CheckErr(err, "Ошибка функции time.Parse:") {
			return
		}

		if time1.After(time2) {
			fmt.Println("Ошибка: дата начала интервала не может быть после даты конца")
			return
		}

		sliceValue := strings.Split(date1, ".")
		date1 = sliceValue[2] + "-" + sliceValue[1] + "-" + sliceValue[0] + " 00:00:00"
		sliceValue = strings.Split(date2, ".")
		date2 = sliceValue[2] + "-" + sliceValue[1] + "-" + sliceValue[0] + " 23:59:59"

		res, err := db.Query("SELECT * FROM `table_of_events` WHERE ? <= `time` AND `time` <= ? ORDER BY `time`", date1, date2)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRowsEvent(res) {
			return
		}

		return
	case "all":
		res, err := db.Query("SELECT * FROM `table_of_events`")
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRowsEvent(res) {
			return
		}

		return
	}
}

func DeleteEvent() {
	var (
		event  Event
		value  string
		flag   bool
		decide string
	)

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	fmt.Println("Введите данные")

	if value, flag = MyScan(fieldId); flag == false {
		return
	}
	event.id, _ = strconv.Atoi(value)

	res, err := db.Query("SELECT * FROM `table_of_events` WHERE `id` = ?", event.id)
	if !CheckErr(err, "Ошибка считывания из БД") {
		return
	}

	if !PrintRowsEvent(res) {
		return
	}

	if decide, flag = MyScan(fieldDecide); flag == false {
		return
	}

	if decide == "y" {
		_, err = db.Exec("DELETE FROM `table_of_events` WHERE `id` = ?", event.id)
		if !CheckErr(err, "Ошибка удаления из БД") {
			return
		}
	}
}

func AddList() {
	var (
		name string
		flag bool
	)

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	fmt.Println("Введите данные")

	if name, flag = MyScan(fieldName); flag == false {
		return
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE `%s` (`id` INT(11) NOT NULL AUTO_INCREMENT, `name` VARCHAR(100) NOT NULL, `description` VARCHAR(100) NULL, PRIMARY KEY (`id`))", name))
	if !CheckErr(err, "Ошибка создания БД") {
		return
	}
}

func ShowList() {
	var (
		name string
		flag bool
	)

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	fmt.Println("Введите данные")

	if name, flag = MyScan(fieldName); flag == false {
		return
	}

	if name == "table_of_events" {
		fmt.Println("Ошибка считывания из БД: нет доступа к данной таблице")
		return
	}

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `%s`", name))
	if !CheckErr(err, "Ошибка считывания из БД") {
		return
	}

	if !PrintRowsList(res) {
		return
	}
}

func ShowLists() {
	var name string

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	res, err := db.Query("SHOW TABLES")
	if !CheckErr(err, "Ошибка считывания из БД") {
		return
	}

	count := 0
	for res.Next() {
		err := res.Scan(&name)
		if !CheckErr(err, "Ошибка функции res.Scan") {
			continue
		}
		if name == "table_of_events" {
			continue
		}
		fmt.Println(name)
		if !CheckErr(err, "Ошибка функции fmt.Fprintf") {
			return
		}
		count++
	}
	if count == 0 {
		fmt.Println("Данные не найдены")
		return
	}
}

func DeleteList() {
	var (
		name   string
		flag   bool
		decide string
	)

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	fmt.Println("Введите данные")

	if name, flag = MyScan(fieldName); flag == false {
		return
	}

	if name == "table_of_events" {
		fmt.Println("Ошибка считывания из БД: нет доступа к данной таблице")
		return
	}

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `%s`", name))
	if !CheckErr(err, "Ошибка считывания из БД") {
		return
	}

	if !PrintRowsList(res) {
		return
	}

	if decide, flag = MyScan(fieldDecide); flag == false {
		return
	}

	if decide == "y" {
		_, err = db.Exec(fmt.Sprintf("DROP TABLE `%s`", name))
		if !CheckErr(err, "Ошибка удаления БД") {
			return
		}
	}
}

func AddPurch() {
	var (
		event    Event
		flag     bool
		listName string
	)

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	fmt.Println("Введите данные")

	if listName, flag = MyScan(fieldListName); flag == false {
		return
	}

	if event.name, flag = MyScan(fieldName); flag == false {
		return
	}

	if event.description, flag = MyScan(fieldDescription); flag == false {
		return
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO `%s`(`name`, `description`) VALUES (?, ?)", listName), event.name, event.description)
	if !CheckErr(err, "Ошибка добавления в БД") {
		return
	}
}

func DeletePurch() {
	var (
		listName string
		event    Event
		value    string
		flag     bool
		decide   string
	)

	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if !CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	fmt.Println("Введите данные")

	if listName, flag = MyScan(fieldListName); flag == false {
		return
	}

	if value, flag = MyScan(fieldId); flag == false {
		return
	}
	event.id, _ = strconv.Atoi(value)

	res, err := db.Query(fmt.Sprintf("SELECT * FROM `%s` WHERE `id` = ?", listName), event.id)
	if !CheckErr(err, "Ошибка считывания из БД") {
		return
	}

	if !PrintRowsList(res) {
		return
	}

	if decide, flag = MyScan(fieldDecide); flag == false {
		return
	}

	if decide == "y" {
		_, err = db.Exec(fmt.Sprintf("DELETE FROM `%s` WHERE `id` = ?", listName), event.id)
		if !CheckErr(err, "Ошибка удаления из БД") {
			return
		}
	}
}

func Help(input string) {
	switch input {
	case "help":
		fmt.Println(help.evt.add)
		fmt.Println(help.evt.dlt)
		fmt.Println(help.evt.shw)

		fmt.Println(help.ext.ext)

		fmt.Println(help.lst.add)
		fmt.Println(help.lst.dlt)
		fmt.Println(help.lst.shw)
		fmt.Println(help.lst.shwLists)

		fmt.Println(help.prch.add)
		fmt.Println(help.prch.dlt)
	case "help evt":
		fmt.Println(help.evt.add)
		fmt.Println(help.evt.dlt)
		fmt.Println(help.evt.shw)
	case "help ext":
		fmt.Println(help.ext.ext)
	case "help lst":
		fmt.Println(help.lst.add)
		fmt.Println(help.lst.dlt)
		fmt.Println(help.lst.shw)
		fmt.Println(help.lst.shwLists)
	case "help prch":
		fmt.Println(help.prch.add)
		fmt.Println(help.prch.dlt)
	case "help evt add":
		fmt.Println(help.evt.add)
	case "help evt dlt":
		fmt.Println(help.evt.dlt)
	case "help evt shw":
		fmt.Println(help.evt.shw)
	case "help lst add":
		fmt.Println(help.lst.add)
	case "help lst dlt":
		fmt.Println(help.lst.dlt)
	case "help lst shw":
		fmt.Println(help.lst.shw)
	case "help lsts shw":
		fmt.Println(help.lst.shwLists)
	case "help prch add":
		fmt.Println(help.prch.add)
	case "help prch dlt":
		fmt.Println(help.prch.dlt)
	default:
		fmt.Printf("Команды '%s' не существует\nДля уточнения информации по командам введите 'help'\n", input)
	}
}

func CheckErr(err error, log string) bool {
	if err != nil {
		fmt.Println(fmt.Errorf("%s: %w", log, err))
		return false
	}
	return true
}

func CheckFormat(str, param, format, log string) bool {
	_, err := time.Parse(format, str)
	if (param == "string" && !regexp.MustCompile(format).MatchString(str)) || (param == "time" && err != nil) {
		fmt.Println(fmt.Errorf("%s", log))
		return false
	}
	if param == "decide" && !StringInSlice(str, strings.Split(format, " ")) {
		fmt.Println(log)
		return false
	}
	return true
}

func MyScan(field Field) (string, bool) {
	in := bufio.NewReader(os.Stdin)
	fmt.Print(fmt.Sprintf("%s: ", field.ruName))
	str, err := in.ReadString('\n')
	str = strings.ReplaceAll(str, "\n", "")
	if !CheckErr(err, "Ошибка чтения in.ReadString") || !CheckFormat(str, field.param, field.format, field.err) {
		return "", false
	}
	return str, true
}

func PrintRowsEvent(res *sql.Rows) bool {
	var event Event

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", "id", "name", "date", "time", "duration", "description")
	if !CheckErr(err, "Ошибка функции fmt.Fprintf") {
		return false
	}

	count := 0
	for res.Next() {
		err := res.Scan(&event.id, &event.myTime, &event.duration, &event.name, &event.description)
		if !CheckErr(err, "Ошибка функции res.Scan") {
			continue
		}
		date, oClock := TimeFromSQL(event.myTime)
		event.duration = strings.Replace(event.duration, ":00", "", 1)
		_, err = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t\n", event.id, event.name, date, oClock, event.duration, event.description)
		if !CheckErr(err, "Ошибка функции fmt.Fprintf") {
			return false
		}
		count++
	}
	if count == 0 {
		fmt.Println("Данные не найдены")
		return false
	}

	err = w.Flush()
	if !CheckErr(err, "Ошибка функции w.Flush") {
		return false
	}

	return true
}

func PrintRowsList(res *sql.Rows) bool {
	var event Event

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t\n", "id", "name", "description")
	if !CheckErr(err, "Ошибка функции fmt.Fprintf") {
		return false
	}

	for res.Next() {
		err := res.Scan(&event.id, &event.name, &event.description)
		if !CheckErr(err, "Ошибка функции res.Scan") {
			continue
		}
		_, err = fmt.Fprintf(w, "%d\t%s\t%s\t\n", event.id, event.name, event.description)
		if !CheckErr(err, "Ошибка функции fmt.Fprintf") {
			return false
		}
	}

	err = w.Flush()
	if !CheckErr(err, "Ошибка функции w.Flush") {
		return false
	}

	return true
}

func TimeFromSQL(myTime string) (string, string) {
	myTimeSlice := strings.Split(myTime, " ")
	date, oClock := myTimeSlice[0], myTimeSlice[1]
	dateSlice := strings.Split(date, "-")
	date = dateSlice[2] + "." + dateSlice[1] + "." + dateSlice[0]
	oClock = strings.Replace(oClock, ":00", "", 1)
	return date, oClock
}

func TimeToSQL(date, oClock string) string {
	dateSlice := strings.Split(date, ".")
	myTime := dateSlice[2] + "-" + dateSlice[1] + "-" + dateSlice[0]
	myTime += " " + oClock + ":00"
	return myTime
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
