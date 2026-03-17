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

	// Для определения области видимости
	scopeStack   []string
	currentScope string
}

func NewLexer(input string, tables *tables.Tables) *Lexer {
	return &Lexer{
		input:        input,
		pos:          0,
		line:         1,
		column:       1,
		lexemes:      make([]models.Lexeme, 0),
		state:        models.STATE_START,
		tables:       tables,
		scopeStack:   []string{"global"},
		currentScope: "global",
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

// Определение типа идентификатора на основе контекста
func (l *Lexer) determineIdentifierType(name string) string {
	// Проверяем, не является ли идентификатор служебным словом
	if _, exists := l.tables.ServiceWords[name]; exists {
		return "служебное_слово"
	}

	// Простой эвристический анализ типа
	// По умолчанию считаем переменной
	idType := models.ID_TYPE_VARIABLE

	// Проверка на функцию (если следующий символ '(')
	nextCh := l.currentChar()
	if nextCh == '(' {
		idType = models.ID_TYPE_FUNCTION
	}

	// Проверка на класс/структуру (по соглашению именования с большой буквы)
	if len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z' {
		idType = models.ID_TYPE_CLASS
	}

	return idType
}

// Обновление области видимости на основе служебных слов
func (l *Lexer) updateScope(lexeme models.Lexeme) {
	switch lexeme.Value {
	case "class", "struct", "interface", "namespace":
		// Ожидаем, что следующим будет идентификатор - имя класса/пространства имен
		l.currentScope = "waiting_for_name"
	case "{":
		// Вход в блок - создаем новую область видимости
		newScope := fmt.Sprintf("%s_block_%d", l.currentScope, len(l.scopeStack))
		l.scopeStack = append(l.scopeStack, newScope)
		l.currentScope = newScope
		l.tables.SetCurrentScope(newScope)
	case "}":
		// Выход из блока
		if len(l.scopeStack) > 1 {
			l.scopeStack = l.scopeStack[:len(l.scopeStack)-1]
			l.currentScope = l.scopeStack[len(l.scopeStack)-1]
			l.tables.SetCurrentScope(l.currentScope)
		}
	}
}

func (l *Lexer) addLexeme(lexType byte, value string) {
	var index int

	switch lexType {
	case models.L_SERVICE:
		index = l.tables.ServiceWords[value]
	case models.L_IDENTIFIER:
		// Определяем тип идентификатора
		idType := l.determineIdentifierType(value)

		// Если ожидаем имя для класса/функции, обновляем область видимости
		if l.currentScope == "waiting_for_name" {
			// Создаем новую область видимости с именем
			newScope := value
			l.scopeStack = append(l.scopeStack, newScope)
			l.currentScope = newScope
			l.tables.SetCurrentScope(newScope)
			l.currentScope = newScope
		}

		// Получаем ID с дополнительной информацией
		index = l.tables.GetIdentifierID(value, l.startLine, l.startColumn, l.currentScope, idType)
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

	// Обновляем область видимости на основе лексемы
	if lexType == models.L_DELIMITER || lexType == models.L_SERVICE {
		l.updateScope(lexeme)
	}

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
