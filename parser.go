// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	//"strconv"
)

type Parser struct {
	Input       string
	Module      llvm.Module
	Builder     llvm.Builder
	l           *Lexer
	currentItem LexItem
	errorCount  int
	depth       int
}

func NewParser(name, input string) *Parser {
	return &Parser{Input: input, Module: llvm.NewModule(name), Builder: llvm.NewBuilder(), l: NewLexer(name, input), currentItem: LexItem{Token: TOK_BLANK, Val: ""}, errorCount: 0, depth: 0}
}

func (p *Parser) NextItem() {
	if p.l != nil && p.currentItem.Token != TOK_EOF {
		p.currentItem = p.l.NextItem()
		if p.currentItem.Token == TOK_ERROR {
			p.RaiseError("Formating error, %v", p.currentItem.Val)
			p.NextItem()
		}
		// Skip all comments
		if p.currentItem.Token == TOK_COMMENT {
			p.NextItem()
		}
	}
}

func (p *Parser) DepthIndent() string {
	s := ""
	for i := 0; i < p.depth; i++ {
		s += "  "
	}
	return fmt.Sprintf("[%s:%d] ", p.l.name, p.l.lineNumber())
}

func (p *Parser) Pos() string {
	return fmt.Sprintf("[%s:%d] ", p.l.name, p.l.lineNumber())
}

func (p *Parser) RaiseError(f string, v ...interface{}) {
	p.errorCount++
	Errorln(p.Pos()+f, v...)
}

func (p *Parser) parseBoolean() *NodeBool {
	return NBool(p.currentItem.Val)
}

func (p *Parser) parseInteger() *NodeInteger {
	return NInteger(p.currentItem.Val)
}

func (p *Parser) parseFloat() *NodeFloat {
	return NFloat(p.currentItem.Val)
}

/*
func (p *Parser) parseIdentifier() NodeExpr {

	identifierName := p.currentItem.Val
	p.NextItem() // skip identifier

	if p.currentItem.Token == TOK_IDENTIFIER {
		return NIdentifier(p.currentItem.Val)
	}

	if p.currentItem.Token == TOK_ENDL {
		return NIdentifier(identifierName)
	}

	// function call
	if p.currentItem.Token == TOK_LPAREN {
		p.NextItem()
		return p.parseFunctionCall(identifierName)
	}

	if p.currentItem.Token == TOK_MUL {
		return NIdentifier(identifierName)//p.parseBinOP(0, )
	}

	p.RaiseError("Expected identifier, got %v", p.currentItem.Token)
	return nil
}*/

func (p *Parser) parseParen() NodeExpr {
	p.NextItem() // skip '('
	exp := p.parseExpression()
	if exp == nil {
		return nil
	}
	if p.currentItem.Token != TOK_RPAREN {
		p.RaiseError("Expected ')', got %v", p.currentItem.Token)
		return nil
	}
	//p.NextItem() // skip ')'
	return exp
}

func (p *Parser) parseBinOP(prec int, LHS NodeExpr) NodeExpr {
	for {
		//fmt.Println(p.currentItem.Val)
		tok_prec := p.currentItem.Token.Precedence()

		//fmt.Printf("Token : %v", LHS)
		if tok_prec < prec {
			//fmt.Printf("LHS: %v", LHS)
			return LHS
		}

		op := p.currentItem.Val
		p.NextItem() // skip operator
		RHS := p.parsePrimary()
		if RHS == nil {
			return nil
		}

		next_prec := p.currentItem.Token.Precedence()
		if tok_prec < next_prec {
			RHS = p.parseBinOP(tok_prec+1, RHS)
			if RHS == nil {
				return nil
			}
		}

		LHS = NBinOp(op, LHS, RHS)

	}
	return nil
}

