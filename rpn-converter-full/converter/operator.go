package converter

// OperatorPriority возвращает приоритет оператора
func OperatorPriority(op string) int {
	priority := map[string]int{
		// Присваивания (самый низкий приоритет)
		"=": 1, "+=": 1, "-=": 1, "*=": 1, "/=": 1, "%=": 1,
		"&=": 1, "|=": 1, "^=": 1, "<<=": 1, ">>=": 1,

		// Тернарный оператор
		"?": 2, ":": 2,

		// Null-coalescing
		"??": 3,

		// Логические операции
		"||": 4,
		"&&": 5,

		// Битовые операции
		"|": 6,
		"^": 7,
		"&": 8,

		// Операции сравнения
		"==": 9, "!=": 9,
		"<": 10, ">": 10, "<=": 10, ">=": 10,
		"is": 10, "as": 10,

		// Сдвиги
		"<<": 11, ">>": 11,

		// Сложение/вычитание (бинарные)
		"bin+": 12, "bin-": 12,

		// Умножение/деление
		"*": 13, "/": 13, "%": 13,

		// Возведение в степень
		"**": 14,

		// Доступ к членам
		".": 15, "->": 15,

		// Вызов функций и индексаторы
		"()": 16, "[]": 16,
	}

	if p, ok := priority[op]; ok {
		return p
	}
	return 0
}

// GetUnaryOperatorPriority возвращает приоритет унарного оператора
func GetUnaryOperatorPriority(op string) int {
	unaryPriority := map[string]int{
		// Унарные операторы (самый высокий приоритет)
		"++": 17, "--": 17, "!": 17, "~": 17, "unary+": 17, "unary-": 17,
	}

	if p, ok := unaryPriority[op]; ok {
		return p
	}
	return 0
}

// IsRightAssoc проверяет, является ли оператор правоассоциативным
func IsRightAssoc(op string) bool {
	rightAssoc := map[string]bool{
		"=": true, "+=": true, "-=": true, "*=": true, "/=": true, "%=": true,
		"&=": true, "|=": true, "^=": true, "<<=": true, ">>=": true,
		"?": true, "??": true,
		"**": true,
		"!":  true, "~": true, "++": true, "--": true,
	}

	return rightAssoc[op]
}

// IsUnaryOperator проверяет, является ли оператор унарным
func IsUnaryOperator(op string) bool {
	unary := map[string]bool{
		"unary+": true, "unary-": true, "!": true, "~": true,
		"++": true, "--": true,
	}

	return unary[op]
}

// IsBinaryOperator проверяет, является ли оператор бинарным
func IsBinaryOperator(op string) bool {
	binary := map[string]bool{
		"bin+": true, "bin-": true, "*": true, "/": true, "%": true, "**": true,
		"==": true, "!=": true, "<": true, ">": true, "<=": true, ">=": true,
		"&&": true, "||": true, "&": true, "|": true, "^": true,
		"<<": true, ">>": true,
		"=": true, "+=": true, "-=": true, "*=": true, "/=": true,
		"?": true, "??": true, "is": true, "as": true,
	}

	return binary[op]
}

// NormalizeOperator нормализует оператор для таблицы приоритетов
func NormalizeOperator(op string, isUnary bool) string {
	if isUnary {
		switch op {
		case "+":
			return "unary+"
		case "-":
			return "unary-"
		default:
			return op
		}
	}

	switch op {
	case "+":
		return "bin+"
	case "-":
		return "bin-"
	default:
		return op
	}
}

// GetOperatorPriorityWithContext возвращает приоритет оператора с учетом контекста
func GetOperatorPriorityWithContext(op string, isUnary bool) int {
	if isUnary {
		return GetUnaryOperatorPriority(op)
	}
	return OperatorPriority(NormalizeOperator(op, false))
}
