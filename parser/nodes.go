package parser

// type VarType int

// VarType
// const (
// 	VT_String VarType = -1
// 	VT_Int    VarType = VarType(Integer)
// 	VT_Float          = -1
// )

type Operator string

const (
	Op_Plus Operator = "+"
	Op_Min           = "-"
	Op_Div           = "/"
	Op_Mul           = "*"
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

type PrintStmtNode struct {
	Inner Node
}

type FunctionCallExprNode struct {
	FuncName  string
	arguments []Node
}