func (p *Parser) parsePrimary() NodeExpr {
	var result NodeExpr

	switch p.currentItem.Token {
	case TOK_IDENTIFIER:
		identifierName := p.currentItem.Val
		p.NextItem()
		if p.currentItem.Token == TOK_LPAREN {
			p.NextItem() // skip '('
			return p.parseFunctionCall(identifierName)
		}
		return NIdentifier(identifierName)
	case TOK_STRING:
		result = NString(p.currentItem.Val)
	case TOK_BOOL:
		result = p.parseBoolean()
	case TOK_INT:
		result = p.parseInteger()
	case TOK_FLOAT:
		result = p.parseFloat()
	case TOK_LPAREN:
		result = p.parseParen()
	case TOK_RBLOCK, TOK_RPAREN:
		return nil
	case TOK_NOT_S:
		return p.parseBinOP(1, nil)
	}

	if result != nil {
		p.NextItem()
		return result
	}

	p.RaiseError("Expected expression, got %v", p.currentItem.Token)
	return nil
}

func (p *Parser) parseExpression() NodeExpr {
	LHS := p.parsePrimary()
	if LHS == nil {
		return nil
	}

	//Debug("LHS: %v\n", LHS)

Loop:
	for {
		//println(p.currentItem.Token.String())
		/*switch p.currentItem.Token {
		case TOK_MUL, TOK_DIV, TOK_PLUS, TOK_MINUS, TOK_NEQ_S, TOK_EQUAL_S: // is operator
			//op := p.currentItem.Val
			//p.NextItem()
			//RHS := p.parseExpression()
			LHS = p.parseBinOP(1, LHS)
			//p.NextItem()
		default:
			break Loop
		}*/
		if p.currentItem.Token.IsOperator() {
			LHS = p.parseBinOP(1, LHS)
		} else {
			break Loop
		}
	}

	//Debug("LHS: %v\nTOken: %v\n", LHS, p.currentItem.Token)
	//return p.parseBinOP(0, LHS)*/
	return LHS
}

func (p *Parser) parseVariable(varName string) *NodeVariable {
	LHS := NIdentifier(varName)
	RHS := p.parseExpression()
	return NVariable(LHS, nil, NAssignement(LHS, RHS))
}

func (p *Parser) parseWhile() *NodeWhile {
	p.NextItem() // skip if keyword

	var condition NodeExpr

	// no condition means while true
	if p.currentItem.Token == TOK_LBLOCK {
		condition = NBool("true")
	} else {
		condition = p.parseExpression()
		if p.currentItem.Token != TOK_LBLOCK {
			p.RaiseError("Expected while body '{', got %v", p.currentItem.Token)
			return nil
		}
	}

	p.NextItem() // skip left block
	body := p.parseBlock()

	if p.currentItem.Token != TOK_RBLOCK {
		p.RaiseError("Expected while closing '}', got %v", p.currentItem.Token)
		return nil
	}
	p.NextItem() // skip right block

	return NWhile(condition, body)

}

func (p *Parser) parseIf() *NodeIf {
	p.NextItem() // skip if keyword

	condition := p.parseExpression()

	if p.currentItem.Token != TOK_LBLOCK {
		p.RaiseError("Expected if body '{', got %v", p.currentItem.Token)
		return nil
	}

	// if
	p.NextItem() // skip left block
	body := p.parseBlock()

	if p.currentItem.Token != TOK_RBLOCK {
		p.RaiseError("Expected if closing '}', got %v", p.currentItem.Token)
		return nil
	}

	p.NextItem() // skip right block

	var elif []*NodeBlock
	var els *NodeBlock

	if p.currentItem.Token == TOK_ELSE {
		p.NextItem() // skip else

		// else if
		if p.currentItem.Token == TOK_IF {
		Loop:
			for {
				p.NextItem() // skip if

				if p.currentItem.Token != TOK_LBLOCK {
					p.RaiseError("Expected else if body '{', got %v", p.currentItem.Token)
					return nil
				}
				p.NextItem() // skip left block

				e := p.parseBlock()
				elif = append(elif, e)

				if p.currentItem.Token != TOK_RBLOCK {
					p.RaiseError("Expected else if '}', got %v", p.currentItem.Token)
					return nil
				}
				p.NextItem() // skip right block

				if p.currentItem.Token != TOK_ELSE {
					return NIf(condition, body, elif, els)
				}

				if p.currentItem.Token != TOK_IF {
					break Loop
				}
			}
		}

		// else
		if p.currentItem.Token != TOK_LBLOCK {
			p.RaiseError("Expected else body '{', got %v", p.currentItem.Token)
			return nil
		}

		p.NextItem() // skip left block
		els = p.parseBlock()

		if p.currentItem.Token != TOK_RBLOCK {
			p.RaiseError("Expected else '}', got %v", p.currentItem.Token)
			return nil
		}
		p.NextItem() // skip right block
	}

	return NIf(condition, body, elif, els)
}

