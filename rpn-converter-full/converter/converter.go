package converter

import (
	"fmt"
	"rpn-converter-full/lexer"
	"strconv"
	"strings"
)

// RPNConverter преобразует токены в обратную польскую запись
type RPNConverter struct {
	output      []string
	stack       []string
	labelCount  int
	tempCount   int
	inLoop      int
	loopLabels  []string
	switchLevel int
}

// NewRPNConverter создает новый конвертер
func NewRPNConverter() *RPNConverter {
	return &RPNConverter{
		output:      make([]string, 0),
		stack:       make([]string, 0),
		labelCount:  0,
		tempCount:   0,
		inLoop:      0,
		loopLabels:  make([]string, 0),
		switchLevel: 0,
	}
}

// Convert преобразует токены в ОПЗ
func (c *RPNConverter) Convert(tokens []lexer.Token) string {
	c.output = make([]string, 0)
	c.stack = make([]string, 0)
	c.labelCount = 0
	c.tempCount = 0
	c.inLoop = 0
	c.loopLabels = make([]string, 0)
	c.switchLevel = 0

	for i := 0; i < len(tokens); i++ {
		tok := tokens[i]

		switch tok.Type {
		// Операнды
		case lexer.TokenIdent, lexer.TokenNumber, lexer.TokenString, lexer.TokenChar, lexer.TokenBool:
			c.outputToken(tok.Value)

		// Бинарные операторы
		case lexer.TokenOperator:
			// Определяем, унарный это оператор или бинарный
			isUnary := c.isUnaryOperator(tok.Value)
			c.processOperatorWithContext(tok.Value, isUnary)

		case lexer.TokenPlusAssign, lexer.TokenMinusAssign,
			lexer.TokenMulAssign, lexer.TokenDivAssign, lexer.TokenModAssign:
			c.processOperator(tok.Value)

		// Присваивание
		case lexer.TokenAssign:
			c.processOperator("=")

		// Инкремент/декремент
		case lexer.TokenInc, lexer.TokenDec:
			c.processPostfixOperator(tok.Value)

		// Скобки
		case lexer.TokenLParen:
			c.pushStack("(")
		case lexer.TokenRParen:
			c.processRParen()

		// Массивы
		case lexer.TokenLBracket:
			c.pushStack("[")
			c.pushStack("АЭМ")
			c.pushStack("2")
		case lexer.TokenRBracket:
			c.processRBracket()

		// Блоки
		case lexer.TokenLBrace:
			c.pushStack("{")
		case lexer.TokenRBrace:
			c.processRBrace()

		// Разделители
		case lexer.TokenComma:
			c.processComma()
		case lexer.TokenSemicolon:
			c.processSemicolon()
		case lexer.TokenColon:
			c.processColon()

		// Тернарный оператор
		case lexer.TokenQuestion:
			c.processOperator("?")

		// Условные операторы
		case lexer.TokenIf:
			c.pushStack("IF")
		case lexer.TokenElse:
			c.processElse()

		// Циклы
		case lexer.TokenFor:
			c.processFor(tokens, &i)
		case lexer.TokenForeach:
			c.processForeach(tokens, &i)
		case lexer.TokenWhile:
			c.processWhile(tokens, &i)
		case lexer.TokenDo:
			c.processDoWhile(tokens, &i)

		// Управляющие операторы
		case lexer.TokenBreak:
			c.processBreak()
		case lexer.TokenContinue:
			c.processContinue()
		case lexer.TokenReturn:
			c.processReturn(tokens, &i)
		case lexer.TokenGoto:
			c.processGoto(tokens, &i)

		// Исключения
		case lexer.TokenTry:
			c.processTry(tokens, &i)
		case lexer.TokenCatch:
			c.processCatch(tokens, &i)
		case lexer.TokenFinally:
			c.processFinally(tokens, &i)
		case lexer.TokenThrow:
			c.processThrow(tokens, &i)

		// Классы и ООП
		case lexer.TokenClass:
			c.processClass(tokens, &i)
		case lexer.TokenStruct:
			c.processStruct(tokens, &i)
		case lexer.TokenInterface:
			c.processInterface(tokens, &i)
		case lexer.TokenEnum:
			c.processEnum(tokens, &i)
		case lexer.TokenDelegate:
			c.processDelegate(tokens, &i)

		// Модификаторы доступа
		case lexer.TokenPublic, lexer.TokenPrivate, lexer.TokenProtected,
			lexer.TokenInternal, lexer.TokenStatic, lexer.TokenReadonly,
			lexer.TokenConst, lexer.TokenVirtual, lexer.TokenOverride,
			lexer.TokenAbstract, lexer.TokenSealed, lexer.TokenPartial:
			c.pushStack(tok.Value)

		// Специальные операторы
		case lexer.TokenNew:
			c.processNew(tokens, &i)
		case lexer.TokenThis:
			c.outputToken("THIS")
		case lexer.TokenBase:
			c.outputToken("BASE")
		case lexer.TokenAs:
			c.processOperator("as")
		case lexer.TokenIs:
			c.processOperator("is")
		case lexer.TokenTypeof:
			c.processTypeof(tokens, &i)
		case lexer.TokenSizeof:
			c.processSizeof(tokens, &i)
		case lexer.TokenChecked, lexer.TokenUnchecked:
			c.processChecked(tok.Value)

		// Асинхронность
		case lexer.TokenAsync:
			c.pushStack("ASYNC")
		case lexer.TokenAwait:
			c.processAwait(tokens, &i)

		// LINQ
		case lexer.TokenFrom, lexer.TokenSelect, lexer.TokenWhere, lexer.TokenIn:
			c.processLinq(tok.Value)

		// Лямбда
		case lexer.TokenLambda:
			c.processLambda(tokens, &i)

		// Пространства имен
		case lexer.TokenNamespace, lexer.TokenUsing:
			c.processNamespace(tokens, &i)

		// Переменные
		case lexer.TokenVar, lexer.TokenDynamic:
			// Игнорируем, тип выводится автоматически

		default:
			if tok.Value != "" {
				c.outputToken(tok.Value)
			}
		}
	}

	// Выталкиваем все оставшиеся операции
	for c.peekStack() != "" {
		if c.peekStack() != "(" && c.peekStack() != "[" && c.peekStack() != "{" {
			c.outputToken(c.popStack())
		} else {
			c.popStack()
		}
	}

	return strings.Join(c.output, " ")
}

