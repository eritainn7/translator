package lexical

import (
	"fmt"
	"translator/lexical/bases/model"
	"translator/lexical/semantics"
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
		scanner.buffer += string(sym)

		scanner.updateState(sym)

		//stateTransition() launch semantic procedures after change state and return code in format(<litter><digit>: I8 W2 R3 etc)
		sequence_leksems += scanner.stateTranstitions(sym) + " "
	}

	return sequence_leksems
}

func (scanner *Scanner) updateState(current_symbol rune) {
	if scanner.pair_current_prev_states["current_state"] == state("S") {
		if scanner.isLetter(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("0"), "prev_state": state("S")}
		} else if scanner.isDigit(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("2"), "prev_state": state("S")}
		} else if scanner.isDot(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("4"), "prev_state": state("S")}
		} else if scanner.isSlash(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("5"), "prev_state": state("S")}
		} else if scanner.isLessThanSym(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("9"), "prev_state": state("S")}
		} else if scanner.isSeparator(current_symbol) || scanner.isOperation(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("8"), "prev_state": state("S")}
		}
		// else {
		// 	scanner.pair_current_prev_states = map[string]state{"current_state": state("F"), "prev_state": state("S")}
		// }
	} else if scanner.pair_current_prev_states["current_state"] == state("0") {
		if scanner.isLetter(current_symbol) {
			scanner.pair_current_prev_states["prev_state"] = state("0")
		} else if scanner.isDigit(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("1"), "prev_state": state("0")}
		} else if scanner.isSeparator(current_symbol) || scanner.isOperation(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("S"), "prev_state": state("0")}
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("1") {
		if scanner.isLetter(current_symbol) || scanner.isDigit(current_symbol) {
			scanner.pair_current_prev_states["prev_state"] = state("1")
		} else if scanner.isSeparator(current_symbol) || scanner.isOperation(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("S"), "prev_state": state("1")}
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("2") {
		if scanner.isDigit(current_symbol) {
			scanner.pair_current_prev_states["prev_state"] = state("2")
		} else if scanner.isDot(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("3"), "prev_state": state("2")}
		} else if scanner.isSeparator(current_symbol) || scanner.isOperation(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("S"), "prev_state": state("2")}
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("3") {
		if scanner.isDigit(current_symbol) {
			scanner.pair_current_prev_states["prev_state"] = state("3")
		} else if scanner.isSeparator(current_symbol) || scanner.isOperation(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("S"), "prev_state": state("3")}
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("4") {
		if scanner.isDigit(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("3"), "prev_state": state("4")}
		} else {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("S"), "prev_state": state("4")}
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("5") {
		if scanner.isSlash(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("6"), "prev_state": state("5")}
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("6") {
		if scanner.isSlash(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("7"), "prev_state": state("6")}
		} else {
			scanner.pair_current_prev_states["prev_state"] = state("6")
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("7") {
		scanner.pair_current_prev_states = map[string]state{"current_state": state("S"), "prev_state": state("7")}
	} else if scanner.pair_current_prev_states["current_state"] == state("8") {
		if !scanner.isSeparator(current_symbol) && !scanner.isOperation(current_symbol) && !scanner.isGreaterThanSym(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("S"), "prev_state": state("8")}
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("9") {
		if scanner.isGreaterThanSym(current_symbol) {
			scanner.pair_current_prev_states = map[string]state{"current_state": state("8"), "prev_state": state("9")}
		}
	} else if scanner.pair_current_prev_states["current_state"] == state("F") {
		fmt.Println("Incorrect program!")
		return
	}

}

func (scanner *Scanner) stateTranstitions(current_symbol rune) string {
	var mark_code string

	switch scanner.pair_current_prev_states["prev_state"] {
	case state("S"):
		{
			if scanner.pair_current_prev_states["current_state"] == state("8") {
				if scanner.isOperation(current_symbol) {
					mark_code = semantics.InitSemanticHandler6(scanner.buffer).Run()
				} else if scanner.isSeparator(current_symbol) {
					mark_code = semantics.InitSemanticHandler4(scanner.buffer).Run()
				}
			}
		}
	case state("0"):
		{
			if scanner.pair_current_prev_states["current_state"] == state("S") {
				mark_code = semantics.InitSemanticHandler2(scanner.buffer).Run()
			}
		}
	case state("1"):
		{
			if scanner.pair_current_prev_states["current_state"] == state("S") {
				mark_code = semantics.InitSemanticHandler1(scanner.buffer).Run()
			}
		}
	case state("2"):
		{
			if scanner.pair_current_prev_states["current_state"] == state("S") {
				mark_code = semantics.InitSemanticHandler3(scanner.buffer).Run()
			}
		}
	case state("3"):
		{
			if scanner.pair_current_prev_states["current_state"] == state("S") {
				mark_code = semantics.InitSemanticHandler3(scanner.buffer).Run()
			}
		}
	case state("4"):
		{
			if scanner.pair_current_prev_states["current_state"] == state("S") {
				mark_code = semantics.InitSemanticHandler4(scanner.buffer).Run()
			}
		}
	case state("6"):
		{
			if scanner.pair_current_prev_states["current_state"] == state("7") {
				semantics.SemanticProcedure5(&scanner.buffer)
			}
		}
	case state("9"):
		{
			if scanner.pair_current_prev_states["current_state"] == state("8") {
				mark_code, _ = semantics.InitSemanticHandler7(scanner.buffer).Run(string(current_symbol))
			} else if scanner.pair_current_prev_states["current_state"] == state("S") {
				mark_code = semantics.InitSemanticHandler6(scanner.buffer).Run()
			}
		}

	}

	scanner.clearBuffer()
	return mark_code
}

func (scanner *Scanner) clearBuffer() {
	scanner.buffer = ""
}

func (scanner *Scanner) isLetter(symbol rune) bool {
	list_letters := "qwertyuiopasdfghjklzxcvbnm"
	for _, letter := range list_letters {
		if letter == symbol {
			return true
		}
	}

	return false
}

func (scanner *Scanner) isDigit(symbol rune) bool {
	list_digits := "0123456789"
	for _, digit := range list_digits {
		if digit == symbol {
			return true
		}
	}

	return false
}

func (scanner *Scanner) isDot(symbol rune) bool {
	return symbol == '.'
}

func (scanner *Scanner) isSlash(symbol rune) bool {
	return symbol == '/'
}

func (scanner *Scanner) isOperation(symbol rune) bool {
	code, _ := model.FindCodeInCSV(model.Operations, string(symbol))
	return code != -1
}

func (scanner *Scanner) isSeparator(symbol rune) bool {
	code, _ := model.FindCodeInCSV(model.Separators, string(symbol))
	return code != -1
}

func (scanner *Scanner) isLessThanSym(symbol rune) bool {
	return symbol == '<'
}

func (scanner *Scanner) isGreaterThanSym(symbol rune) bool {
	return symbol == '>'
}
