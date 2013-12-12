// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
)

func ExampleLexer() {

	src := `
	/*
		1st Hello world script in zlang
	*/
	extern println(string)

	// pointer support
	extern foo(@int) @

	square is func(int x) int {
		return x*x
	}

	main is func() int {
		println("Hello, 世界") // support UTF-8 encoding

		// result = square
		result is square(2)
		if result neq 4 {
			println("FAILURE: \"not 4\"")
		} else {
			println("SUCCESS")
		}

		return 0
	}


	`
	l := NewLexer("ExampleLexer", src)
	for {
		it := l.NextItem()
		fmt.Println(it.Token, "\t\t", it)
		if it.Token == TOK_EOF || it.Token == TOK_ERROR {
			break
		}
	}

	// output:
	//
	// ENDL 		 "\n"
	// COMMENT 		 "/*\n\t\t1st H"...
	// ENDL 		 "\n"
	// extern 		 "extern"
	// IDENTIFIER 		 "println"
	// ( 		 "("
	// IDENTIFIER 		 "string"
	// ) 		 ")"
	// ENDL 		 "\n"
	// ENDL 		 "\n"
	// COMMENT 		 "// pointer"...
	// ENDL 		 "\n"
	// extern 		 "extern"
	// IDENTIFIER 		 "foo"
	// ( 		 "("
	// IDENTIFIER 		 "@int"
	// ) 		 ")"
	// IDENTIFIER 		 "@"
	// ENDL 		 "\n"
	// ENDL 		 "\n"
	// IDENTIFIER 		 "square"
	// is 		 "is"
	// func 		 "func"
	// ( 		 "("
	// IDENTIFIER 		 "int"
	// IDENTIFIER 		 "x"
	// ) 		 ")"
	// IDENTIFIER 		 "int"
	// { 		 "{"
	// ENDL 		 "\n"
	// return 		 "return"
	// IDENTIFIER 		 "x"
	// * 		 "*"
	// IDENTIFIER 		 "x"
	// ENDL 		 "\n"
	// } 		 "}"
	// ENDL 		 "\n"
	// ENDL 		 "\n"
	// IDENTIFIER 		 "main"
	// is 		 "is"
	// func 		 "func"
	// ( 		 "("
	// ) 		 ")"
	// IDENTIFIER 		 "int"
	// { 		 "{"
	// ENDL 		 "\n"
	// IDENTIFIER 		 "println"
	// ( 		 "("
	// STRING 		 "\"Hello, 世界"...
	// ) 		 ")"
	// COMMENT 		 "// support"...
	// ENDL 		 "\n"
	// ENDL 		 "\n"
	// COMMENT 		 "// result "...
	// ENDL 		 "\n"
	// IDENTIFIER 		 "result"
	// is 		 "is"
	// IDENTIFIER 		 "square"
	// ( 		 "("
	// INT 		 "2"
	// ) 		 ")"
	// ENDL 		 "\n"
	// if 		 "if"
	// IDENTIFIER 		 "result"
	// neq 		 "neq"
	// INT 		 "4"
	// { 		 "{"
	// ENDL 		 "\n"
	// IDENTIFIER 		 "println"
	// ( 		 "("
	// STRING 		 "\"FAILURE: "...
	// ) 		 ")"
	// ENDL 		 "\n"
	// } 		 "}"
	// else 		 "else"
	// { 		 "{"
	// ENDL 		 "\n"
	// IDENTIFIER 		 "println"
	// ( 		 "("
	// STRING 		 "\"SUCCESS\""
	// ) 		 ")"
	// ENDL 		 "\n"
	// } 		 "}"
	// ENDL 		 "\n"
	// ENDL 		 "\n"
	// return 		 "return"
	// INT 		 "0"
	// ENDL 		 "\n"
	// } 		 "}"
	// ENDL 		 "\n"
	// ENDL 		 "\n"
	// ENDL 		 "\n"
	// EOF 		 EOF
}