func (p *Parser) parseFunctionCall(funcName string) *NodeCall {
	var args []NodeExpr
Loop:
	for {
		arg := p.parseExpression()
		if arg != nil {
			//Debug("->arg:%v\n", arg)
			args = append(args, arg)
		}

		if p.currentItem.Token == TOK_RPAREN {
			break Loop
		}

		if p.currentItem.Token != TOK_COMMA {
			p.RaiseError("Expected ')' or ',' in arg list, got %v", p.currentItem.Token)
			return nil
		}

		p.NextItem() // skip comma
	}

	p.NextItem() // skip ')'
	return NCall(NIdentifier(funcName), args)
}

func (p *Parser) parseFunction(funcName string) *NodeFunction {
	p.NextItem() // skip func keyword

	proto := p.parsePrototype(funcName)
	if proto == nil {
		return nil
	}

	if p.currentItem.Token != TOK_LBLOCK {
		p.RaiseError("Expected '{' after function definition, got %v", p.currentItem.Token)
		return nil
	}
	p.NextItem() // skip left block
	body := p.parseBlock()
	p.NextItem() // skip right block
	return NFunction(proto, body)
}

func (p *Parser) parseReturn() *NodeReturn {
	p.NextItem() // skip return keyword
	return NReturn(p.parseExpression())
}

func (p *Parser) parsePrototype(funcName string) *NodePrototype {
	fName := NIdentifier(funcName)

	if p.currentItem.Token != TOK_LPAREN {
		p.RaiseError("Expected '(' in prototype, got %v", p.currentItem.Token)
		return nil
	}

	p.NextItem() // skip '('

	var args []*NodeVariable
Loop:
	for {
		if p.currentItem.Token == TOK_IDENTIFIER {
			var varType, varName *NodeIdentifier
			varType = NIdentifier(p.currentItem.Val)
			p.NextItem()
			if p.currentItem.Token == TOK_IDENTIFIER {
				varName = NIdentifier(p.currentItem.Val)
				p.NextItem()
			}
			variable := NVariable(varName, varType, nil)
			args = append(args, variable)

			if p.currentItem.Token == TOK_COMMA {
				p.NextItem() // skip ,
			} else {
				break Loop
			}
		} else {
			break Loop
		}

	}

	if p.currentItem.Token != TOK_RPAREN {
		p.RaiseError("Expected ')' in prototype, got %v", p.currentItem.Token)
		return nil
	}

	p.NextItem() // skip ')'

	typeName := NIdentifier("")
	// handle return type
	if p.currentItem.Token == TOK_IDENTIFIER {
		typeName = NIdentifier(p.currentItem.Val)
		p.NextItem()
	}

	return NPrototype(fName, typeName, args)
}

func (p *Parser) parseExtern() *NodePrototype {
	p.NextItem() // skip extern keyword
	if p.currentItem.Token != TOK_IDENTIFIER {
		p.RaiseError("Expected function name in prototype, got %v", p.currentItem.Token)
		return nil
	}
	fName := p.currentItem.Val
	p.NextItem() // skip func name
	return p.parsePrototype(fName)
}

