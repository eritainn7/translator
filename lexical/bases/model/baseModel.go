package model

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func FindCodeInCSV[T MutableBase | ImmutableBase](base T, text_leksem string) (int, error) {
	file, err := os.Open(endpoints[base])

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error reading records")
	}

	for _, record := range records {
		code_leksem := unespace(record[0])
		if code_leksem == text_leksem {
			code_leksem, err := strconv.Atoi(record[1])
			return code_leksem, err
		}
	}

	return -1, fmt.Errorf("Not available")
}

func AddCodeInCSV[T MutableBase](base T) error {

	return nil
}

func unespace(text_leksem string) string {
	var unescape_text_leksem string
	unescape_text_leksem = strings.ReplaceAll(text_leksem, "`", "")
	return unescape_text_leksem

}
