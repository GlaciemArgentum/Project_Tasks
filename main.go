package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

func main() {
	db, err := sql.Open("mysql", "root:12345678@tcp(127.0.0.1:3306)/db")
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			fmt.Printf("Ошибка закрытия БД: %v\n", err)
		}
	}(db)
	if CheckErr(err, "Ошибка открытия БД") {
		panic("panic")
	}

	in := bufio.NewReader(os.Stdin)
	var input string

	for {
		input, err = in.ReadString('\n')
		input = strings.ReplaceAll(input, "\n", "")
		if CheckErr(err, "Ошибка чтения in.ReadString") {
			continue
		}

		switch input {
		case `/addEvent`:
			var event Event
			fmt.Print("Введите данные\nНазвание: ")
			event.name, err = in.ReadString('\n')
			event.name = strings.ReplaceAll(event.name, "\n", "")
			if CheckErr(err, "Ошибка чтения in.ReadString") || CheckLogic(event.name, `^.{1,}$`, "Название должно содиржать хотя бы один символ") {
				break
			}

			fmt.Print("Дата: ")
			event.dataTime, err = in.ReadString('\n')
			event.dataTime = strings.ReplaceAll(event.dataTime, "\n", "")
			if CheckErr(err, "Ошибка чтения in.ReadString") || CheckDataFormat(event.dataTime, "02.01.2006", "Неверный формат даты") {
				break
			}
			dataSlice := strings.Split(event.dataTime, ".")
			event.dataTime = dataSlice[2] + "-" + dataSlice[1] + "-" + dataSlice[0]

			fmt.Print("Время: ")
			time, err := in.ReadString('\n')
			time = strings.ReplaceAll(time, "\n", "")
			if CheckErr(err, "Ошибка чтения in.ReadString") || CheckDataFormat(time, "15:04", "Неверный формат времени") {
				break
			}
			event.dataTime += " " + time + ":00"

			fmt.Print("Продолжительность: ")
			event.duration, err = in.ReadString('\n')
			event.duration = strings.ReplaceAll(event.duration, "\n", "")
			if CheckErr(err, "Ошибка чтения in.ReadString") || CheckDataFormat(event.duration, "15:04", "Неверный формат продолжительности") {
				break
			}
			event.duration += ":00"

			fmt.Print("Описание: ")
			event.description, err = in.ReadString('\n')
			event.description = strings.ReplaceAll(event.description, "\n", "")
			if CheckErr(err, "Ошибка чтения in.ReadString") {
				break
			}

			_, err = db.Exec(fmt.Sprintf("INSERT INTO table_of_events(time_begin, duration, name, description) VALUES ('%s', '%s', '%s', '%s')", event.dataTime, event.duration, event.name, event.description))
			if CheckErr(err, "Ошибка добавления данных в БД") {
				break
			}

		case `/findEvent -id`:
			var id int
			_, err := fmt.Sscanf(input, "/findEvent id %d", id)
			CheckErr(err, "Ошибка Sscanf")

		}
	}
}

// /addEvent Встреча 20.09.2000 02:30 Описание

// /showEvent -id 1
// /showEvent -dataTime 2023-12-23T05:50:00
// /showEvent -interval 2023-12-23T05:50:00 2023-12-23T05:55:00

// /deleteEvent -id 1
// /deleteEvent -dataTime 2023-12-23T05:50:00

// /addList Список покупок

// /addPurchase Список_покупок Покупка

// /showList Список покупок

// /deleteList Список покупок

// /deletePurchase Список_покупок Покупка

// /help

// /finish
