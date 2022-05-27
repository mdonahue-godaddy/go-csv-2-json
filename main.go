package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"unicode"

	"github.com/kugoucode/golib/gostruct"
)

func main() {
	log.SetOutput(os.Stdout)

	csvFileName := "example.csv"

	csvFile, cerr := os.Open(csvFileName)
	if cerr != nil {
		log.Panicf("os.Open('%s') error: %v", csvFileName, cerr)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	header, herr := reader.Read()
	if herr != nil {
		log.Panicf("Header reader.Read() error: %v", herr)
	}

	if len(header) <= 0 {
		log.Panic("No header columns found")
	}

	builder := gostruct.New()

	for _, column := range header {
		fieldName := gostruct.NewStringField(column)
		//fieldTag := fmt.Sprintf("json:\"%s\" validate:\"required\"", column)
		//fieldName = fieldName.SetTag("`" + fieldTag + "`")
		builder.AddField(*fieldName)
	}

	build := builder.Build()

	var records []interface{}

	for {
		row, rerr := reader.Read()

		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			log.Panicf("Row reader.Read() error: %v", rerr)
		}

		b := build.New()
		for idx, column := range row {
			b.SetString(header[idx], column)
		}

		records = append(records, b.Addr())

		fmt.Printf("Interface %T: %+v\n", b.Interface(), b.Interface())
		fmt.Printf("Addr %T: %+v\n", b.Addr(), b.Addr())
	}

	jsonBytes, jerr := json.MarshalIndent(records, "", "  ")
	if jerr != nil {
		log.Panicf("json.MarshalIndent(records, '', '') error: %v", jerr)
	}

	jsonFileName := fmt.Sprintf("%s.json", csvFileName)
	if ferr := ioutil.WriteFile(jsonFileName, jsonBytes, 0644); ferr != nil {
		log.Panicf("ioutil.WriteFile(jsonFileName, jsonBytes, 0644) error: %v", ferr)
	}

	log.Println(string(jsonBytes))
}

func IsUpper(s string) bool {
	for _, r := range s {
		if !unicode.IsUpper(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}
