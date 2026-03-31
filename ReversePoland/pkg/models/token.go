package models

// TokenType представляет тип лексемы
type TokenType int

const (
	TokenIdent TokenType = iota
	TokenNumber
	TokenString
	TokenOperator
	TokenLParen
	TokenRParen
	TokenLBracket
	TokenRBracket
	TokenLBrace
	TokenRBrace
	TokenComma
	TokenSemicolon
	TokenAssign
	TokenIf
	TokenThen
	TokenElse
	TokenGoto
	TokenProc
	TokenEnd
	TokenDcl
	TokenFor
	TokenWhile
	TokenDo
	TokenSwitch
	TokenCase
	TokenDefault
	TokenBreak
	TokenContinue
	TokenReturn
	TokenIncrement
	TokenDecrement
	TokenQuestion
	TokenColon
	TokenNullCoalescing
	TokenLambda
	TokenEOF

	// Ключевые слова для классов
	TokenClass
	TokenPublic
	TokenPrivate
	TokenProtected
	TokenInternal
	TokenStatic
	TokenVirtual
	TokenOverride
	TokenAbstract
	TokenSealed
	TokenPartial
	TokenInterface
	TokenImplements
	TokenExtends
	TokenNew
	TokenThis
	TokenBase
	TokenConstructor
	TokenDestructor
	TokenProperty
	TokenGet
	TokenSet
	TokenInit
	TokenEvent
	TokenDelegate
)

// Token представляет лексему
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

// TokenTypeName возвращает строковое представление типа токена
func (t TokenType) String() string {
	names := map[TokenType]string{
		TokenIdent:     "IDENT",
		TokenNumber:    "NUMBER",
		TokenString:    "STRING",
		TokenOperator:  "OPERATOR",
		TokenLParen:    "LPAREN",
		TokenRParen:    "RPAREN",
		TokenLBracket:  "LBRACKET",
		TokenRBracket:  "RBRACKET",
		TokenLBrace:    "LBRACE",
		TokenRBrace:    "RBRACE",
		TokenComma:     "COMMA",
		TokenSemicolon: "SEMICOLON",
		TokenAssign:    "ASSIGN",
		TokenIf:        "IF",
		TokenThen:      "THEN",
		TokenElse:      "ELSE",
		TokenFor:       "FOR",
		TokenWhile:     "WHILE",
		TokenClass:     "CLASS",
		TokenPublic:    "PUBLIC",
		TokenPrivate:   "PRIVATE",
		TokenReturn:    "RETURN",
		TokenEOF:       "EOF",
	}

	if name, ok := names[t]; ok {
		return name
	}
	return "UNKNOWN"
}
