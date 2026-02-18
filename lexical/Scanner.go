package lexical

import (
	"translator/shared"
)

type Scanner struct {
	buffer                   string
	pair_current_prev_states map[string]state
	path_to_file             string
}

func InitScanner(path_to_file string) *Scanner {
	return &Scanner{
		buffer: "",
		pair_current_prev_states: map[string]state{
			"current_state": state("S"),
		},
		path_to_file: path_to_file,
	}
}

func (scanner *Scanner) Run() string {
	sequence_leksems := ""

	fileReader := shared.InitFileReader(scanner.path_to_file)
	text_program := fileReader.GetTextFile()

	for _, sym := range text_program {

	}

	return sequence_leksems
}

func (scanner *Scanner) clearBuffer() {
	scanner.buffer = ""
}

func (scanner *Scanner) isLetter(symbol rune) bool {
	ans := false

	return ans
}

func (scanner *Scanner) isDigit(symbol rune) bool {
	ans := false

	return ans
}

func (scanner *Scanner) isDot(symbol rune) bool {
	ans := false

	return ans
}

func (scanner *Scanner) isSlash(symbol rune) bool {
	ans := false

	return ans
}

func (scanner *Scanner) isSeparator(symbol rune) bool {
	ans := false

	return ans
}

func (scanner *Scanner) isLessThanSym(symbol rune) bool {
	ans := false

	return ans
}

func (scanner *Scanner) isGreaterThanSym(symbol rune) bool {
	ans := false

	return ans
}
