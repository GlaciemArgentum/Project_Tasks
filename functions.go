package main

import (
	_ "database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"regexp"
	"time"
)

type Event struct {
	name        string
	dataTime    string
	duration    string
	description string
}

func CheckErr(err error, log string) bool {
	if err != nil {
		fmt.Println(fmt.Errorf("%s: %w", log, err))
		return true
	}
	return false
}

func CheckDataFormat(date, format string, log string) bool {
	_, err := time.Parse(format, date)
	if err != nil {
		fmt.Println(fmt.Errorf("%s", log))
		return true
	}
	return false
}

func CheckLogic(str string, sample string, log string) bool {
	if !regexp.MustCompile(sample).MatchString(str) {
		fmt.Println(fmt.Errorf("%s", log))
		return true
	}
	return false
}
