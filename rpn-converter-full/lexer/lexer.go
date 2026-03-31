package lexer

import (
	"strings"
	"unicode"
)

// Lexer представляет лексический анализатор
type Lexer struct {
	input string
	pos   int
	line  int
	col   int
	ch    rune
}

// NewLexer создает новый лексический анализатор
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
		pos:   0,
		line:  1,
		col:   0,
	}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.pos < len(l.input) {
		l.ch = rune(l.input[l.pos])
		l.pos++
		l.col++
	} else {
		l.ch = 0
	}
}

func (l *Lexer) peekChar() rune {
	if l.pos < len(l.input) {
		return rune(l.input[l.pos])
	}
	return 0
}

func (l *Lexer) skipWhitespace() {
	for l.ch != 0 && (unicode.IsSpace(l.ch) || l.ch == '\n' || l.ch == '\r') {
		if l.ch == '\n' {
			l.line++
			l.col = 0
		}
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos - 1
	for l.ch != 0 && (unicode.IsLetter(l.ch) || unicode.IsDigit(l.ch) || l.ch == '_') {
		l.readChar()
	}
	return l.input[start : l.pos-1]
}

func (l *Lexer) readNumber() string {
	start := l.pos - 1
	for l.ch != 0 && (unicode.IsDigit(l.ch) || l.ch == '.' || l.ch == 'x' || l.ch == 'X' ||
		l.ch == 'b' || l.ch == 'B' || l.ch == 'e' || l.ch == 'E' || l.ch == '+' || l.ch == '-') {
		l.readChar()
	}
	return l.input[start : l.pos-1]
}

func (l *Lexer) readString() string {
	start := l.pos
	l.readChar() // пропускаем открывающую кавычку
	for l.ch != 0 && l.ch != '"' {
		if l.ch == '\\' {
			l.readChar() // пропускаем escape-символ
		}
		l.readChar()
	}
	end := l.pos - 1
	l.readChar() // пропускаем закрывающую кавычку
	return l.input[start:end]
}

func (l *Lexer) readCharLiteral() string {
	l.readChar() // пропускаем '
	value := string(l.ch)
	l.readChar() // пропускаем символ
	if l.ch == '\'' {
		l.readChar() // пропускаем '
	}
	return value
}

func (l *Lexer) skipLineComment() {
	for l.ch != 0 && l.ch != '\n' {
		l.readChar()
	}
}

func (l *Lexer) skipBlockComment() {
	l.readChar() // пропускаем *
	for !(l.ch == '*' && l.peekChar() == '/') && l.ch != 0 {
		if l.ch == '\n' {
			l.line++
			l.col = 0
		}
		l.readChar()
	}
	l.readChar() // пропускаем *
	l.readChar() // пропускаем /
}

// NextToken возвращает следующий токен
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	if l.ch == 0 {
		return Token{Type: TokenEOF, Line: l.line, Col: l.col}
	}

	// строки
	if l.ch == '"' {
		return Token{Type: TokenString, Value: l.readString(), Line: l.line, Col: l.col}
	}

	// символы
	if l.ch == '\'' {
		return Token{Type: TokenChar, Value: l.readCharLiteral(), Line: l.line, Col: l.col}
	}

	// идентификаторы и ключевые слова
	if unicode.IsLetter(l.ch) || l.ch == '_' {
		ident := l.readIdentifier()
		if tokType, ok := keywords[strings.ToLower(ident)]; ok {
			return Token{Type: tokType, Value: ident, Line: l.line, Col: l.col}
		}
		return Token{Type: TokenIdent, Value: ident, Line: l.line, Col: l.col}
	}

	// числа
	if unicode.IsDigit(l.ch) {
		return Token{Type: TokenNumber, Value: l.readNumber(), Line: l.line, Col: l.col}
	}

	// операторы и специальные символы
	tok := l.readOperator()
	l.readChar()
	return tok
}

