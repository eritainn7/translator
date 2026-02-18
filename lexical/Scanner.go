package lexical

import (
	"translator/shared"
)

type Scanner struct {
	buffer                   string
	pair_current_prev_states [2]state
	path_to_file             string
}

func InitScanner(path_to_file string) *Scanner {
	return &Scanner{
		buffer:                   "",
		pair_current_prev_states: [2]state{state("S")},
		path_to_file:             path_to_file,
	}
}

func (scanner *Scanner) Run() string {
	sequence_leksems := ""

	fileReader := shared.InitFileReader(scanner.path_to_file)
	text_program := fileReader.GetTextFile()

	for i, sym := range text_program {
		//In this place all main logic!!!
	}

	return sequence_leksems
}

func (scanner *Scanner) clearBuffer() {
	scanner.buffer = ""
}
