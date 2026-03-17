package lexer

import (
	"fmt"
	"lexer/models"
	"lexer/tables"
)

// Лексический анализатор
type Lexer struct {
	input       string
	pos         int
	line        int
	column      int
	lexemes     []models.Lexeme
	state       int
	curLexeme   string
	startLine   int
	startColumn int
	tables      *tables.Tables
}

func NewLexer(input string, tables *tables.Tables) *Lexer {
	return &Lexer{
		input:   input,
		pos:     0,
		line:    1,
		column:  1,
		lexemes: make([]models.Lexeme, 0),
		state:   models.STATE_START,
		tables:  tables,
	}
}

func (l *Lexer) currentChar() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return rune(l.input[l.pos])
}

func (l *Lexer) nextChar() {
	if l.pos < len(l.input) {
		if l.input[l.pos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.pos++
	}
}

func (l *Lexer) addLexeme(lexType byte, value string) {
	var index int
	switch lexType {
	case models.L_SERVICE:
		index = l.tables.ServiceWords[value]
	case models.L_IDENTIFIER:
		index = l.tables.GetIdentifierID(value)
	case models.L_CONSTANT:
		index = l.tables.GetConstantID(value)
	case models.L_DELIMITER:
		index = l.tables.Delimiters[value]
	case models.L_OPERATION:
		index = l.tables.Operations[value]
	case models.L_ERROR:
		index = 0
	}

	lexeme := models.Lexeme{
		Type:   lexType,
		Index:  index,
		Value:  value,
		Line:   l.startLine,
		Column: l.startColumn,
	}

	l.lexemes = append(l.lexemes, lexeme)

	// Вывод в требуемом формате
	if lexType != models.L_ERROR {
		fmt.Printf("%c%d ", lexType, index)
	}
}

func (l *Lexer) isDigit(ch rune) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) isHexDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')
}

func (l *Lexer) isBinDigit(ch rune) bool {
	return ch == '0' || ch == '1'
}

func (l *Lexer) isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func (l *Lexer) isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func (l *Lexer) isDelimiter(ch rune) bool {
	_, exists := l.tables.Delimiters[string(ch)]
	return exists
}

func (l *Lexer) Analyze() []models.Lexeme {
	for l.pos < len(l.input) {
		l.processChar()
	}

	// Обработка возможного незавершенного состояния
	l.processEndOfFile()

	// Добавляем маркер конца файла
	l.startLine = l.line
	l.startColumn = l.column
	l.addLexeme(models.L_EOF, "EOF")

	return l.lexemes
}

func (l *Lexer) processEndOfFile() {
	// Обработка незавершенных лексем при достижении конца файла
	switch l.state {
	case models.STATE_IDENTIFIER:
		if _, exists := l.tables.ServiceWords[l.curLexeme]; exists {
			l.addLexeme(models.L_SERVICE, l.curLexeme)
		} else {
			l.addLexeme(models.L_IDENTIFIER, l.curLexeme)
		}
	case models.STATE_NUMBER, models.STATE_NUMBER_HEX,
		models.STATE_NUMBER_BIN, models.STATE_NUMBER_REAL:
		l.addLexeme(models.L_CONSTANT, l.curLexeme)
	case models.STATE_OPERATION:
		l.addLexeme(models.L_OPERATION, l.curLexeme)
	case models.STATE_DELIMITER:
		l.addLexeme(models.L_DELIMITER, l.curLexeme)
	case models.STATE_CHAR, models.STATE_STRING:
		l.addLexeme(models.L_ERROR, "Unterminated literal")
	}
}

func (l *Lexer) GetLexemes() []models.Lexeme {
	return l.lexemes
}
