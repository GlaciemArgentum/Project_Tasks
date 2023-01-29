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

var (
	fieldId          = Field{"id", "id", "string", `^[1-9][0-9]*$`, "id должен быть натуральным числом"}
	fieldName        = Field{"name", "Имя", "string", `^.{1,}$`, "Название должно содержать хотя бы один символ"}
	fieldDescription = Field{"description", "Описание", "string", `^.{0,}$`, ""}
	fieldDuration    = Field{"duration", "Продолжительность", "time", "15:04", "Неверный формат продолжительности"}
	fieldOClock      = Field{"oClock", "Время", "time", "15:04", "Неверный формат времени"}
	fieldDate        = Field{"date", "Дата", "time", "02.01.2006", "Неверный формат даты"}
	fieldDecide      = Field{"decide", "Вы уверены, что хотите удалить это событие (y/n)", "decide", "y n", "Неверный формат ответа: должен быть y/n"}
)

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

	_, err = db.Exec("INSERT INTO table_of_events(time, duration, name, description) VALUES (?, ?, ?, ?)", event.myTime, event.duration, event.name, event.description)
	if !CheckErr(err, "Ошибка добавления в БД") {
		return
	}
}

func FindEvent(param string) {
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

		res, err := db.Query("SELECT * FROM table_of_events WHERE name = ?", value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRows(res) {
			return
		}

		return
	case "description":
		if value, flag = MyScan(fieldDescription); flag == false {
			return
		}

		res, err := db.Query("SELECT * FROM table_of_events WHERE description = ?", value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRows(res) {
			return
		}

		return
	case "date":
		if value, flag = MyScan(fieldDate); flag == false {
			return
		}
		sliceValue := strings.Split(value, ".")
		value = sliceValue[2] + "-" + sliceValue[1] + "-" + sliceValue[0]

		res, err := db.Query("SELECT * FROM table_of_events WHERE DATE_FORMAT(time, '%Y-%m-%d') = ? ORDER BY time", value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRows(res) {
			return
		}

		return
	case "duration":
		if value, flag = MyScan(fieldDuration); flag == false {
			return
		}
		value += ":00"

		res, err := db.Query("SELECT * FROM table_of_events WHERE duration = ?", value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRows(res) {
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

		res, err := db.Query("SELECT * FROM table_of_events WHERE time <= ? AND ? <= time + INTERVAL duration HOUR_SECOND ORDER BY time", value, value)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRows(res) {
			return
		}

		return
	case "id":
		if value, flag = MyScan(fieldId); flag == false {
			return
		}
		event.id, _ = strconv.Atoi(value)

		res, err := db.Query("SELECT * FROM table_of_events WHERE id = ?", event.id)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRows(res) {
			return
		}

		return
	case "interval":
		var date1, date2 string

		fmt.Println("Даты начала и конца интервала:")

		if date1, flag = MyScan(fieldDate); flag == false {
			return
		}
		if date2, flag = MyScan(fieldDate); flag == false {
			return
		}

		time1, err := time.Parse("02.01.2006", date1)
		if !CheckErr(err, "Ошибка time.Parse:") {
			return
		}
		time2, err := time.Parse("02.01.2006", date2)
		if !CheckErr(err, "Ошибка time.Parse:") {
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

		res, err := db.Query("SELECT * FROM table_of_events WHERE ? <= time AND time <= ? ORDER BY time", date1, date2)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRows(res) {
			return
		}

		return
	case "all":
		res, err := db.Query("SELECT * FROM table_of_events")
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}

		if !PrintRows(res) {
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

	res, err := db.Query("SELECT * FROM table_of_events WHERE id = ?", event.id)
	if !CheckErr(err, "Ошибка считывания из БД") {
		return
	}

	if !PrintRows(res) {
		return
	}

	if decide, flag = MyScan(fieldDecide); flag == false {
		return
	}

	if decide == "y" {
		_, err = db.Exec("DELETE FROM table_of_events WHERE id = ?", event.id)
		if !CheckErr(err, "Ошибка считывания из БД") {
			return
		}
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

func PrintRows(res *sql.Rows) bool {
	var event Event

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", "id", "name", "date", "time", "duration", "description")
	if !CheckErr(err, "Ошибка fmt.Fprintf") {
		return false
	}

	count := 0
	for res.Next() {
		err := res.Scan(&event.id, &event.myTime, &event.duration, &event.name, &event.description)
		if !CheckErr(err, "Ошибка res.Scan") {
			continue
		}
		date, oClock := TimeFromSQL(event.myTime)
		event.duration = strings.Replace(event.duration, ":00", "", 1)
		_, err = fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\t\n", event.id, event.name, date, oClock, event.duration, event.description)
		if !CheckErr(err, "Ошибка fmt.Fprintf") {
			return false
		}
		count++
	}
	if count == 0 {
		fmt.Println("Данные не найдены")
		return false
	}

	err = w.Flush()
	if !CheckErr(err, "Ошибка w.Flush") {
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