// isUnaryOperator определяет, является ли оператор унарным в текущем контексте
func (c *RPNConverter) isUnaryOperator(op string) bool {
	// Проверяем, является ли оператор потенциально унарным
	if op != "+" && op != "-" && op != "!" && op != "~" {
		return false
	}

	// Если стек пуст или последний элемент в выходе - оператор или открывающая скобка,
	// то это унарный оператор
	if len(c.output) == 0 {
		return true
	}

	lastOutput := c.output[len(c.output)-1]
	// Если последний элемент - оператор или открывающая скобка
	if lastOutput == "(" || lastOutput == "[" || lastOutput == "{" ||
		IsBinaryOperator(NormalizeOperator(lastOutput, false)) {
		return true
	}

	return false
}

func (c *RPNConverter) processOperatorWithContext(op string, isUnary bool) {
	normalizedOp := NormalizeOperator(op, isUnary)

	for {
		top := c.peekStack()
		if top == "" {
			break
		}

		topPriority := GetOperatorPriorityWithContext(top, IsUnaryOperator(top))
		opPriority := GetOperatorPriorityWithContext(normalizedOp, isUnary)

		if topPriority > opPriority ||
			(topPriority == opPriority && !IsRightAssoc(op) && top != "(" && top != "[") {
			c.outputToken(c.popStack())
		} else {
			break
		}
	}
	c.pushStack(normalizedOp)
}

func (c *RPNConverter) processOperator(op string) {
	c.processOperatorWithContext(op, false)
}

func (c *RPNConverter) processPostfixOperator(op string) {
	// Постфиксный оператор: a++ -> a ++
	temp := c.newTemp()
	c.outputToken(temp)
	c.outputToken(":=")
	c.outputToken(op)
	c.pushStack(temp)
}

