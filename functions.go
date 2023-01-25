package main

import (
	"bufio"
	"database/sql"
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"regexp"
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
	myType string
	format string
	err    string
}

var (
	id          = Field{"id", "id", "string", `^[1-9][0-9]*$`, "id должно быть натуральным числом"}
	name        = Field{"name", "Имя", "string", `^.{1,}$`, "Название должно содержать хотя бы один символ"}
	description = Field{"description", "Описание", "string", `^.{0,}$`, "Название должно содержать хотя бы один символ"}
	duration    = Field{"duration", "Продолжительность", "time", `^.{1,}$`, "Название должно содержать хотя бы один символ"}
	oClock      = Field{"oClock", "Время", "time", `^.{1,}$`, "Название должно содержать хотя бы один символ"}
	date        = Field{"date", "Дата", "time", `^.{1,}$`, "Название должно содержать хотя бы один символ"}
)

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
	return true
}

func MyScan(log, param, format, errLogic string) (string, bool) {
	in := bufio.NewReader(os.Stdin)
	fmt.Print(log)
	str, err := in.ReadString('\n')
	str = strings.ReplaceAll(str, "\n", "")
	if !CheckErr(err, "Ошибка чтения in.ReadString") || !CheckFormat(str, param, format, errLogic) {
		return "", false
	}
	return str, true
}

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

	if event.name, flag = MyScan("Введите данные\nНазвание: ", "string", `^.{1,}$`, "Название должно содержать хотя бы один символ"); flag == false {
		return
	}

	if date, flag = MyScan("Дата: ", "time", "02.01.2006", "Неверный формат даты"); flag == false {
		return
	}
	if oClock, flag = MyScan("Время: ", "time", "15:04", "Неверный формат времени"); flag == false {
		return
	}
	event.myTime = TimeToSQL(date, oClock)

	if event.duration, flag = MyScan("Продолжительность: ", "time", "15:04", "Неверный формат продолжительности"); flag == false {
		return
	}
	event.duration += ":00"

	if event.description, flag = MyScan("Описание: ", "string", `^.{0,}$`, " "); flag == false {
		return
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO table_of_events(time, duration, name, description) VALUES ('%s', '%s', '%s', '%s')", event.myTime, event.duration, event.name, event.description))
	if !CheckErr(err, "Ошибка добавления в БД") {
		return
	}
}

func FindEvent(param string) {
	var (
		event Event
		field string
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

	switch param {
	case "name":
		if field, flag = MyScan("Введите данные\nИмя: ", "string", `^.{0,}$`, ""); flag == false {
			return
		}
	case "description":
		if field, flag = MyScan("Введите данные\nОписание: ", "string", `^.{0,}$`, ""); flag == false {
			return
		}
	case "id":
		if field, flag = MyScan("Введите данные\nid: ", "string", `^[1-9][0-9]*$`, "id должно быть натуральным числом"); flag == false {
			return
		}
	case "date":
		if field, flag = MyScan("Введите данные\nДата: ", "time", "02.01.2006", "Неверный формат даты"); flag == false {
			return
		}
	case "duration":
		if field, flag = MyScan("Введите данные\nПродолжительность: ", "time", "15:04", "Неверный формат продолжительности"); flag == false {
			return
		}
	case "time":
		var date, oClock string
		if date, flag = MyScan("Введите данные\nДата: ", "time", "02.01.2006", "Неверный формат даты"); flag == false {
			return
		}
		if oClock, flag = MyScan("Введите данные\nВремя: ", "time", "15:04", "Неверный формат времени"); flag == false {
			return
		}
		field = TimeToSQL(date, oClock)
	}

	/*
		if field, flag = MyScan("Введите данные\nfield: ", "string", `^[1-9][0-9]*$`, "field должно быть натуральным числом"); flag == false {
			return
		}
		event.id, _ = strconv.Atoi(field)

		res, err := db.Query(fmt.Sprintf("SELECT * FROM table_of_events WHERE field = %d", event.id)) */
	if !CheckErr(err, "Ошибка считывания из БД") {
		return
	}

	if !PrintRows(res) {
		return
	}
}

func PrintRows(res *sql.Rows) bool {
	var event Event

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 3, ' ', 0)

	_, err := fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\t\n", "id", "name", "date", "time", "duration", "description")
	if !CheckErr(err, "Ошибка fmt.Fprintf") {
		return false
	}

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
