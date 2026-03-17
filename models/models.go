package models

import "fmt"

// Типы лексем
const (
	L_SERVICE    = 'W' // служебное слово
	L_IDENTIFIER = 'I' // идентификатор
	L_CONSTANT   = 'C' // константа
	L_DELIMITER  = 'R' // разделитель
	L_OPERATION  = 'O' // операция
	L_EOF        = 'E' // конец файла
	L_ERROR      = 'F' // ошибка
)

// Состояния автомата
const (
	STATE_START = iota
	STATE_IDENTIFIER
	STATE_NUMBER
	STATE_NUMBER_HEX
	STATE_NUMBER_BIN
	STATE_NUMBER_REAL
	STATE_NUMBER_REAL_EXP
	STATE_CHAR
	STATE_CHAR_ESCAPE
	STATE_STRING
	STATE_STRING_ESCAPE
	STATE_COMMENT_LINE
	STATE_COMMENT_BLOCK
	STATE_COMMENT_BLOCK_END
	STATE_OPERATION
	STATE_DELIMITER
	STATE_ERROR
)

// Лексема
type Lexeme struct {
	Type   byte
	Index  int
	Value  string
	Line   int
	Column int
}

func (l Lexeme) String() string {
	return fmt.Sprintf("%c%d", l.Type, l.Index)
}

func (l Lexeme) FullInfo() string {
	return fmt.Sprintf("Type: %c, Index: %d, Value: %q, Line: %d, Column: %d",
		l.Type, l.Index, l.Value, l.Line, l.Column)
}
