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
		case input == "evt add":
			AddEvent()
		case strings.HasPrefix(input, "evt shw"):
			param := strings.Replace(input, "evt shw ", "", 1)
			if !StringInSlice(param, []string{"id", "name", "date", "time", "dur", "desc", "intv", "all"}) {
				fmt.Println("Ошибка. Такого параметра нет")
				break
			}
			ShowEvent(param)
		case input == "evt dlt":
			DeleteEvent()
		case input == "lst add":
			AddList()
		case input == "lst shw":
			ShowList()
		case input == "lsts shw":
			ShowLists()
		case input == "lst dlt":
			DeleteList()
		case input == "prch add":
			AddPurch()
		case input == "prch dlt":
			DeletePurch()
		case input == "ext":
			fmt.Println("Программа завершена")
			return
		//case strings.HasPrefix(input, "help"):
		default:
			fmt.Printf("Команды '%s' не существует\nДля уточнения информации по командам введите 'help'\n", input)
		}
	}
}

// добавить редактирование
// сделать не полное совпадение в shw, а частичное
// Решить проблему при смене языка
