package parser

type Operator string

const (
	Op_Plus       Operator = "+"
	Op_Min                 = "-"
	Op_Div                 = "/"
	Op_Mul                 = "*"
	Op_CompEqual           = "=="
	Op_CompNEqual          = "!="
	Op_LT                  = "<"
	Op_GT                  = ">"
	Op_LTE                 = "<="
	Op_GTE                 = ">="
	Op_Or                  = "||"
	Op_And                 = "&&"
)

const (
	Op_Equal      Operator = "="
	Op_PlusEqual           = "+="
	Op_MinEqual            = "-="
	Op_TimesEqual          = "*="
	Op_DivEqual            = "/="
)

type Node interface{}

// Left (operator) Right
type BinaryExpressionNode struct {
	Left     Node
	Operator Operator
	Right    Node
}

//
type IfNode struct {
	Condition Node
	// Can be nil if no more else block
	NextElseBlock Node
	Body          []Node
}

type ElseNode struct {
	Body []Node
}

// (identifier)
type VariableNode struct {
	Name string
}

// VarName := Rhs
type VarDeclNode struct {
	VarName string
	Rhs     Node
}

// VarName = Rhs
type VarAssignNode struct {
	VarName string
	Op      Operator
	Rhs     Node
}

type ArrayIndexAssignmentNode struct {
	Array VariableNode
	Index Node
	Rhs   Node
}

// (int)
type IntegerLiteralNode struct {
	Value int64
}

type NilLiteralNode struct{}

type StringLiteralNode struct {
	String string
}

type FloatLiteralNode struct {
	Value float64
}

type BoolLiteralNode struct {
	Value bool
}

// [(expr,)*]
type ArrayLiteralNode struct {
	Values []Node
}

// !(expr)
type BinaryNotNode struct {
	Expr Node
}

// (identifier)(args...)
type FunctionCallExprNode struct {
	FuncName  string
	Arguments []Node
}

type MethodCallExprNode struct {
	CalledOn     Node
	FunctionCall FunctionCallExprNode
}

type FunctionDeclNode struct {
	FuncName      string
	ArgumentNames []string
	Body          []Node
}

// (ident)[(int)]
type ArrayIndexNode struct {
	Array Node
	Index Node
}

type ReturnNode struct {
	Expr Node
}

type WhileNode struct {
	Condition Node
	Body      []Node
}
