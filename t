parsing myVar
>> = Ident
Parsing expresssion
&parser.Token{Type:3, Value:"5"}
parsing 

parsing myVar
>> = Ident
Parsing expresssion
&parser.Token{Type:10, Value:"("}
Parsing open par
token::: &{10 (}
bracket level: 0
Collected inner token: &{0 myVar}
Collected inner token: &{6 +}
Collected inner token: &{3 6}
After collecting inner brackets, iter now at &parser.Token{Type:9, Value:"*"}
Parsing expresssion
&parser.Token{Type:0, Value:"myVar"}
Expression length: 3
nextToken: &{0 myVar}
nextToken: &{6 +}
Found operator
[0xc0000102e8]
Parsing expresssion
&parser.Token{Type:0, Value:"myVar"}
Expresion: []interface {}{parser.VariableNode{Name:"myVar"}, (*parser.Token)(0xc000010300)}
nextToken: &{3 6}
nextToken: <nil>
End? false
lhs end: [0xc000010318 0xc000010468]
Parsing expresssion
&parser.Token{Type:3, Value:"6"}
Final Expresion: []interface {}{parser.VariableNode{Name:"myVar"}, (*parser.Token)(0xc000010300), parser.IntegerLiteralNode{Value:6}}
At 0 {myVar}
At 1 &{6 +}
is token
At 2 {6}
Checking operator 7
At 0 {myVar}
At 1 &{6 +}
is token
Expression length: 18
nextToken: &{9 *}
Found operator
[]
Expresion: []interface {}{parser.BinaryExpressionNode{Left:parser.VariableNode{Name:"myVar"}, Operator:"+", Right:parser.IntegerLiteralNode{Value:6}}, (*parser.Token)(0xc0000104f8)}
nextToken: &{3 2}
nextToken: &{5 
}
nextToken: &{4 noot!}
nextToken: &{10 (}
nextToken: &{0 myVar}
nextToken: &{11 )}
nextToken: <nil>
End? false
lhs end: [0xc000010570 0xc0000105a0 0xc0000105d0 0xc000010600 0xc000010630 0xc000010660 0xc000010690]
Parsing expresssion
&parser.Token{Type:3, Value:"2"}
Final Expresion: []interface {}{parser.BinaryExpressionNode{Left:parser.VariableNode{Name:"myVar"}, Operator:"+", Right:parser.IntegerLiteralNode{Value:6}}, (*parser.Token)(0xc0000104f8), parser.IntegerLiteralNode{Value:2}}
At 0 {{myVar} + {6}}
At 1 &{9 *}
is token
At 2 {2}
Checking operator 7
At 0 {{myVar} + {6}}
At 1 &{9 *}
is token
At 2 {2}
Checking operator 6
At 0 {{myVar} + {6}}
At 1 &{9 *}
is token
At 2 {2}
Checking operator 8
At 0 {{myVar} + {6}}
At 1 &{9 *}
is token
Tree: []parser.Node{parser.VarDeclNode{VarName:"myVar", Rhs:parser.IntegerLiteralNode{Value:5}}, parser.VarAssignNode{VarName:"myVar", Rhs:parser.BinaryExpressionNode{Left:parser.BinaryExpressionNode{Left:parser.VariableNode{Name:"myVar"}, Operator:"+", Right:parser.IntegerLiteralNode{Value:6}}, Operator:"*", Right:parser.IntegerLiteralNode{Value:2}}}}