func (c *RPNConverter) processRParen() {
	for c.peekStack() != "" && c.peekStack() != "(" {
		c.outputToken(c.popStack())
	}
	if c.peekStack() == "(" {
		c.popStack()
	}
}

func (c *RPNConverter) processRBracket() {
	for c.peekStack() != "[" && c.peekStack() != "" {
		c.outputToken(c.popStack())
	}
	if c.peekStack() == "[" {
		c.popStack()
	}
	if c.peekStack() == "АЭМ" {
		c.popStack()
		counter := c.popStack()
		c.outputToken(counter)
		c.outputToken("АЭМ")
	}
}

func (c *RPNConverter) processRBrace() {
	for c.peekStack() != "" && c.peekStack() != "{" {
		c.outputToken(c.popStack())
	}
	if c.peekStack() == "{" {
		c.popStack()
	}
}

func (c *RPNConverter) processComma() {
	for c.peekStack() != "" && c.peekStack() != "(" && c.peekStack() != "[" &&
		c.peekStack() != "АЭМ" && c.peekStack() != "Ф" {
		c.outputToken(c.popStack())
	}

	if c.peekStack() == "АЭМ" {
		c.popStack()
		counter := c.popStack()
		newCounter := fmt.Sprintf("%d", atoi(counter)+1)
		c.pushStack("АЭМ")
		c.pushStack(newCounter)
	} else if c.peekStack() == "Ф" {
		c.popStack()
		counter := c.popStack()
		newCounter := fmt.Sprintf("%d", atoi(counter)+1)
		c.pushStack("Ф")
		c.pushStack(newCounter)
	}
}

func (c *RPNConverter) processSemicolon() {
	for c.peekStack() != "" && c.peekStack() != "(" && c.peekStack() != "[" && c.peekStack() != "{" {
		c.outputToken(c.popStack())
	}
}

func (c *RPNConverter) processColon() {
	if c.peekStack() == "?" {
		c.popStack()
		c.processTernary()
	} else {
		c.processOperator(":")
	}
}

func (c *RPNConverter) processTernary() {
	label1 := c.newLabel()
	label2 := c.newLabel()

	condition := c.popStack()
	c.outputToken(condition)
	c.outputToken(label1)
	c.outputToken("УПЛ")

	trueExpr := c.popStack()
	c.outputToken(trueExpr)
	c.outputToken(label2)
	c.outputToken("БП")
	c.outputToken(label1)
	c.outputToken(":")

	falseExpr := c.popStack()
	c.outputToken(falseExpr)
	c.outputToken(label2)
	c.outputToken(":")

	c.pushStack(trueExpr)
}

func (c *RPNConverter) processElse() {
	// Выталкиваем все до IF
	for c.peekStack() != "" && !strings.HasPrefix(c.peekStack(), "IF") {
		c.outputToken(c.popStack())
	}
	if strings.HasPrefix(c.peekStack(), "IF") {
		ifLabel := c.popStack()
		label2 := c.newLabel()
		c.outputToken(label2)
		c.outputToken("БП")
		c.outputToken(ifLabel[2:]) // извлекаем метку
		c.pushStack("IF" + label2)
	}
}

func (c *RPNConverter) processFor(tokens []lexer.Token, i *int) {
	c.pushStack("FOR")
	c.inLoop++

	labelStart := c.newLabel()
	labelEnd := c.newLabel()
	c.loopLabels = append(c.loopLabels, labelStart, labelEnd)

	// Пропускаем "for"
	*i++

	// Инициализация
	if *i < len(tokens) && tokens[*i].Type == lexer.TokenLParen {
		*i++
		// Парсим инициализацию
		for *i < len(tokens) && tokens[*i].Type != lexer.TokenSemicolon {
			c.processToken(tokens[*i])
			*i++
		}
		if *i < len(tokens) && tokens[*i].Type == lexer.TokenSemicolon {
			c.processSemicolon()
			*i++
		}

		// Условие
		c.outputToken(labelStart)
		c.outputToken(":")

		conditionStart := *i
		for *i < len(tokens) && tokens[*i].Type != lexer.TokenSemicolon {
			c.processToken(tokens[*i])
			*i++
		}
		if conditionStart < *i {
			c.outputToken(labelEnd)
			c.outputToken("УПЛ")
		}

		if *i < len(tokens) && tokens[*i].Type == lexer.TokenSemicolon {
			*i++
		}

		// Инкремент
		for *i < len(tokens) && tokens[*i].Type != lexer.TokenRParen {
			c.processToken(tokens[*i])
			*i++
		}

		if *i < len(tokens) && tokens[*i].Type == lexer.TokenRParen {
			*i++
		}
	}

	c.pushStack("FOR_BODY" + labelStart + ":" + labelEnd)
}

