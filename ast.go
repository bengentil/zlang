// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	"strconv"
)

//var NamedValues map[string]*llvm.Value

type Node interface {
	CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)
}

// ********************************************
// Expressions
// ********************************************
type (
	NodeExpr interface {
		Node
		exprNode()
	}

	NodeInteger struct {
		Node
		NodeExpr
		Value int
	}

	NodeFloat struct {
		Node
		NodeExpr
		Value float32
	}

	NodeString struct {
		Node
		NodeExpr
		Value string
	}

	NodeIdentifier struct {
		Node
		NodeExpr
		Name string
	}

	NodeBinOperator struct {
		Node
		NodeExpr
		Operator string
		LHS      NodeExpr
		RHS      NodeExpr
	}

	NodeAssignement struct {
		Node
		NodeExpr
		LHS *NodeIdentifier
		RHS NodeExpr
	}

	NodeCall struct {
		Node
		NodeExpr
		Name *NodeIdentifier
		Args []NodeExpr
	}

	NodeBlock struct {
		Node
		NodeExpr
		Statements []NodeStmt
		Depth      int
	}
)

//
// Make nodes creation easier
//

func NInteger(value string) *NodeInteger {
	val, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &NodeInteger{Value: val}
}

func NFloat(value string) *NodeFloat {
	return &NodeFloat{Value: 0.0}
}

func NString(value string) *NodeString {
	return &NodeString{Value: value}
}

func NIdentifier(name string) *NodeIdentifier {
	return &NodeIdentifier{Name: name}
}

func NBinOp(operator string, lhs, rhs NodeExpr) *NodeBinOperator {
	return &NodeBinOperator{Operator: operator, LHS: lhs, RHS: rhs}
}

func NAssignement(lhs *NodeIdentifier, rhs NodeExpr) *NodeAssignement {
	return &NodeAssignement{LHS: lhs, RHS: rhs}
}

func NCall(name *NodeIdentifier, args []NodeExpr) *NodeCall {
	return &NodeCall{Name: name, Args: args}
}

func NBlock(stmts []NodeStmt, depth int) *NodeBlock {
	return &NodeBlock{Statements: stmts, Depth: depth}
}

//
// Make nodes printable
//
func (n *NodeInteger) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeInteger\",\"value\":%v}", n.Value))
}

func (n *NodeFloat) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeFloat\",\"value\":%v}", n.Value))
}

func (n *NodeString) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeString\",\"name\":%q}", n.Value))
}

func (n *NodeIdentifier) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeIdentifier\",\"name\":%q}", n.Name))
}

func (n *NodeBinOperator) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeBinOperator\",\"op\":%q, \"lhs\":%v, \"rhs\":%v}", n.Operator, n.LHS, n.RHS))
}

func (n *NodeAssignement) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeAssignement\",\"lhs\":%v, \"rhs\":%v}", n.LHS, n.RHS))
}

func (n *NodeCall) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeCall\",\"name\":%v, \"args\":%v}", n.Name, n.Args))
}

func (n *NodeBlock) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeBlock\",\"statements\":%v,\"depth\":%d}", n.Statements, n.Depth))
}

// exprNode() ensures that only statement nodes can be
// assigned to a NodeExpr.
func (n *NodeInteger) exprNode()     {}
func (n *NodeFloat) exprNode()       {}
func (n *NodeIdentifier) exprNode()  {}
func (n *NodeBinOperator) exprNode() {}
func (n *NodeAssignement) exprNode() {}
func (n *NodeCall) exprNode()        {}
func (n *NodeBlock) exprNode()       {}

// ********************************************
// Statements
// ********************************************
type (
	NodeStmt interface {
		Node
		stmtNode()
	}

	NodePrototype struct {
		Node
		NodeStmt
		Type *NodeIdentifier
		Name *NodeIdentifier
		Args []*NodeVariable
	}

	NodeFunction struct {
		Node
		NodeStmt
		Proto *NodePrototype
		Body  *NodeBlock
	}

	NodeReturn struct {
		Node
		NodeStmt
		Value NodeExpr
	}

	NodeExpression struct {
		Node
		NodeStmt
		Expression NodeExpr
	}

	NodeVariable struct {
		Node
		NodeStmt
		Type       *NodeIdentifier
		Name       *NodeIdentifier
		AssignExpr *NodeAssignement
	}
)

