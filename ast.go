// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	"strconv"
)

type ContextValue struct {
	FunctionName string
	VariableName string
}

var contextVariable map[ContextValue]*llvm.Value

func getContextVariable(fName, varName string) *llvm.Value {
	return contextVariable[ContextValue{fName, varName}]
}

func setContextVariable(fName, varName string, value *llvm.Value) {
	contextVariable[ContextValue{fName, varName}] = value
}

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

	NodeBool struct {
		Node
		NodeExpr
		Value bool
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

func NBool(value string) *NodeBool {
	val, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}
	return &NodeBool{Value: val}
}

func NInteger(value string) *NodeInteger {
	val, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &NodeInteger{Value: val}
}

func NFloat(value string) *NodeFloat {
	val, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return nil
	}
	return &NodeFloat{Value: float32(val)}
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
func (n *NodeBool) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeBool\",\"value\":%v}", n.Value))
}

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
func (n *NodeBool) exprNode()        {}
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
		String() string
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

	NodeIf struct {
		Node
		NodeStmt
		Condition NodeExpr
		Body      *NodeBlock
		Elif      []*NodeBlock // TODO: support for else if
		Else      *NodeBlock
	}

	NodeWhile struct {
		Node
		NodeStmt
		Condition NodeExpr
		Body      *NodeBlock
	}

	NodeBreak struct {
		Node
		NodeStmt
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

func NIf(cond NodeExpr, body *NodeBlock, elif []*NodeBlock, els *NodeBlock) *NodeIf {
	return &NodeIf{Condition: cond, Body: body, Elif: elif, Else: els}
}

func NWhile(cond NodeExpr, body *NodeBlock) *NodeWhile {
	return &NodeWhile{Condition: cond, Body: body}
}

func NBreak() *NodeBreak {
	return &NodeBreak{}
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

func (n *NodeIf) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeIf\",\"condition\":%v,\"body\":%v,\"elif\":%v,\"else\":%v}", n.Condition, n.Body, n.Elif, n.Else))
}

func (n *NodeWhile) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeWhile\",\"condition\":%v,\"body\":%v}", n.Condition, n.Body))
}

func (n *NodeBreak) String() string {
	return JNil(fmt.Sprintf("{\"__type\":\"NodeBreak\"}"))
}

// stmtNode() ensures that only statement nodes can be
// assigned to a NodeStmt.
func (n *NodePrototype) stmtNode()  {}
func (n *NodeFunction) stmtNode()   {}
func (n *NodeReturn) stmtNode()     {}
func (n *NodeExpression) stmtNode() {}
func (n *NodeVariable) stmtNode()   {}
func (n *NodeIf) stmtNode()         {}
func (n *NodeWhile) stmtNode()      {}
func (n *NodeBreak) stmtNode()      {}

// ********************************************
// Code generation
// ********************************************
func (n *NodeBool) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) {
	var i llvm.Value
	if n.Value {
		i = llvm.ConstInt(llvm.Int1Type(), 1, false)
	} else {
		i = llvm.ConstInt(llvm.Int1Type(), 0, false)
	}
	return &i, nil
}
func (n *NodeInteger) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) {
	i := llvm.ConstInt(llvm.Int32Type(), uint64(n.Value), false)
	return &i, nil
}
func (n *NodeFloat) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) {
	f := llvm.ConstFloat(llvm.FloatType(), float64(n.Value))
	return &f, nil
}
func (n *NodeString) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	/*s := llvm.ConstString(n.Value, true)
	gv := llvm.GlobalVariable()*/
	s := builder.CreateGlobalString(n.Value, ".str")
	v := llvm.ConstInBoundsGEP(s, []llvm.Value{llvm.ConstInt(llvm.Int32Type(), 0, false), llvm.ConstInt(llvm.Int32Type(), 0, false)})
	return &v, nil
}
func (n *NodeIdentifier) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	fName := mod.LastFunction().Name()
	v := builder.CreateLoad(*getContextVariable(fName, n.Name), n.Name)
	return &v, nil
}
func (n *NodeBinOperator) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	l, err := n.LHS.CodeGen(mod, builder)
	r, err := n.RHS.CodeGen(mod, builder)

	if l == nil || r == nil || err != nil {
		return nil, err
	}

	// TODO: check difference between CreateMul, CreateFMul
	// CreateUDiv, CreateSDiv
	switch n.Operator {
	case "*":
		res := builder.CreateMul(*l, *r, "multmp")
		return &res, nil
	case "/":
		res := builder.CreateSDiv(*l, *r, "divtmp")
		return &res, nil
	case "+":
		res := builder.CreateAdd(*l, *r, "addtmp")
		return &res, nil
	case "-":
		res := builder.CreateSub(*l, *r, "subtmp")
		return &res, nil
	case "eq":
		res := builder.CreateICmp(llvm.IntEQ, *l, *r, "cmptmp")
		return &res, nil
	case "neq":
		res := builder.CreateICmp(llvm.IntNE, *l, *r, "cmptmp")
		return &res, nil
	case "lt":
		res := builder.CreateICmp(llvm.IntSLT, *l, *r, "cmptmp")
		return &res, nil
	case "le":
		res := builder.CreateICmp(llvm.IntSLE, *l, *r, "cmptmp")
		return &res, nil
	case "gt":
		res := builder.CreateICmp(llvm.IntSGT, *l, *r, "cmptmp")
		return &res, nil
	case "ge":
		res := builder.CreateICmp(llvm.IntSGE, *l, *r, "cmptmp")
		return &res, nil
	case "and":
		res := builder.CreateAnd(*l, *r, "cmptmp")
		return &res, nil
	case "or":
		res := builder.CreateOr(*l, *r, "cmptmp")
		return &res, nil
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
		if err != nil {
			return ret, err
		}
		//ret.Dump()
	}
	return ret, nil
}

func (n *NodeCall) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	var args []llvm.Value

	if n.Name == nil {
		return nil, fmt.Errorf("Empty identifier")
	}

	f := mod.NamedFunction(n.Name.Name)
	if f.IsNil() {
		return nil, fmt.Errorf("Function %s not found", n.Name.Name)
	}

	for _, exp := range n.Args {
		v, err := exp.CodeGen(mod, builder)
		if err != nil {
			return nil, err
		}
		args = append(args, *v)
	}
	c := builder.CreateCall(f, args, "")
	return &c, nil
}

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

	// set args name
	for i, param := range f.Params() {
		//fmt.Println(n.Args[i])
		if n.Args[i].Name != nil {
			param.SetName(n.Args[i].Name.Name)
		}
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

	// alloc mem for params
	for _, param := range f.Params() {
		alloc := builder.CreateAlloca(param.Type(), param.Name())
		builder.CreateStore(param, alloc)
		setContextVariable(f.Name(), param.Name(), &alloc)
	}

	_, err = n.Body.CodeGen(mod, builder)
	if err != nil {
		f.EraseFromParentAsFunction()
		return nil, err
	}

	/*err = llvm.VerifyFunction(*f, llvm.PrintMessageAction)
	if err != nil {
		return f, err
	}*/

	return f, nil
}
func (n *NodeReturn) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {

	retVal, err := n.Value.CodeGen(mod, builder)
	if err != nil || retVal == nil {
		return nil, err
	}
	//DebugDump(retVal)
	builder.CreateRet(*retVal)

	return retVal, nil
}

func (n *NodeExpression) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	exp, err := n.Expression.CodeGen(mod, builder)
	return exp, err
}

func (n *NodeVariable) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	var lhs llvm.Value
	fName := mod.LastFunction().Name()
	if getContextVariable(fName, n.Name.Name) == nil { // variable not initialised
		lhs = builder.CreateAlloca(llvm.Int32Type(), n.Name.Name)
		setContextVariable(fName, n.Name.Name, &lhs)
	} else {
		lhs = *getContextVariable(fName, n.Name.Name)
	}

	rhs, err := n.AssignExpr.RHS.CodeGen(mod, builder)
	if err != nil {
		return nil, nil
	}

	v := builder.CreateStore(*rhs, lhs)
	/*lhs_load := builder.CreateLoad(lhs, n.Name.Name)
	setContextVariable(fName, n.Name.Name, &lhs_load)*/
	return &v, nil
}

func (n *NodeIf) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {
	cond, err := n.Condition.CodeGen(mod, builder)
	if err != nil || cond == nil {
		return nil, err
	}

	f := builder.GetInsertBlock().Parent()

	ifblk := llvm.AddBasicBlock(f, "ifcond")
	var elseblk llvm.BasicBlock
	if n.Else != nil {
		elseblk = llvm.AddBasicBlock(f, "else")
	}
	endif := llvm.AddBasicBlock(f, "endif")

	if n.Else != nil {
		builder.CreateCondBr(*cond, ifblk, elseblk)
	} else {
		builder.CreateCondBr(*cond, ifblk, endif)
	}

	builder.SetInsertPointAtEnd(ifblk)
	_, err = n.Body.CodeGen(mod, builder)
	if err != nil {
		return nil, err
	}
	ifblk = builder.GetInsertBlock()

	if ifblk.LastInstruction().IsATerminatorInst().IsNil() {
		builder.CreateBr(endif)
	}

	//DebugDump(body)

	if n.Else != nil {
		builder.SetInsertPointAtEnd(elseblk)
		_, err = n.Else.CodeGen(mod, builder)
		if err != nil {
			return nil, err
		}
		//DebugDump(els)

		elseblk = builder.GetInsertBlock()
		if elseblk.LastInstruction().IsATerminatorInst().IsNil() {
			builder.CreateBr(endif)
		}
	}

	builder.SetInsertPointAtEnd(endif)
	return cond, nil
}

func (n *NodeWhile) CodeGen(mod *llvm.Module, builder *llvm.Builder) (*llvm.Value, error) {

	f := builder.GetInsertBlock().Parent()
	whileloop := llvm.AddBasicBlock(f, "whilecond")
	whilebody := llvm.AddBasicBlock(f, "whilebody")
	endloop := llvm.AddBasicBlock(f, "endwhile")

	builder.CreateBr(whileloop) //go into loop

	builder.SetInsertPointAtEnd(whileloop)

	cond, err := n.Condition.CodeGen(mod, builder)
	if err != nil || cond == nil {
		return nil, err
	}

	builder.CreateCondBr(*cond, whilebody, endloop)

	builder.SetInsertPointAtEnd(whilebody)
	_, err = n.Body.CodeGen(mod, builder)
	if err != nil {
		return nil, err
	}

	builder.CreateBr(whileloop) // loop

	builder.SetInsertPointAtEnd(endloop)
	return cond, nil
}

func (n *NodeBreak) CodeGen(*llvm.Module, *llvm.Builder) (*llvm.Value, error) { return nil, nil }