func (l *Lexer) readOperator() Token {
	switch l.ch {
	case '(':
		return Token{Type: TokenLParen, Value: "(", Line: l.line, Col: l.col}
	case ')':
		return Token{Type: TokenRParen, Value: ")", Line: l.line, Col: l.col}
	case '[':
		return Token{Type: TokenLBracket, Value: "[", Line: l.line, Col: l.col}
	case ']':
		return Token{Type: TokenRBracket, Value: "]", Line: l.line, Col: l.col}
	case '{':
		return Token{Type: TokenLBrace, Value: "{", Line: l.line, Col: l.col}
	case '}':
		return Token{Type: TokenRBrace, Value: "}", Line: l.line, Col: l.col}
	case ',':
		return Token{Type: TokenComma, Value: ",", Line: l.line, Col: l.col}
	case ';':
		return Token{Type: TokenSemicolon, Value: ";", Line: l.line, Col: l.col}
	case ':':
		return Token{Type: TokenColon, Value: ":", Line: l.line, Col: l.col}
	case '?':
		if l.peekChar() == '?' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "??", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenQuestion, Value: "?", Line: l.line, Col: l.col}
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "==", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '>' {
			l.readChar()
			return Token{Type: TokenLambda, Value: "=>", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenAssign, Value: "=", Line: l.line, Col: l.col}
	case '+':
		if l.peekChar() == '+' {
			l.readChar()
			return Token{Type: TokenInc, Value: "++", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenPlusAssign, Value: "+=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "+", Line: l.line, Col: l.col}
	case '-':
		if l.peekChar() == '-' {
			l.readChar()
			return Token{Type: TokenDec, Value: "--", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenMinusAssign, Value: "-=", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '>' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "->", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "-", Line: l.line, Col: l.col}
	case '*':
		if l.peekChar() == '*' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "**", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenMulAssign, Value: "*=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "*", Line: l.line, Col: l.col}
	case '/':
		if l.peekChar() == '/' {
			l.readChar()
			l.skipLineComment()
			return l.NextToken()
		}
		if l.peekChar() == '*' {
			l.readChar()
			l.skipBlockComment()
			return l.NextToken()
		}
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenDivAssign, Value: "/=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "/", Line: l.line, Col: l.col}
	case '%':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenModAssign, Value: "%=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "%", Line: l.line, Col: l.col}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "!=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "!", Line: l.line, Col: l.col}
	case '<':
		if l.peekChar() == '<' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				return Token{Type: TokenOperator, Value: "<<=", Line: l.line, Col: l.col}
			}
			return Token{Type: TokenOperator, Value: "<<", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "<=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "<", Line: l.line, Col: l.col}
	case '>':
		if l.peekChar() == '>' {
			l.readChar()
			if l.peekChar() == '=' {
				l.readChar()
				return Token{Type: TokenOperator, Value: ">>=", Line: l.line, Col: l.col}
			}
			return Token{Type: TokenOperator, Value: ">>", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenOperator, Value: ">=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: ">", Line: l.line, Col: l.col}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "&&", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "&=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "&", Line: l.line, Col: l.col}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "||", Line: l.line, Col: l.col}
		}
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "|=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "|", Line: l.line, Col: l.col}
	case '^':
		if l.peekChar() == '=' {
			l.readChar()
			return Token{Type: TokenOperator, Value: "^=", Line: l.line, Col: l.col}
		}
		return Token{Type: TokenOperator, Value: "^", Line: l.line, Col: l.col}
	case '~':
		return Token{Type: TokenOperator, Value: "~", Line: l.line, Col: l.col}
	default:
		return Token{Type: TokenOperator, Value: string(l.ch), Line: l.line, Col: l.col}
	}
}

// Tokenize возвращает все токены
func (l *Lexer) Tokenize() []Token {
	tokens := make([]Token, 0)
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}

// keywords карта ключевых слов
var keywords = map[string]TokenType{
	"if":        TokenIf,
	"else":      TokenElse,
	"switch":    TokenSwitch,
	"case":      TokenCase,
	"default":   TokenDefault,
	"for":       TokenFor,
	"foreach":   TokenForeach,
	"while":     TokenWhile,
	"do":        TokenDo,
	"break":     TokenBreak,
	"continue":  TokenContinue,
	"return":    TokenReturn,
	"goto":      TokenGoto,
	"try":       TokenTry,
	"catch":     TokenCatch,
	"finally":   TokenFinally,
	"throw":     TokenThrow,
	"class":     TokenClass,
	"struct":    TokenStruct,
	"interface": TokenInterface,
	"enum":      TokenEnum,
	"delegate":  TokenDelegate,
	"event":     TokenEvent,
	"public":    TokenPublic,
	"private":   TokenPrivate,
	"protected": TokenProtected,
	"internal":  TokenInternal,
	"static":    TokenStatic,
	"readonly":  TokenReadonly,
	"const":     TokenConst,
	"virtual":   TokenVirtual,
	"override":  TokenOverride,
	"abstract":  TokenAbstract,
	"sealed":    TokenSealed,
	"partial":   TokenPartial,
	"new":       TokenNew,
	"this":      TokenThis,
	"base":      TokenBase,
	"as":        TokenAs,
	"is":        TokenIs,
	"typeof":    TokenTypeof,
	"sizeof":    TokenSizeof,
	"checked":   TokenChecked,
	"unchecked": TokenUnchecked,
	"lock":      TokenLock,
	"using":     TokenUsing,
	"namespace": TokenNamespace,
	"import":    TokenImport,
	"var":       TokenVar,
	"dynamic":   TokenDynamic,
	"async":     TokenAsync,
	"await":     TokenAwait,
	"yield":     TokenYield,
	"where":     TokenWhere,
	"select":    TokenSelect,
	"from":      TokenFrom,
	"in":        TokenIn,
	"true":      TokenBool,
	"false":     TokenBool,
}
