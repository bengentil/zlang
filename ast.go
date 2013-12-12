// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	"strconv"
)

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
		NodeExpr
		Value int
	}

	NodeFloat struct {
		NodeExpr
		Value float32
	}

	NodeIdentifier struct {
		NodeExpr
		Name string
	}

	NodeBinOperator struct {
		NodeExpr
		Operator string
		LHS      NodeExpr
		RHS      NodeExpr
	}

	NodeAssignement struct {
		NodeExpr
		LHS *NodeIdentifier
		RHS NodeExpr
	}

	NodeBlock struct {
		NodeExpr
		Statements []NodeStmt
	}

	NodeCallExpr struct {
		NodeExpr
		Name *NodeIdentifier
		Args []NodeExpr
	}
)

// exprNode() ensures that only statement nodes can be
// assigned to a NodeExpr.
func (n *NodeInteger) exprNode()     {}
func (n *NodeFloat) exprNode()       {}
func (n *NodeIdentifier) exprNode()  {}
func (n *NodeBinOperator) exprNode() {}
func (n *NodeAssignement) exprNode() {}
func (n *NodeBlock) exprNode()       {}
func (n *NodeCallExpr) exprNode()    {}

// ********************************************
// Statements
// ********************************************
type (
	NodeStmt interface {
		Node
		stmtNode()
	}

	NodePrototype struct {
		NodeStmt
		Type *NodeIdentifier
		Name *NodeIdentifier
		Args []*NodeVariable
	}

	NodeFunction struct {
		NodeStmt
		Proto *NodePrototype
		Body  *NodeBlock
	}

	NodeReturn struct {
		NodeStmt
		Value NodeExpr
	}

	NodeExpStmt struct {
		NodeStmt
		Expression NodeExpr
	}

	NodeVariable struct {
		NodeStmt
		Type       *NodeIdentifier
		Name       *NodeIdentifier
		AssignExpr NodeExpr
	}
)

// stmtNode() ensures that only statement nodes can be
// assigned to a NodeStmt.
func (n *NodePrototype) stmtNode() {}
func (n *NodeFunction) stmtNode()  {}
func (n *NodeReturn) stmtNode()    {}
func (n *NodeExpStmt) stmtNode()   {}
func (n *NodeVariable) stmtNode()  {}

// ********************************************
// Initiate nodes
// ********************************************

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

func NIdentifier(name string) *NodeIdentifier {
	return &NodeIdentifier{Name: name}
}

func NVariable(name, typ *NodeIdentifier) *NodeVariable {
	return &NodeVariable{Name: name, Type: typ}
}

func NPrototype(name, typ *NodeIdentifier, args []*NodeVariable) *NodePrototype {
	return &NodePrototype{Name: name, Type: typ, Args: args}
}

func NFunction(proto *NodePrototype, body *NodeBlock) *NodeFunction {
	return &NodeFunction{Proto: proto, Body: body}
}

func NReturn(exp NodeExpr) *NodeReturn {
	return &NodeReturn{Value: exp}
}

func NExprStmt(exp NodeExpr) *NodeExpStmt {
	return &NodeExpStmt{Expression: exp}
}

func NAssignement(lhs *NodeIdentifier, rhs NodeExpr) *NodeAssignement {
	return &NodeAssignement{LHS: lhs, RHS: rhs}
}

func NBlock(stmts []NodeStmt) *NodeBlock {
	return &NodeBlock{Statements: stmts}
}

func NCallExpr(name *NodeIdentifier, args []NodeExpr) *NodeCallExpr {
	return &NodeCallExpr{Name: name, Args: args}
}

func NBinOp(operator string, lhs, rhs NodeExpr) *NodeBinOperator {
	return &NodeBinOperator{Operator: operator, LHS: lhs, RHS: rhs}
}

// ********************************************
// Make nodes printable
// ********************************************
func (n *NodeIdentifier) String() string {
	return fmt.Sprintf("%v", n.Name)
}

func (n *NodeVariable) String() string {
	return fmt.Sprintf("Var{name:%v, type:%v}", n.Name, n.Type)
}

func (n *NodePrototype) String() string {
	return fmt.Sprintf("Proto{name:%v, type:%v, args:%v}", n.Name, n.Type, n.Args)
}

func (n *NodeFunction) String() string {
	return fmt.Sprintf("Func{%v}", n.Proto)
}

func (n *NodeReturn) String() string {
	return fmt.Sprintf("Return{%v}", n.Value)
}

// ********************************************
// Code generation
// ********************************************
func (n *NodeInteger) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)     { return nil, nil }
func (n *NodeFloat) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)       { return nil, nil }
func (n *NodeIdentifier) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)  { return nil, nil }
func (n *NodeBinOperator) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }
func (n *NodeAssignement) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }

func (n *NodeBlock) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) {
	return nil, nil
}

func (n *NodeCallExpr) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }

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

	retVal, err := n.Body.CodeGen(mod, builder)
	if err != nil {
		f.EraseFromParentAsFunction()
		return nil, err
	}

	// TODO: use LLVM verifyFunction

	builder.CreateRet(*retVal)
	return f, nil
}
func (n *NodeExpStmt) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error)  { return nil, nil }
func (n *NodeVariable) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }
