package models

// ASTNode представляет узел абстрактного синтаксического дерева
type ASTNode interface {
	GetType() string
}

// ProgramNode представляет программу
type ProgramNode struct {
	Classes    []*ClassNode
	Statements []StatementNode
}

func (p *ProgramNode) GetType() string { return "Program" }

// ClassNode представляет класс
type ClassNode struct {
	Name       string
	Modifiers  []string
	BaseClass  string
	Interfaces []string
	Fields     []*FieldNode
	Properties []*PropertyNode
	Methods    []*MethodNode
	IsPartial  bool
	IsStatic   bool
	IsAbstract bool
	IsSealed   bool
}

func (c *ClassNode) GetType() string { return "Class" }

// FieldNode представляет поле класса
type FieldNode struct {
	Name      string
	Type      string
	Modifiers []string
	IsStatic  bool
	IsConst   bool
	InitValue ExpressionNode
}

func (f *FieldNode) GetType() string { return "Field" }

// PropertyNode представляет свойство
type PropertyNode struct {
	Name      string
	Type      string
	Modifiers []string
	HasGet    bool
	HasSet    bool
	HasInit   bool
	IsAuto    bool
	GetBody   []StatementNode
	SetBody   []StatementNode
}

func (p *PropertyNode) GetType() string { return "Property" }

// MethodNode представляет метод
type MethodNode struct {
	Name       string
	ReturnType string
	Parameters []*ParameterNode
	Modifiers  []string
	Body       []StatementNode
	IsVirtual  bool
	IsOverride bool
	IsAbstract bool
	IsStatic   bool
}

func (m *MethodNode) GetType() string { return "Method" }

// ParameterNode представляет параметр метода
type ParameterNode struct {
	Name string
	Type string
}

// StatementNode представляет оператор
type StatementNode interface {
	ASTNode
	IsStatement()
}

// ExpressionStatement представляет выражение как оператор
type ExpressionStatement struct {
	Expr ExpressionNode
}

func (e *ExpressionStatement) GetType() string { return "ExpressionStatement" }
func (e *ExpressionStatement) IsStatement()    {}

// AssignmentStatement представляет оператор присваивания
type AssignmentStatement struct {
	Left  ExpressionNode
	Right ExpressionNode
}

func (a *AssignmentStatement) GetType() string { return "Assignment" }
func (a *AssignmentStatement) IsStatement()    {}

// IfStatement представляет условный оператор
type IfStatement struct {
	Condition ExpressionNode
	ThenBody  []StatementNode
	ElseBody  []StatementNode
}

func (i *IfStatement) GetType() string { return "If" }
func (i *IfStatement) IsStatement()    {}

// ForStatement представляет цикл for
type ForStatement struct {
	Init      StatementNode
	Condition ExpressionNode
	Post      StatementNode
	Body      []StatementNode
}

func (f *ForStatement) GetType() string { return "For" }
func (f *ForStatement) IsStatement()    {}

// WhileStatement представляет цикл while
type WhileStatement struct {
	Condition ExpressionNode
	Body      []StatementNode
}

func (w *WhileStatement) GetType() string { return "While" }
func (w *WhileStatement) IsStatement()    {}

// ReturnStatement представляет оператор return
type ReturnStatement struct {
	Value ExpressionNode
}

func (r *ReturnStatement) GetType() string { return "Return" }
func (r *ReturnStatement) IsStatement()    {}

// ExpressionNode представляет выражение
type ExpressionNode interface {
	ASTNode
	IsExpression()
}

// BinaryExpression представляет бинарное выражение
type BinaryExpression struct {
	Left     ExpressionNode
	Operator string
	Right    ExpressionNode
}

func (b *BinaryExpression) GetType() string { return "Binary" }
func (b *BinaryExpression) IsExpression()   {}

// UnaryExpression представляет унарное выражение
type UnaryExpression struct {
	Operator string
	Operand  ExpressionNode
}

func (u *UnaryExpression) GetType() string { return "Unary" }
func (u *UnaryExpression) IsExpression()   {}

// Identifier представляет идентификатор
type Identifier struct {
	Name string
}

func (i *Identifier) GetType() string { return "Identifier" }
func (i *Identifier) IsExpression()   {}

// Literal представляет литерал (число, строка)
type Literal struct {
	Type  string
	Value string
}

func (l *Literal) GetType() string { return "Literal" }
func (l *Literal) IsExpression()   {}

// MemberAccess представляет доступ к члену (obj.member)
type MemberAccess struct {
	Object ExpressionNode
	Member string
}

func (m *MemberAccess) GetType() string { return "MemberAccess" }
func (m *MemberAccess) IsExpression()   {}

// ArrayAccess представляет доступ к элементу массива
type ArrayAccess struct {
	Array ExpressionNode
	Index ExpressionNode
}

func (a *ArrayAccess) GetType() string { return "ArrayAccess" }
func (a *ArrayAccess) IsExpression()   {}

// MethodCall представляет вызов метода
type MethodCall struct {
	Method    ExpressionNode
	Arguments []ExpressionNode
}

func (m *MethodCall) GetType() string { return "MethodCall" }
func (m *MethodCall) IsExpression()   {}