func (p *Parser) parseStatement() NodeStmt {
	identifierName := p.currentItem.Val
	p.NextItem() // skip identifier

	switch p.currentItem.Token {
	case TOK_LPAREN: // function call
		p.NextItem() // skip '('
		return NExpression(p.parseFunctionCall(identifierName))
	case TOK_ASSIGN, TOK_ASSIGN_S: // variable or function definition
		p.NextItem()
		if p.currentItem.Token == TOK_FUNC {
			return p.parseFunction(identifierName)
		}

		return p.parseVariable(identifierName)
	case TOK_COMMENT:
		p.NextItem() // skip comment
		if p.currentItem.Token == TOK_IDENTIFIER {
			return p.parseStatement()
		}
	}

	p.RaiseError("Unexpected token '%v', expect function call or variable assignation", p.currentItem.Token)
	return nil
}

/*
func (p *Parser) parseTopExpression() *NodeFunction {
	body := p.parseExpression()
	proto := NPrototype(NIdentifier(""), NIdentifier(""), []*NodeVariable{})
	return NFunction(proto, body)
}

func (p *Parser) handleExtern() {
	if ast := p.parseExtern(); ast != nil {
		Debug("Parsed an extern func:\n\t%q\n", ast)
		f, err := ast.CodeGen(&p.Module, &p.Builder)
		if err != nil {
			p.RaiseError("%v", err)
			return
		}
		DebugDump(f)
	} else {
		p.NextItem()
	}
}

func (p *Parser) handleExpression() {
	if ast := p.parseExpression(); ast != nil {
		Debug("Parsed an expression:\n\t%q\n", ast)
		f, err := ast.CodeGen(&p.Module, &p.Builder)
		if err != nil {
			p.RaiseError("%v", err)
			return
		}
		DebugDump(f)
	} else {
		p.NextItem()
	}
}*/

func (p *Parser) parseBlock() *NodeBlock {
	var statements []NodeStmt
	var stmt NodeStmt

	p.depth++

	if p.depth < 1 {
		p.RaiseError("Unexpected depth %d, should be > 1", p.depth)
		return nil
	}

Loop:
	for {
		stmt = nil
		switch p.currentItem.Token {
		case TOK_RBLOCK:
			// expected for depth != 1
			if p.depth == 1 {
				p.RaiseError("Unexpected '}' on top level block")
			}
			break Loop
		case TOK_EOF:
			// expected for depth == 1
			if p.depth != 1 {
				p.RaiseError("Unexpected end of file, missing '}'")
			}
			break Loop
		case TOK_EXTERN:
			if p.depth == 1 {
				stmt = p.parseExtern()
			} else {
				p.RaiseError("Extern function can't be declared inside function")
				break Loop
			}
		case TOK_IDENTIFIER:
			stmt = p.parseStatement()
		case TOK_RETURN:
			stmt = p.parseReturn()
		case TOK_ENDL, TOK_COMMENT:
			p.NextItem()
		case TOK_IF:
			stmt = p.parseIf()
		case TOK_WHILE:
			stmt = p.parseWhile()
		case TOK_BREAK:
			stmt = NBreak()
			p.NextItem()
		default:
			p.RaiseError("Unexpected token '%v'", p.currentItem.Token)
			break Loop
		}

		if stmt != nil {
			//Debug("STMT[%d]: %v\n", p.depth, stmt)
			statements = append(statements, stmt)
		}
	}

	p.depth--
	return NBlock(statements, p.depth+1)
}

func (p *Parser) Parse() (*NodeBlock, error) {
	p.NextItem()           // get the first item
	root := p.parseBlock() // parse top block

	if p.errorCount == 1 {
		return nil, fmt.Errorf("1 error")
	}
	if p.errorCount > 1 {
		return nil, fmt.Errorf("%d errors", p.errorCount)
	}
	return root, nil
}