func (c *RPNConverter) processForeach(tokens []lexer.Token, i *int) {
	c.pushStack("FOREACH")
	c.inLoop++

	labelStart := c.newLabel()
	labelEnd := c.newLabel()
	c.loopLabels = append(c.loopLabels, labelStart, labelEnd)

	// Пропускаем "foreach"
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenLParen {
		*i++
		// Переменная
		if *i < len(tokens) && tokens[*i].Type == lexer.TokenVar {
			*i++
		}
		if *i < len(tokens) && tokens[*i].Type == lexer.TokenIdent {
			itemVar := tokens[*i].Value
			*i++
			if *i < len(tokens) && tokens[*i].Type == lexer.TokenIn {
				*i++
			}
			// Коллекция
			for *i < len(tokens) && tokens[*i].Type != lexer.TokenRParen {
				c.processToken(tokens[*i])
				*i++
			}
			c.outputToken("GET_ENUMERATOR")
			c.outputToken(labelStart)
			c.outputToken(":")
			c.outputToken("MOVE_NEXT")
			c.outputToken(labelEnd)
			c.outputToken("УПЛ")
			c.outputToken("GET_CURRENT")
			c.outputToken(itemVar)
			c.outputToken(":=")
		}

		if *i < len(tokens) && tokens[*i].Type == lexer.TokenRParen {
			*i++
		}
	}

	c.pushStack("FOREACH_BODY" + labelStart + ":" + labelEnd)
}

func (c *RPNConverter) processWhile(tokens []lexer.Token, i *int) {
	c.pushStack("WHILE")
	c.inLoop++

	labelStart := c.newLabel()
	labelEnd := c.newLabel()
	c.loopLabels = append(c.loopLabels, labelStart, labelEnd)

	// Пропускаем "while"
	*i++

	c.outputToken(labelStart)
	c.outputToken(":")

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenLParen {
		*i++
		// Парсим условие
		for *i < len(tokens) && tokens[*i].Type != lexer.TokenRParen {
			c.processToken(tokens[*i])
			*i++
		}
		c.outputToken(labelEnd)
		c.outputToken("УПЛ")

		if *i < len(tokens) && tokens[*i].Type == lexer.TokenRParen {
			*i++
		}
	}

	c.pushStack("WHILE_BODY" + labelStart + ":" + labelEnd)
}

func (c *RPNConverter) processDoWhile(tokens []lexer.Token, i *int) {
	c.pushStack("DO")
	c.inLoop++

	labelStart := c.newLabel()
	labelEnd := c.newLabel()
	c.loopLabels = append(c.loopLabels, labelStart, labelEnd)

	// Пропускаем "do"
	*i++

	c.outputToken(labelStart)
	c.outputToken(":")

	c.pushStack("DO_BODY" + labelStart + ":" + labelEnd)

	// Ищем while
	for *i < len(tokens) && tokens[*i].Type != lexer.TokenWhile {
		*i++
	}

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenWhile {
		*i++
		if *i < len(tokens) && tokens[*i].Type == lexer.TokenLParen {
			*i++
			for *i < len(tokens) && tokens[*i].Type != lexer.TokenRParen {
				c.processToken(tokens[*i])
				*i++
			}
			c.outputToken(labelStart)
			c.outputToken("БП")
			c.outputToken(labelEnd)
			c.outputToken(":")

			if *i < len(tokens) && tokens[*i].Type == lexer.TokenRParen {
				*i++
			}
		}
	}
}

