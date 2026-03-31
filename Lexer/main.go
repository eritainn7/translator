package main

import (
	"bufio"
	"fmt"
	"lexer/lexer"
	"lexer/shared"
	"lexer/tables"
	"os"
	"path/filepath"
)

func main() {
	// Проверка аргументов командной строки
	if len(os.Args) < 2 {
		fmt.Println("Usage: lexer <filename>")
		fmt.Println("Example: lexer program.cs")
		return
	}

	// Определяем путь к данным
	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("Error getting executable path: %v\n", err)
		return
	}

	dataPath := filepath.Join(filepath.Dir(execPath), "data")

	// Создаем папку data если её нет
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		fmt.Printf("Error creating data directory: %v\n", err)
		return
	}

	// Создаем папку output если её нет
	if err := os.MkdirAll("output", 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// Загружаем таблицы из CSV
	tbl, err := tables.NewTables(dataPath)
	if err != nil {
		fmt.Printf("Error loading tables: %v\n", err)
		return
	}

	// Открытие файла с исходным кодом
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error opening file: %v\n", err)
		return
	}
	defer file.Close()

	// Чтение содержимого файла
	var input string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	// Создание и запуск лексического анализатора
	fmt.Println("Результат лексического анализа:")
	fmt.Println("------------------------")

	lexer := lexer.NewLexer(input, tbl)
	list_lexems := lexer.Analyze()

	// Вывод пустой строки после лексем
	fmt.Println()

	// Вывод информации об идентификаторах
	fmt.Println("\nТаблица идентификаторов:")
	fmt.Println("------------------------")
	fmt.Printf("%-5s %-20s %-20s %-15s %-6s %-6s\n",
		"ID", "Идентификатор", "Область видимости", "Тип", "Линия", "Колонка")
	fmt.Println("---------------------------------------------------------------")

	for name, id := range tbl.Identifiers {
		if info, exists := tbl.IdentifierInfo[name]; exists {
			fmt.Printf("%-5d %-20s %-20s %-15s %-6d %-6d\n",
				id, name, info.Scope, info.Type, info.Line, info.Column)
		} else {
			fmt.Printf("%-5d %-20s %-20s %-15s %-6d %-6d\n",
				id, name, "unknown", "unknown", 0, 0)
		}
	}

	// Вывод последовательности лексем в файл
	file_output, err := os.Create("output/lexems.txt")
	if err != nil {
		fmt.Println("Unable to create file:", err)
		os.Exit(1)
	}
	defer file_output.Close()

	var output_string string
	for _, lexem := range list_lexems {
		output_string += lexem.String() + " "
	}
	file_output.Write([]byte(output_string))

	// Формирование таблицы кодов лексем
	shared.SaveLexemesToCSV("output/seq_codes.csv", list_lexems)

	fmt.Println("\n\nРезультаты сохранены в папке output/")
	fmt.Printf("Таблица идентификаторов сохранена в %s/identifiers.csv\n", dataPath)
}
