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

type Node interface{}

// Left (operator) Right
type BinaryExpressionNode struct {
	Left     Node
	Operator Operator
	Right    Node
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
	Rhs     Node
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

// (identifier)(args...)
type FunctionCallExprNode struct {
	FuncName  string
	Arguments []Node
}

type FunctionDeclNode struct {
	FuncName      string
	ArgumentNames []string
	Body          []Node
}

type ReturnNode struct {
	Expr Node
}
