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
	records, err := getCurrentBase(base)
	if err != nil {
		fmt.Println("Error with getting Base")
		return -1, fmt.Errorf("Not available")
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

func AddCodeInCSV[T MutableBase](base T, text_leksem string) error {
	code, _ := FindCodeInCSV(base, text_leksem)
	if code > -1 {
		return fmt.Errorf("code and leksem is exist in base!")
	}

	file, err := os.Create(endpoints[base])

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	value_new_code_leksem, _ := calculateLengthBase(base)
	value_new_code_leksem++

	data := []string{text_leksem, strconv.Itoa(value_new_code_leksem)}

	err = writer.Write(data)
	if err != nil {
		fmt.Println("Error: ", err)
		return nil
	}
	writer.Flush()
	return nil
}

func unespace(text_leksem string) string {
	var unescape_text_leksem string
	unescape_text_leksem = strings.ReplaceAll(text_leksem, "`", "")
	return unescape_text_leksem

}

func calculateLengthBase[T MutableBase | ImmutableBase](base T) (int, error) {
	file, err := os.Open(endpoints[base])

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		return 0, fmt.Errorf("Error reading records")

	}

	return len(records), nil
}

func getCurrentBase[T MutableBase | ImmutableBase](base T) ([][]string, error) {
	file, err := os.Open(endpoints[base])

	if err != nil {
		log.Fatal("Error while reading the file", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()

	if err != nil {
		fmt.Println("Error reading records")
		err = fmt.Errorf("Error reading records")
		return [][]string{}, err
	}

	return records, nil
}