func (c *RPNConverter) processBreak() {
	if c.inLoop > 0 && len(c.loopLabels) >= 2 {
		// Выход из цикла
		c.outputToken(c.loopLabels[len(c.loopLabels)-1])
		c.outputToken("БП")
	} else if c.switchLevel > 0 {
		// Выход из switch
		c.outputToken("BREAK")
		c.outputToken("SWITCH_END")
	}
}

func (c *RPNConverter) processContinue() {
	if c.inLoop > 0 && len(c.loopLabels) >= 2 {
		// Переход к следующей итерации
		c.outputToken(c.loopLabels[len(c.loopLabels)-2])
		c.outputToken("БП")
	}
}

func (c *RPNConverter) processReturn(tokens []lexer.Token, i *int) {
	c.outputToken("RETURN")
	*i++

	// Парсим возвращаемое значение
	for *i < len(tokens) && tokens[*i].Type != lexer.TokenSemicolon {
		c.processToken(tokens[*i])
		*i++
	}
}

func (c *RPNConverter) processGoto(tokens []lexer.Token, i *int) {
	*i++
	if *i < len(tokens) && tokens[*i].Type == lexer.TokenIdent {
		c.outputToken(tokens[*i].Value)
		c.outputToken("БП")
		*i++
	}
}

func (c *RPNConverter) processTry(tokens []lexer.Token, i *int) {
	c.outputToken("TRY")
	c.pushStack("TRY")
	*i++
}

func (c *RPNConverter) processCatch(tokens []lexer.Token, i *int) {
	c.outputToken("CATCH")
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenLParen {
		*i++
		// Тип исключения
		for *i < len(tokens) && tokens[*i].Type != lexer.TokenRParen {
			c.processToken(tokens[*i])
			*i++
		}
		if *i < len(tokens) && tokens[*i].Type == lexer.TokenRParen {
			*i++
		}
	}
}

func (c *RPNConverter) processFinally(tokens []lexer.Token, i *int) {
	c.outputToken("FINALLY")
	*i++
}

func (c *RPNConverter) processThrow(tokens []lexer.Token, i *int) {
	c.outputToken("THROW")
	*i++

	// Парсим выражение исключения
	for *i < len(tokens) && tokens[*i].Type != lexer.TokenSemicolon {
		c.processToken(tokens[*i])
		*i++
	}
}

func (c *RPNConverter) processClass(tokens []lexer.Token, i *int) {
	c.outputToken("CLASS")
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenIdent {
		c.outputToken(tokens[*i].Value)
		c.outputToken("DEFINE_CLASS")
		*i++
	}

	c.pushStack("CLASS")
}

func (c *RPNConverter) processStruct(tokens []lexer.Token, i *int) {
	c.outputToken("STRUCT")
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenIdent {
		c.outputToken(tokens[*i].Value)
		c.outputToken("DEFINE_STRUCT")
		*i++
	}
}

func (c *RPNConverter) processInterface(tokens []lexer.Token, i *int) {
	c.outputToken("INTERFACE")
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenIdent {
		c.outputToken(tokens[*i].Value)
		c.outputToken("DEFINE_INTERFACE")
		*i++
	}
}

func (c *RPNConverter) processEnum(tokens []lexer.Token, i *int) {
	c.outputToken("ENUM")
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenIdent {
		c.outputToken(tokens[*i].Value)
		c.outputToken("DEFINE_ENUM")
		*i++
	}
}

func (c *RPNConverter) processDelegate(tokens []lexer.Token, i *int) {
	c.outputToken("DELEGATE")
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenIdent {
		c.outputToken(tokens[*i].Value)
		c.outputToken("DEFINE_DELEGATE")
		*i++
	}
}

func (c *RPNConverter) processNew(tokens []lexer.Token, i *int) {
	c.outputToken("NEW")
	*i++

	// Тип
	if *i < len(tokens) && tokens[*i].Type == lexer.TokenIdent {
		typeName := tokens[*i].Value
		c.outputToken(typeName)
		*i++

		// Аргументы конструктора
		if *i < len(tokens) && tokens[*i].Type == lexer.TokenLParen {
			c.pushStack("(")
			*i++
			argCount := 0

			for *i < len(tokens) && tokens[*i].Type != lexer.TokenRParen {
				c.processToken(tokens[*i])
				argCount++
				if *i < len(tokens) && tokens[*i].Type == lexer.TokenComma {
					*i++
				}
			}

			c.processRParen()
			c.outputToken(fmt.Sprintf("%d", argCount))
			c.outputToken("CONSTRUCTOR_CALL")

			if *i < len(tokens) && tokens[*i].Type == lexer.TokenRParen {
				*i++
			}
		}
	}
}

