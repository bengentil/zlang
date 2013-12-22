// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
)

type Token int

const (
	TOK_ERROR   Token = iota
	TOK_COMMENT       // /* comment */ or // comment
	TOK_BLANK         // space, tab
	TOK_ENDL          // \r, \n
	TOK_IDENTIFIER
	TOK_EOF

	TOK_BOOL  // true / false
	TOK_FLOAT // 1.0
	TOK_INT   // 1
	//TOK_CHAR   // 'a'
	TOK_BYTE   // 0x20
	TOK_STRING // "azerty"

	operator_begin
	/*
		TOK_ASSIGN // =
		TOK_NOT    // !
		TOK_EQUAL  // ==
		TOK_NEQ    // !=
		TOK_LT     // <
		TOK_LE     // <=
		TOK_GT     // >
		TOK_GE     // >=

		TOK_AND  // &&
		TOK_OR   // ||
		TOK_NAND //
		TOK_NOR  //
		TOK_XOR  //
	*/

	TOK_ASSIGN_S // is (=)
	TOK_NOT_S    // not (!)
	TOK_EQUAL_S  // eq (==)
	TOK_NEQ_S    // neq (!=)
	TOK_LT_S     // lt (<)
	TOK_LE_S     // le (<=)
	TOK_GT_S     // gt (>)
	TOK_GE_S     // ge (>=)

	TOK_AND_S  // and
	TOK_OR_S   // or
	TOK_NAND_S // nand
	TOK_NOR_S  // nor
	TOK_XOR_S  // xor
	TOK_LSH_S  // lshift
	TOK_RSH_S  // rshift

	TOK_PLUS  // +
	TOK_MINUS // -
	TOK_MUL   // *
	TOK_DIV   // /

	TOK_INC // ++
	TOK_DEC // --

	//TOK_POINTER // @
	operator_end

	delimiter_begin
	TOK_LPAREN   // (
	TOK_RPAREN   // )
	TOK_LBLOCK   // {
	TOK_RBLOCK   // }
	TOK_LBRACKET // [
	TOK_RBRACKET // ]
	TOK_DOT      // .
	TOK_COMMA    // ,
	delimiter_end

	keyword_begin
	TOK_FUNC   // function declaration
	TOK_EXTERN // extern (C bindings)
	TOK_RETURN // return
	TOK_IF     // if
	TOK_ELSE   // else
	TOK_SWITCH // switch
	TOK_CASE   // case
	TOK_BREAK  // break
	TOK_FOR    // for
	TOK_WHILE  // while
	keyword_end
)

var tokens = [...]string{
	TOK_ERROR:      "ERROR",
	TOK_COMMENT:    "COMMENT",
	TOK_BLANK:      "BLANK",
	TOK_ENDL:       "ENDL",
	TOK_IDENTIFIER: "IDENTIFIER",
	TOK_EOF:        "EOF",

	TOK_BOOL:  "BOOL",
	TOK_FLOAT: "FLOAT",
	TOK_INT:   "INT",
	TOK_BYTE:  "BYTE",
	//TOK_CHAR:   "CHAR",
	TOK_STRING: "STRING",

	// operator_begin
	/*
		TOK_ASSIGN: "=",
		TOK_NOT:    "!",
		TOK_EQUAL:  "==",
		TOK_NEQ:    "!=",
		TOK_LT:     "<",
		TOK_LE:     "<=",
		TOK_GT:     ">",
		TOK_GE:     ">=",
	*/

	TOK_ASSIGN_S: "is",
	TOK_NOT_S:    "not",
	TOK_EQUAL_S:  "eq",
	TOK_NEQ_S:    "neq",
	TOK_LT_S:     "lt",
	TOK_LE_S:     "le",
	TOK_GT_S:     "gt",
	TOK_GE_S:     "ge",

	TOK_AND_S:  "and",
	TOK_OR_S:   "or",
	TOK_NAND_S: "nand",
	TOK_NOR_S:  "nor",
	TOK_XOR_S:  "xor",
	TOK_LSH_S:  "lshift",
	TOK_RSH_S:  "rshift",

	TOK_PLUS:  "+",
	TOK_MINUS: "-",
	TOK_MUL:   "*",
	TOK_DIV:   "/",

	TOK_INC: "++",
	TOK_DEC: "--",

	//TOK_POINTER: "@",
	// operator_end

	// delimiter_begin
	TOK_LPAREN:   "(",
	TOK_RPAREN:   ")",
	TOK_LBLOCK:   "{",
	TOK_RBLOCK:   "}",
	TOK_LBRACKET: "[",
	TOK_RBRACKET: "]",
	TOK_DOT:      ".",
	TOK_COMMA:    ",",
	// delimiter_end

	// keyword_begin
	TOK_FUNC:   "func",
	TOK_EXTERN: "extern",
	TOK_RETURN: "return",
	TOK_IF:     "if",
	TOK_ELSE:   "else",
	TOK_SWITCH: "switch",
	TOK_CASE:   "case",
	TOK_BREAK:  "break",
	TOK_FOR:    "for",
	TOK_WHILE:  "while",
	// keyword_end
}

// Get a printable token
func (tok Token) String() string {
	s := tokens[tok]
	if s == "" {
		return fmt.Sprintf("TOK_%d", int(tok))
	}
	return s
}

// used to resolve tokens
var keywords map[string]Token
var operators map[string]Token
var delimiters map[string]Token

func resolveOperator(op string) Token {
	if tok_op, is_operator := operators[op]; is_operator {
		return tok_op
	}

	return TOK_ERROR // not an operator
}

func resolveDelimiter(del string) Token {
	if tok_del, is_delimiter := delimiters[del]; is_delimiter {
		return tok_del
	}

	return TOK_ERROR // not an delimiter
}

// identify if a string is an operator(is, not, and...), a keyword(if, for...) or an identifier
func resolveIdentifier(identifier string) Token {
	if tok_op, is_operator := operators[identifier]; is_operator {
		return tok_op
	}
	if tok_key, is_keyword := keywords[identifier]; is_keyword {
		return tok_key
	}

	if identifier == "true" || identifier == "false" {
		return TOK_BOOL
	}

	return TOK_IDENTIFIER
}

func (t Token) Precedence() int {
	switch t {
	case TOK_AND_S, TOK_OR_S, TOK_NOT_S, TOK_NAND_S, TOK_NOR_S, TOK_XOR_S:
		return 2
	case TOK_EQUAL_S, TOK_NEQ_S, TOK_LT_S, TOK_LE_S, TOK_GT_S, TOK_GE_S:
		return 3
	case TOK_LSH_S, TOK_RSH_S:
		return 4
	case TOK_PLUS, TOK_MINUS:
		return 5
	case TOK_MUL, TOK_DIV:
		return 6
	}
	return 0
}

func (t Token) IsOperator() bool {
	if t > operator_begin && t < operator_end {
		return true
	}
	return false
}
