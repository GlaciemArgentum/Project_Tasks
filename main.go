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
	//defer CheckErr(db.Close(), "Ошибка закрытия БД")
	CheckErr(err, "Ошибка открытия БД")

	in := bufio.NewReader(os.Stdin)
	var input string
	for {
		input, _ = in.ReadString('\n')
		switch {

		case strings.HasPrefix(input, "/addEvent "):
			var event Event
			_, err := fmt.Sscanf(input, "/addEvent %s %s %s %s", &event.name, &event.dataTime, &event.duration, &event.description)
			CheckErr(err, "Ошибка Sscanf")
			event.dataTime = strings.Replace(event.dataTime, "T", " ", 1)
			_, err = db.Exec(fmt.Sprintf("INSERT INTO table_of_events(time_begin, duration, name, description) VALUES ('%s', '%s', '%s', '%s')", event.dataTime, event.duration, event.name, event.description))
			CheckErr(err, "Ошибка добавления данных в БД")

		case strings.HasPrefix(input, "/findEvent id "):
			var id int
			_, err := fmt.Sscanf(input, "/findEvent id %d", id)
			CheckErr(err, "Ошибка Sscanf")

		}
	}
}

// /addEvent Встреча 2023-12-23T05:50:00 02:30:00 Описание

// /showEvent id 1
// /showEvent dataTime 2023-12-23T05:50:00
// /showEvent interval 2023-12-23T05:50:00 2023-12-23T05:55:00

// /deleteEvent id 1
// /deleteEvent dataTime 2023-12-23T05:50:00

// /addList Список покупок

// /addPurchase Список_покупок Покупка

// /showList Список покупок

// /deleteList Список покупок

// /deletePurchase Список_покупок Покупка

// /help

// /finish