func (c *RPNConverter) processTypeof(tokens []lexer.Token, i *int) {
	c.outputToken("TYPEOF")
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenLParen {
		*i++
		for *i < len(tokens) && tokens[*i].Type != lexer.TokenRParen {
			c.processToken(tokens[*i])
			*i++
		}
		if *i < len(tokens) && tokens[*i].Type == lexer.TokenRParen {
			*i++
		}
	}
}

func (c *RPNConverter) processSizeof(tokens []lexer.Token, i *int) {
	c.outputToken("SIZEOF")
	*i++

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenLParen {
		*i++
		for *i < len(tokens) && tokens[*i].Type != lexer.TokenRParen {
			c.processToken(tokens[*i])
			*i++
		}
		if *i < len(tokens) && tokens[*i].Type == lexer.TokenRParen {
			*i++
		}
	}
}

func (c *RPNConverter) processChecked(operator string) {
	c.outputToken(operator)
	c.pushStack("CHECKED")
}

func (c *RPNConverter) processAwait(tokens []lexer.Token, i *int) {
	c.outputToken("AWAIT")
	*i++
	if *i < len(tokens) {
		c.processToken(tokens[*i])
	}
}

func (c *RPNConverter) processLinq(keyword string) {
	c.outputToken(strings.ToUpper(keyword))
}

func (c *RPNConverter) processLambda(tokens []lexer.Token, i *int) {
	c.outputToken("LAMBDA")
	*i++

	// Параметры
	for *i < len(tokens) && tokens[*i].Type != lexer.TokenLambda {
		if tokens[*i].Type == lexer.TokenIdent {
			c.outputToken(tokens[*i].Value)
		}
		*i++
	}

	if *i < len(tokens) && tokens[*i].Type == lexer.TokenLambda {
		*i++
	}

	// Тело лямбды
	c.pushStack("LAMBDA_BODY")
}

func (c *RPNConverter) processNamespace(tokens []lexer.Token, i *int) {
	c.outputToken(strings.ToUpper(tokens[*i].Value))
	*i++

	for *i < len(tokens) && tokens[*i].Type != lexer.TokenLBrace &&
		tokens[*i].Type != lexer.TokenSemicolon {
		c.processToken(tokens[*i])
		*i++
	}
}

func (c *RPNConverter) processToken(tok lexer.Token) {
	switch tok.Type {
	case lexer.TokenIdent, lexer.TokenNumber, lexer.TokenString, lexer.TokenChar:
		c.outputToken(tok.Value)
	case lexer.TokenOperator:
		isUnary := c.isUnaryOperator(tok.Value)
		c.processOperatorWithContext(tok.Value, isUnary)
	case lexer.TokenLParen:
		c.pushStack("(")
	case lexer.TokenRParen:
		c.processRParen()
	default:
		if tok.Value != "" {
			c.outputToken(tok.Value)
		}
	}
}

func (c *RPNConverter) outputToken(token string) {
	c.output = append(c.output, token)
}

func (c *RPNConverter) pushStack(val string) {
	c.stack = append(c.stack, val)
}

func (c *RPNConverter) popStack() string {
	if len(c.stack) == 0 {
		return ""
	}
	val := c.stack[len(c.stack)-1]
	c.stack = c.stack[:len(c.stack)-1]
	return val
}

func (c *RPNConverter) peekStack() string {
	if len(c.stack) == 0 {
		return ""
	}
	return c.stack[len(c.stack)-1]
}

func (c *RPNConverter) newLabel() string {
	c.labelCount++
	return fmt.Sprintf("L%d", c.labelCount)
}

func (c *RPNConverter) newTemp() string {
	c.tempCount++
	return fmt.Sprintf("T%d", c.tempCount)
}

func atoi(s string) int {
	result, _ := strconv.Atoi(s)
	return result
}
