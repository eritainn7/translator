package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"rpn-converter-full/converter"
	"rpn-converter-full/lexer"
	"strings"
)

func main() {
	var (
		inputFile  = flag.String("input", "", "Входной файл с кодом C#")
		showTokens = flag.Bool("tokens", false, "Показать токены")
		verbose    = flag.Bool("verbose", false, "Подробный вывод")
	)
	flag.Parse()

	var code string

	if *inputFile != "" {
		data, err := os.ReadFile(*inputFile)
		if err != nil {
			fmt.Printf("Ошибка чтения файла: %v\n", err)
			os.Exit(1)
		}
		code = string(data)
	} else {
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			code = strings.Join(lines, "\n")
		} else {
			fmt.Println("Введите код C# (Ctrl+D для завершения):")
			scanner := bufio.NewScanner(os.Stdin)
			var lines []string
			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			code = strings.Join(lines, "\n")
		}
	}

	if code == "" {
		fmt.Println("Нет входных данных")
		flag.Usage()
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("=== Лексический анализ ===")
	}

	lex := lexer.NewLexer(code)
	tokens := lex.Tokenize()

	if *showTokens {
		fmt.Println("\n=== Токены ===")
		for i, tok := range tokens {
			if tok.Type != lexer.TokenEOF {
				fmt.Printf("%3d: %-15s %q (line %d, col %d)\n",
					i, tok.Type.String(), tok.Value, tok.Line, tok.Col)
			}
		}
		fmt.Println()
	}

	if *verbose {
		fmt.Printf("Всего токенов: %d\n", len(tokens))
		fmt.Println("\n=== Преобразование в ОПЗ ===")
	}

	rpnConverter := converter.NewRPNConverter()
	rpn := rpnConverter.Convert(tokens)

	fmt.Println("\n=== Обратная польская запись (ОПЗ) ===")
	fmt.Println(rpn)

	if *verbose {
		fmt.Printf("\n=== Статистика ===\n")
		fmt.Printf("Длина ОПЗ: %d символов\n", len(rpn))
		fmt.Printf("Количество элементов: %d\n", len(strings.Fields(rpn)))
	}
}
