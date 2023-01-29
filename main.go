package main

import (
	"bufio"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("\nВведите команду: ")
		input, err := in.ReadString('\n')
		input = strings.ReplaceAll(input, "\n", "")
		if !CheckErr(err, "Ошибка чтения in.ReadString") {
			continue
		}

		switch {
		case input == "/addEvent":
			AddEvent()
		case strings.HasPrefix(input, "/findEvent"):
			param := strings.Replace(input, "/findEvent -", "", 1)
			if !StringInSlice(param, []string{"id", "name", "date", "time", "duration", "description", "interval", "all"}) {
				fmt.Println("Ошибка. Такого параметра нет")
				break
			}
			FindEvent(param)
		case input == "/deleteEvent":
			DeleteEvent()
		}
	}
}

// /addEvent Встреча 20.09.2000 02:30 Описание

// /findEvent -id 1
// /findEvent -myTime 2023-12-23T05:50:00
// /findEvent -interval 2023-12-23T05:50:00 2023-12-23T05:55:00

// /deleteEvent -id 1

// /addList Список покупок

// /addPurchase Список_покупок Покупка

// /showList Список покупок

// /deleteList Список покупок

// /deletePurchase Список_покупок Покупка

// /help

// /finish
