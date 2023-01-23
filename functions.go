package main

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

type Event struct {
	name        string
	dataTime    string
	duration    string
	description string
}

func CheckErr(err error, log string) {
	if err != nil {
		fmt.Println(fmt.Errorf("%s: %w\n", log, err))
	}
}
