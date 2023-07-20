package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	"golang.org/x/text/encoding/charmap"
	"io"
	"log"
	"os"
	"time"
)

type ValCurse struct {
	XMLName xml.Name `xml:"ValCurs"`
	Valute  []struct {
		NumCode  int    `xml:"NumCode"`
		CharCode string `xml:"CharCode"`
		Nominal  int    `xml:"Nominal"`
		Name     string `xml:"Name"`
		Value    string `xml:"Value"`
	} `xml:"Valute"`
}

// Преобразование даты из YYYY-MM-DD в DD/MM/YYYY
func dateFormatter(date string) string {
	layout := "2006-01-02"
	desiredLayout := "02/01/2006"
	data, err := time.Parse(layout, date)
	if err != nil {
		fmt.Println("Произошла ошибка при разборе даты:", err)
		return ""
	}

	// Преобразование объекта time.Time в нужный формат строки
	formattedDate := data.Format(desiredLayout)
	return formattedDate
}

/*
/*

Вариант, который должен был работать, но при запросе через http.Get в теле оказывается
Client does not have access rights to the content so server is rejecting to give proper response
*/
/*
func getXMLFromCBR(URL string) []byte {
	// Отправляем GET-запрос к API
	response, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}

	// Читаем ответ от сервера
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	//fmt.Println(string(body))
	return body
}
*/

// Парсинг ответа от сервера с помощью библиотеки geziyor
func getXMLFromCBR(URL string) {
	geziyor.NewGeziyor(&geziyor.Options{
		StartURLs:   []string{URL},
		ParseFunc:   parseFunc,
		Exporters:   []export.Exporter{},
		LogDisabled: true,
	}).Start()
}

func parseFunc(g *geziyor.Geziyor, r *client.Response) {

	file, err := os.Create("file.txt")

	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(string(r.Body))
}

func main() {

	code := flag.String("code", "USD", "Код валюты")
	date := flag.String("date", "2022-10-08", "Дата")
	flag.Parse()

	// Проверка введённой даты
	dt1, err := time.Parse("2006-01-02", *date)
	if !dt1.Before(time.Now()) {
		fmt.Println("Введённая дата превышает настоящее число")
		return
	}

	formattedDate := dateFormatter(*date)

	// Составление URL
	baseURL := "http://www.cbr.ru/scripts/XML_daily.asp?"
	data_req := "date_req=" + formattedDate
	URL := baseURL + data_req

	getXMLFromCBR(URL)

	// Открытие файла и декод из windows-1251
	f, err := os.Open("file.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	d := xml.NewDecoder(f)
	d.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	valCourse := ValCurse{}
	err = d.Decode(&valCourse)
	if err != nil {
		log.Fatalf("decode: %v", err)
	}

	// Поиск данных в xml
	flag1 := false
	for _, valNode := range valCourse.Valute {
		if valNode.CharCode == *code {
			fmt.Printf("%s (%s): %s", valNode.CharCode, valNode.Name, valNode.Value)
			flag1 = true
		}
	}
	if !flag1 {
		fmt.Printf("Введённый вами код: %s - не найден", *code)
	}
}