//
// Make nodes creation easier
//
func NPrototype(name, typ *NodeIdentifier, args []*NodeVariable) *NodePrototype {
	return &NodePrototype{Name: name, Type: typ, Args: args}
}

func NFunction(proto *NodePrototype, body *NodeBlock) *NodeFunction {
	return &NodeFunction{Proto: proto, Body: body}
}

func NReturn(exp NodeExpr) *NodeReturn {
	return &NodeReturn{Value: exp}
}

func NExpression(exp NodeExpr) *NodeExpression {
	return &NodeExpression{Expression: exp}
}

func NVariable(name, typ *NodeIdentifier, assign *NodeAssignement) *NodeVariable {
	return &NodeVariable{Name: name, Type: typ, AssignExpr: assign}
}

//
// Make nodes printable
//

func (n *NodePrototype) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodePrototype\",\"name\":%v, \"type\":%v, \"args\":%v}", n.Name, n.Type, n.Args))
}

func (n *NodeFunction) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeFunction\",\"proto\":%v,\"body\":%v}", n.Proto, n.Body))
}

func (n *NodeReturn) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeReturn\",\"value\":%v}", n.Value))
}

func (n *NodeExpression) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeExpression\",\"expression\":%v}", n.Expression))
}

func (n *NodeVariable) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeVariable\",\"name\":%v,\"type\":%v,\"assign_expr\":%v}", n.Name, n.Type, n.AssignExpr))
}

// stmtNode() ensures that only statement nodes can be
// assigned to a NodeStmt.
func (n *NodePrototype) stmtNode()  {}
func (n *NodeFunction) stmtNode()   {}
func (n *NodeReturn) stmtNode()     {}
func (n *NodeExpression) stmtNode() {}
func (n *NodeVariable) stmtNode()   {}

// ********************************************
// Code generation
// ********************************************
func (n *NodeInteger) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)    { return nil, nil }
func (n *NodeFloat) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)      { return nil, nil }
func (n *NodeString) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)     { return nil, nil }
func (n *NodeIdentifier) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }
func (n *NodeBinOperator) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	l, err := n.LHS.CodeGen()
	r, err := n.RHS.CodeGen()

	if l == nil || r == nil || err != nil {
		return nil, err
	}

	switch n.Operator {
	case "*":
		return builder.CreateMul(l, r, "multmp"), nil
	}
	return nil, nil
}
func (n *NodeAssignement) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }

func (n *NodeBlock) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {

	/*blockName := "entry" // function block

	if n.Depth == 1 { // root anonymous block
		blockName = "rootblk"
	}*/

	var ret *llvm.Value

	for _, s := range n.Statements {
		ret, err := s.CodeGen(mod, builder)
		if err != nil || ret == nil {
			return ret, err
		}
		if n.Depth == 1 {
			DebugDump(ret)
		}
	}
	return ret, nil
}

func (n *NodeCall) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }

func (n *NodePrototype) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	var f_args []llvm.Type
	for _, arg := range n.Args {
		f_args = append(f_args, LLVMType(arg.Type.Name))
	}

	f_name := n.Name.Name
	f_ltype := LLVMType(n.Type.Name)
	f_type := llvm.FunctionType(f_ltype, f_args, false)
	f := llvm.AddFunction(*mod, f_name, f_type)

	if f.Name() != f_name {
		f.EraseFromParentAsFunction()
		return nil, fmt.Errorf("Redefinition of function %s", f_name)
	}

	return &f, nil
}

func (n *NodeFunction) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	f, err := n.Proto.CodeGen(mod, builder)
	if err != nil {
		return nil, err
	}

	entry := llvm.AddBasicBlock(*f, "entry")
	builder.SetInsertPointAtEnd(entry)

	_, err = n.Body.CodeGen(mod, builder)
	if err != nil {
		f.EraseFromParentAsFunction()
		return nil, err
	}

	// TODO: use LLVM verifyFunction

	return f, nil
}
func (n *NodeReturn) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {

	retVal, err := n.Value.CodeGen(mod, builder)
	if err != nil || retVal == nil {
		return nil, err
	}
	builder.CreateRet(*retVal)

	return retVal, nil
}
func (n *NodeExpression) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }
func (n *NodeVariable) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)   { return nil, nil }
