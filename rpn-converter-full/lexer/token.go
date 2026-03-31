package lexer

// TokenType представляет тип лексемы
type TokenType int

const (
	TokenIdent TokenType = iota
	TokenNumber
	TokenString
	TokenChar
	TokenBool
	TokenOperator
	TokenLParen
	TokenRParen
	TokenLBracket
	TokenRBracket
	TokenLBrace
	TokenRBrace
	TokenComma
	TokenSemicolon
	TokenColon
	TokenQuestion
	TokenAssign
	TokenPlusAssign
	TokenMinusAssign
	TokenMulAssign
	TokenDivAssign
	TokenModAssign
	TokenInc
	TokenDec
	TokenIf
	TokenElse
	TokenSwitch
	TokenCase
	TokenFor
	TokenForeach
	TokenWhile
	TokenDo
	TokenBreak
	TokenContinue
	TokenReturn
	TokenGoto
	TokenTry
	TokenCatch
	TokenFinally
	TokenThrow
	TokenClass
	TokenStruct
	TokenInterface
	TokenEnum
	TokenDelegate
	TokenEvent
	TokenProperty
	TokenIndexer
	TokenPublic
	TokenPrivate
	TokenProtected
	TokenInternal
	TokenStatic
	TokenReadonly
	TokenConst
	TokenVirtual
	TokenOverride
	TokenAbstract
	TokenSealed
	TokenPartial
	TokenNew
	TokenThis
	TokenBase
	TokenAs
	TokenIs
	TokenTypeof
	TokenSizeof
	TokenDefault
	TokenChecked
	TokenUnchecked
	TokenLock
	TokenUsing
	TokenNamespace
	TokenImport
	TokenVar
	TokenDynamic
	TokenAsync
	TokenAwait
	TokenYield
	TokenWhere
	TokenSelect
	TokenFrom
	TokenIn
	TokenLambda
	TokenEOF
)

// Token представляет лексему
type Token struct {
	Type  TokenType
	Value string
	Line  int
	Col   int
}

// String возвращает строковое представление типа токена
func (t TokenType) String() string {
	names := map[TokenType]string{
		TokenIdent:     "IDENT",
		TokenNumber:    "NUMBER",
		TokenString:    "STRING",
		TokenChar:      "CHAR",
		TokenBool:      "BOOL",
		TokenOperator:  "OPERATOR",
		TokenLParen:    "LPAREN",
		TokenRParen:    "RPAREN",
		TokenLBracket:  "LBRACKET",
		TokenRBracket:  "RBRACKET",
		TokenLBrace:    "LBRACE",
		TokenRBrace:    "RBRACE",
		TokenComma:     "COMMA",
		TokenSemicolon: "SEMICOLON",
		TokenColon:     "COLON",
		TokenQuestion:  "QUESTION",
		TokenAssign:    "ASSIGN",
		TokenInc:       "INC",
		TokenDec:       "DEC",
		TokenIf:        "IF",
		TokenElse:      "ELSE",
		TokenSwitch:    "SWITCH",
		TokenCase:      "CASE",
		TokenDefault:   "DEFAULT",
		TokenFor:       "FOR",
		TokenForeach:   "FOREACH",
		TokenWhile:     "WHILE",
		TokenDo:        "DO",
		TokenBreak:     "BREAK",
		TokenContinue:  "CONTINUE",
		TokenReturn:    "RETURN",
		TokenGoto:      "GOTO",
		TokenTry:       "TRY",
		TokenCatch:     "CATCH",
		TokenFinally:   "FINALLY",
		TokenThrow:     "THROW",
		TokenClass:     "CLASS",
		TokenStruct:    "STRUCT",
		TokenInterface: "INTERFACE",
		TokenEnum:      "ENUM",
		TokenDelegate:  "DELEGATE",
		TokenPublic:    "PUBLIC",
		TokenPrivate:   "PRIVATE",
		TokenProtected: "PROTECTED",
		TokenInternal:  "INTERNAL",
		TokenStatic:    "STATIC",
		TokenVirtual:   "VIRTUAL",
		TokenOverride:  "OVERRIDE",
		TokenAbstract:  "ABSTRACT",
		TokenSealed:    "SEALED",
		TokenPartial:   "PARTIAL",
		TokenNew:       "NEW",
		TokenThis:      "THIS",
		TokenBase:      "BASE",
		TokenAs:        "AS",
		TokenIs:        "IS",
		TokenTypeof:    "TYPEOF",
		TokenLock:      "LOCK",
		TokenUsing:     "USING",
		TokenNamespace: "NAMESPACE",
		TokenVar:       "VAR",
		TokenAsync:     "ASYNC",
		TokenAwait:     "AWAIT",
		TokenYield:     "YIELD",
		TokenLambda:    "LAMBDA",
		TokenEOF:       "EOF",
	}

	if name, ok := names[t]; ok {
		return name
	}
	return "UNKNOWN"
}
