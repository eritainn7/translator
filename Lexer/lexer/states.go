package lexer

import (
	"lexer/models"
)

func (l *Lexer) processChar() {
	ch := l.currentChar()

	switch l.state {
	case models.STATE_START:
		l.processStart(ch)
	case models.STATE_IDENTIFIER:
		l.processIdentifier(ch)
	case models.STATE_NUMBER:
		l.processNumber(ch)
	case models.STATE_NUMBER_HEX:
		l.processNumberHex(ch)
	case models.STATE_NUMBER_BIN:
		l.processNumberBin(ch)
	case models.STATE_NUMBER_REAL:
		l.processNumberReal(ch)
	case models.STATE_NUMBER_REAL_EXP:
		l.processNumberRealExp(ch)
	case models.STATE_CHAR:
		l.processCharLiteral(ch)
	case models.STATE_CHAR_ESCAPE:
		l.processCharEscape(ch)
	case models.STATE_STRING:
		l.processStringLiteral(ch)
	case models.STATE_STRING_ESCAPE:
		l.processStringEscape(ch)
	case models.STATE_COMMENT_LINE:
		l.processCommentLine(ch)
	case models.STATE_COMMENT_BLOCK:
		l.processCommentBlock(ch)
	case models.STATE_COMMENT_BLOCK_END:
		l.processCommentBlockEnd(ch)
	case models.STATE_OPERATION:
		l.processOperation(ch)
	case models.STATE_DELIMITER:
		l.processDelimiter(ch)
	}
}

func (l *Lexer) processStart(ch rune) {
	switch {
	case l.isWhitespace(ch):
		l.nextChar()

	case l.isLetter(ch):
		l.state = models.STATE_IDENTIFIER
		l.curLexeme = string(ch)
		l.startLine = l.line
		l.startColumn = l.column
		l.nextChar()

	case l.isDigit(ch):
		l.state = models.STATE_NUMBER
		l.curLexeme = string(ch)
		l.startLine = l.line
		l.startColumn = l.column
		l.nextChar()

	case ch == '\'':
		l.state = models.STATE_CHAR
		l.curLexeme = ""
		l.startLine = l.line
		l.startColumn = l.column
		l.nextChar()

	case ch == '"':
		l.state = models.STATE_STRING
		l.curLexeme = ""
		l.startLine = l.line
		l.startColumn = l.column
		l.nextChar()

	case ch == '/':
		l.state = models.STATE_OPERATION
		l.curLexeme = string(ch)
		l.startLine = l.line
		l.startColumn = l.column
		l.nextChar()

	case l.isDelimiter(ch):
		l.state = models.STATE_DELIMITER
		l.curLexeme = string(ch)
		l.startLine = l.line
		l.startColumn = l.column
		l.nextChar()

	default:
		// Проверяем на операции
		if _, exists := l.tables.Operations[string(ch)]; exists {
			l.state = models.STATE_OPERATION
			l.curLexeme = string(ch)
			l.startLine = l.line
			l.startColumn = l.column
			l.nextChar()
		} else {
			// Недопустимый символ
			l.startLine = l.line
			l.startColumn = l.column
			l.addLexeme(models.L_ERROR, string(ch))
			l.nextChar()
		}
	}
}

func (l *Lexer) processIdentifier(ch rune) {
	if l.isLetter(ch) || l.isDigit(ch) {
		l.curLexeme += string(ch)
		l.nextChar()
	} else {
		// Проверяем, является ли идентификатор служебным словом
		if _, exists := l.tables.ServiceWords[l.curLexeme]; exists {
			l.addLexeme(models.L_SERVICE, l.curLexeme)
		} else {
			l.addLexeme(models.L_IDENTIFIER, l.curLexeme)
		}
		l.state = models.STATE_START
	}
}

func (l *Lexer) processNumber(ch rune) {
	// Проверка на шестнадцатеричный или двоичный формат
	if len(l.curLexeme) == 1 && l.curLexeme == "0" {
		if ch == 'x' || ch == 'X' {
			l.curLexeme += string(ch)
			l.state = models.STATE_NUMBER_HEX
			l.nextChar()
			return
		} else if ch == 'b' || ch == 'B' {
			l.curLexeme += string(ch)
			l.state = models.STATE_NUMBER_BIN
			l.nextChar()
			return
		}
	}

	if l.isDigit(ch) {
		l.curLexeme += string(ch)
		l.nextChar()
	} else if ch == '.' {
		l.curLexeme += string(ch)
		l.state = models.STATE_NUMBER_REAL
		l.nextChar()
	} else if ch == 'e' || ch == 'E' {
		l.curLexeme += string(ch)
		l.state = models.STATE_NUMBER_REAL_EXP
		l.nextChar()
	} else {
		l.addLexeme(models.L_CONSTANT, l.curLexeme)
		l.state = models.STATE_START
	}
}

func (l *Lexer) processNumberHex(ch rune) {
	if l.isHexDigit(ch) {
		l.curLexeme += string(ch)
		l.nextChar()
	} else {
		l.addLexeme(models.L_CONSTANT, l.curLexeme)
		l.state = models.STATE_START
	}
}

