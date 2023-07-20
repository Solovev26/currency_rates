package main

import (
	"flag"
	"fmt"
)

func main() {

	code := flag.String("code", "USD", "Код валюты")
	date := flag.String("date", "2022-10-08", "Дата")

	flag.Parse()

	fmt.Println(*code)
	fmt.Println(*date)
}
