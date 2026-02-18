package main

import (
	"fmt"
	"translator/lexical"
)

func main() {
	scanner := lexical.InitScanner("mocks/prog.cs")
	seq_leksems := scanner.Run()

	fmt.Println(seq_leksems)
}