func (l *Lexer) processNumberBin(ch rune) {
	if l.isBinDigit(ch) {
		l.curLexeme += string(ch)
		l.nextChar()
	} else {
		l.addLexeme(models.L_CONSTANT, l.curLexeme)
		l.state = models.STATE_START
	}
}

func (l *Lexer) processNumberReal(ch rune) {
	if l.isDigit(ch) {
		l.curLexeme += string(ch)
		l.nextChar()
	} else if ch == 'e' || ch == 'E' {
		l.curLexeme += string(ch)
		l.state = models.STATE_NUMBER_REAL_EXP
		l.nextChar()
	} else {
		l.addLexeme(models.L_CONSTANT, l.curLexeme)
		l.state = models.STATE_START
	}
}

func (l *Lexer) processNumberRealExp(ch rune) {
	if ch == '+' || ch == '-' || l.isDigit(ch) {
		l.curLexeme += string(ch)
		l.nextChar()
		l.state = models.STATE_NUMBER_REAL
	} else {
		l.addLexeme(models.L_CONSTANT, l.curLexeme)
		l.state = models.STATE_START
	}
}

func (l *Lexer) processCharLiteral(ch rune) {
	if ch == '\\' {
		l.curLexeme += string(ch)
		l.state = models.STATE_CHAR_ESCAPE
		l.nextChar()
	} else if ch == '\'' {
		l.addLexeme(models.L_CONSTANT, "'"+l.curLexeme+"'")
		l.state = models.STATE_START
		l.nextChar()
	} else {
		l.curLexeme += string(ch)
		l.state = models.STATE_CHAR
		l.nextChar()
	}
}

func (l *Lexer) processCharEscape(ch rune) {
	l.curLexeme += string(ch)
	l.state = models.STATE_CHAR
	l.nextChar()
}

func (l *Lexer) processStringLiteral(ch rune) {
	if ch == '\\' {
		l.curLexeme += string(ch)
		l.state = models.STATE_STRING_ESCAPE
		l.nextChar()
	} else if ch == '"' {
		l.addLexeme(models.L_CONSTANT, "\""+l.curLexeme+"\"")
		l.state = models.STATE_START
		l.nextChar()
	} else {
		l.curLexeme += string(ch)
		l.nextChar()
	}
}

func (l *Lexer) processStringEscape(ch rune) {
	l.curLexeme += string(ch)
	l.state = models.STATE_STRING
	l.nextChar()
}

func (l *Lexer) processCommentLine(ch rune) {
	if ch == '\n' {
		l.state = models.STATE_START
		l.nextChar()
	} else {
		l.nextChar()
	}
}

func (l *Lexer) processCommentBlock(ch rune) {
	if ch == '*' {
		l.state = models.STATE_COMMENT_BLOCK_END
		l.nextChar()
	} else {
		l.nextChar()
	}
}

func (l *Lexer) processCommentBlockEnd(ch rune) {
	if ch == '/' {
		l.state = models.STATE_START
		l.nextChar()
	} else {
		l.state = models.STATE_COMMENT_BLOCK
	}
}

func (l *Lexer) processOperation(ch rune) {
	// Проверка на двухсимвольные операции
	nextCh := rune(0)
	if l.pos < len(l.input) {
		nextCh = rune(l.input[l.pos])
	}

	if nextCh != 0 {
		twoCharOp := l.curLexeme + string(nextCh)
		if _, exists := l.tables.Operations[twoCharOp]; exists {
			l.curLexeme = twoCharOp
			l.nextChar()
			l.nextChar()
			l.addLexeme(models.L_OPERATION, l.curLexeme)
			l.state = models.STATE_START
			return
		}
	}

	// Проверка на комментарии
	if l.curLexeme == "/" && nextCh == '/' {
		l.state = models.STATE_COMMENT_LINE
		l.nextChar()
		l.nextChar()
	} else if l.curLexeme == "/" && nextCh == '*' {
		l.state = models.STATE_COMMENT_BLOCK
		l.nextChar()
		l.nextChar()
	} else {
		// Односимвольная операция
		l.addLexeme(models.L_OPERATION, l.curLexeme)
		l.state = models.STATE_START
	}
}

func (l *Lexer) processDelimiter(ch rune) {
	// Проверка на двухсимвольные разделители
	nextCh := rune(0)
	if l.pos < len(l.input) {
		nextCh = rune(l.input[l.pos])
	}

	if nextCh != 0 {
		twoCharDelim := l.curLexeme + string(nextCh)
		if _, exists := l.tables.Delimiters[twoCharDelim]; exists {
			l.curLexeme = twoCharDelim
			l.nextChar()
			l.nextChar()
			l.addLexeme(models.L_DELIMITER, l.curLexeme)
			l.state = models.STATE_START
			return
		}
	}

	// Односимвольный разделитель
	l.addLexeme(models.L_DELIMITER, l.curLexeme)
	l.state = models.STATE_START
}
