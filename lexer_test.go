// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
)

func runTest(name, src string) {
	l := NewLexer(name, src)
	for {
		it := l.NextItem()
		fmt.Println(it.Token, "\t\t", it)
		if it.Token == TOK_EOF || it.Token == TOK_ERROR {
			break
		}
	}
}

func ExampleLexArray() {

	src := `
	main is func() int {
		// int a[][]
		// a[0][0] = 1
		// a[1][0] = 2, a[1][1] = 3
		// a[2][0] = 4
		a is [[1], [2, 3], [4]]

		f is [f(), f(), 1]

		z is [0x5A, 0x6C, 0x61, 0x6E, 0x67, 0x21]

		a[2][0] is 3

		return z[0]
	}
	`

	runTest("ExampleLexArray", src)
	//output:
	//
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//ENDL 		 "\n"
	//COMMENT 		 "// int a[]"...
	//ENDL 		 "\n"
	//COMMENT 		 "// a[0][0]"...
	//ENDL 		 "\n"
	//COMMENT 		 "// a[1][0]"...
	//ENDL 		 "\n"
	//COMMENT 		 "// a[2][0]"...
	//ENDL 		 "\n"
	//IDENTIFIER 		 "a"
	//is 		 "is"
	//[ 		 "["
	//[ 		 "["
	//INT 		 "1"
	//] 		 "]"
	//, 		 ","
	//[ 		 "["
	//INT 		 "2"
	//, 		 ","
	//INT 		 "3"
	//] 		 "]"
	//, 		 ","
	//[ 		 "["
	//INT 		 "4"
	//] 		 "]"
	//] 		 "]"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "f"
	//is 		 "is"
	//[ 		 "["
	//IDENTIFIER 		 "f"
	//( 		 "("
	//) 		 ")"
	//, 		 ","
	//IDENTIFIER 		 "f"
	//( 		 "("
	//) 		 ")"
	//, 		 ","
	//INT 		 "1"
	//] 		 "]"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "z"
	//is 		 "is"
	//[ 		 "["
	//BYTE 		 "0x5A"
	//, 		 ","
	//BYTE 		 "0x6C"
	//, 		 ","
	//BYTE 		 "0x61"
	//, 		 ","
	//BYTE 		 "0x6E"
	//, 		 ","
	//BYTE 		 "0x67"
	//, 		 ","
	//BYTE 		 "0x21"
	//] 		 "]"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "a"
	//[ 		 "["
	//INT 		 "2"
	//] 		 "]"
	//[ 		 "["
	//INT 		 "0"
	//] 		 "]"
	//is 		 "is"
	//INT 		 "3"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "z"
	//[ 		 "["
	//INT 		 "0"
	//] 		 "]"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//EOF 		 EOF

}

func ExampleLexBitOps() {

	src := `
	main is func() int {
		i is 4 lshift 1

		if i nand 4 {
			i is 4 rshift 1
		}

		i is 4 xor 5

		if i nor 4 {

		}
	}
	`

	runTest("ExampleLexBitOps", src)
	//output:
	//
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "i"
	//is 		 "is"
	//INT 		 "4"
	//lshift 		 "lshift"
	//INT 		 "1"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//if 		 "if"
	//IDENTIFIER 		 "i"
	//nand 		 "nand"
	//INT 		 "4"
	//{ 		 "{"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "i"
	//is 		 "is"
	//INT 		 "4"
	//rshift 		 "rshift"
	//INT 		 "1"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "i"
	//is 		 "is"
	//INT 		 "4"
	//xor 		 "xor"
	//INT 		 "5"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//if 		 "if"
	//IDENTIFIER 		 "i"
	//nor 		 "nor"
	//INT 		 "4"
	//{ 		 "{"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//EOF 		 EOF

}

func ExampleLexOperator() {

	src := `
	main is func() int {
		if 4 eq 4 {
			// eq
		} 
		if 4 neq 5 {
			// neq
		}
		if 4 lt 5 {
			// lt
		}
		if 4 le 5 {
			// le
		}
		if 4 gt 5 {
			// gt
		}
		if 4 ge 5 {
			// ge
		}
		if 4 ge 5 or 4 le 5 and 4 eq 4 {
			// ge
		}
		if not 0 {

		}
	}
	`

	runTest("ExampleLexOperator", src)
	//output:
	//
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//ENDL 		 "\n"
	//if 		 "if"
	//INT 		 "4"
	//eq 		 "eq"
	//INT 		 "4"
	//{ 		 "{"
	//ENDL 		 "\n"
	//COMMENT 		 "// eq"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//if 		 "if"
	//INT 		 "4"
	//neq 		 "neq"
	//INT 		 "5"
	//{ 		 "{"
	//ENDL 		 "\n"
	//COMMENT 		 "// neq"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//if 		 "if"
	//INT 		 "4"
	//lt 		 "lt"
	//INT 		 "5"
	//{ 		 "{"
	//ENDL 		 "\n"
	//COMMENT 		 "// lt"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//if 		 "if"
	//INT 		 "4"
	//le 		 "le"
	//INT 		 "5"
	//{ 		 "{"
	//ENDL 		 "\n"
	//COMMENT 		 "// le"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//if 		 "if"
	//INT 		 "4"
	//gt 		 "gt"
	//INT 		 "5"
	//{ 		 "{"
	//ENDL 		 "\n"
	//COMMENT 		 "// gt"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//if 		 "if"
	//INT 		 "4"
	//ge 		 "ge"
	//INT 		 "5"
	//{ 		 "{"
	//ENDL 		 "\n"
	//COMMENT 		 "// ge"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//if 		 "if"
	//INT 		 "4"
	//ge 		 "ge"
	//INT 		 "5"
	//or 		 "or"
	//INT 		 "4"
	//le 		 "le"
	//INT 		 "5"
	//and 		 "and"
	//INT 		 "4"
	//eq 		 "eq"
	//INT 		 "4"
	//{ 		 "{"
	//ENDL 		 "\n"
	//COMMENT 		 "// ge"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//if 		 "if"
	//not 		 "not"
	//INT 		 "0"
	//{ 		 "{"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//EOF 		 EOF

}

func ExampleLexWhile() {

	src := `
	main is func() int {
		while {
			while 4 neq 5 {
				break
			}
			break
		}
	}
	`

	runTest("ExampleLexWhile", src)
	//output:
	//
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//ENDL 		 "\n"
	//while 		 "while"
	//{ 		 "{"
	//ENDL 		 "\n"
	//while 		 "while"
	//INT 		 "4"
	//neq 		 "neq"
	//INT 		 "5"
	//{ 		 "{"
	//ENDL 		 "\n"
	//break 		 "break"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//break 		 "break"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//EOF 		 EOF

}

func ExampleLexBool() {

	src := `
	foo is func(int64 i) bool {
		if i gt 1 {
			return true
		} 
		return false
	}

	main is func() int {
		if foo() {
			return 1
		}
		return 0
	}
	`

	runTest("ExampleLexBool", src)
	//output:
	//
	//ENDL 		 "\n"
	//IDENTIFIER 		 "foo"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//IDENTIFIER 		 "int64"
	//IDENTIFIER 		 "i"
	//) 		 ")"
	//IDENTIFIER 		 "bool"
	//{ 		 "{"
	//ENDL 		 "\n"
	//if 		 "if"
	//IDENTIFIER 		 "i"
	//gt 		 "gt"
	//INT 		 "1"
	//{ 		 "{"
	//ENDL 		 "\n"
	//return 		 "return"
	//BOOL 		 "true"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//return 		 "return"
	//BOOL 		 "false"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//ENDL 		 "\n"
	//if 		 "if"
	//IDENTIFIER 		 "foo"
	//( 		 "("
	//) 		 ")"
	//{ 		 "{"
	//ENDL 		 "\n"
	//return 		 "return"
	//INT 		 "1"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//return 		 "return"
	//INT 		 "0"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//EOF 		 EOF

}

func ExampleLexInteger() {

	src := `
	foo is func(int64 i) {
		return i
	}

	square is func(int32 x) int32 {
		return x*x
	}

	main is func() int {
		result is (square(2)+1)*2
		result is 1+20/square(2)
		return result
	}
	`

	runTest("ExampleLexInteger", src)

	//output:
	//
	//ENDL 		 "\n"
	//IDENTIFIER 		 "foo"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//IDENTIFIER 		 "int64"
	//IDENTIFIER 		 "i"
	//) 		 ")"
	//{ 		 "{"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "i"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "square"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//IDENTIFIER 		 "int32"
	//IDENTIFIER 		 "x"
	//) 		 ")"
	//IDENTIFIER 		 "int32"
	//{ 		 "{"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "x"
	//* 		 "*"
	//IDENTIFIER 		 "x"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "result"
	//is 		 "is"
	//( 		 "("
	//IDENTIFIER 		 "square"
	//( 		 "("
	//INT 		 "2"
	//) 		 ")"
	//+ 		 "+"
	//INT 		 "1"
	//) 		 ")"
	//* 		 "*"
	//INT 		 "2"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "result"
	//is 		 "is"
	//INT 		 "1"
	//+ 		 "+"
	//INT 		 "20"
	/// 		 "/"
	//IDENTIFIER 		 "square"
	//( 		 "("
	//INT 		 "2"
	//) 		 ")"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "result"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//EOF 		 EOF

}

func ExampleLexFloat() {

	src := `
	foo is func(float64 i) {
		return i
	}

	square is func(float32 x) float32 {
		return x*x
	}

	main is func() float {
		result is (square(2.)+1.0)*2.
		result is 1.0+20./square(2.0)
		return 23.3.
	}
	`

	runTest("ExampleLexFloat", src)

	//output:
	//
	//ENDL 		 "\n"
	//IDENTIFIER 		 "foo"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//IDENTIFIER 		 "float64"
	//IDENTIFIER 		 "i"
	//) 		 ")"
	//{ 		 "{"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "i"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "square"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//IDENTIFIER 		 "float32"
	//IDENTIFIER 		 "x"
	//) 		 ")"
	//IDENTIFIER 		 "float32"
	//{ 		 "{"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "x"
	//* 		 "*"
	//IDENTIFIER 		 "x"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "float"
	//{ 		 "{"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "result"
	//is 		 "is"
	//( 		 "("
	//IDENTIFIER 		 "square"
	//( 		 "("
	//FLOAT 		 "2."
	//) 		 ")"
	//+ 		 "+"
	//FLOAT 		 "1.0"
	//) 		 ")"
	//* 		 "*"
	//FLOAT 		 "2."
	//ENDL 		 "\n"
	//IDENTIFIER 		 "result"
	//is 		 "is"
	//FLOAT 		 "1.0"
	//+ 		 "+"
	//FLOAT 		 "20."
	/// 		 "/"
	//IDENTIFIER 		 "square"
	//( 		 "("
	//FLOAT 		 "2.0"
	//) 		 ")"
	//ENDL 		 "\n"
	//return 		 "return"
	//ERROR 		 misformated float 23.3.

}

func ExampleLexFunction() {

	src := `
	extern puts(string)
	extern fib_c(int, @int, int) int
	extern foo(@int) @

	square is func(int x) @int { // comment
		return x*x
	}

	main is func() int {
		puts("Hello, 世界") // support UTF-8 encoding

		// result = square
		result is square(2)

		return result
	}
	`

	runTest("ExampleLexFunction", src)

	//output:
	//
	//ENDL 		 "\n"
	//extern 		 "extern"
	//IDENTIFIER 		 "puts"
	//( 		 "("
	//IDENTIFIER 		 "string"
	//) 		 ")"
	//ENDL 		 "\n"
	//extern 		 "extern"
	//IDENTIFIER 		 "fib_c"
	//( 		 "("
	//IDENTIFIER 		 "int"
	//, 		 ","
	//IDENTIFIER 		 "@int"
	//, 		 ","
	//IDENTIFIER 		 "int"
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//ENDL 		 "\n"
	//extern 		 "extern"
	//IDENTIFIER 		 "foo"
	//( 		 "("
	//IDENTIFIER 		 "@int"
	//) 		 ")"
	//IDENTIFIER 		 "@"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "square"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//IDENTIFIER 		 "int"
	//IDENTIFIER 		 "x"
	//) 		 ")"
	//IDENTIFIER 		 "@int"
	//{ 		 "{"
	//COMMENT 		 "// comment"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "x"
	//* 		 "*"
	//IDENTIFIER 		 "x"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "puts"
	//( 		 "("
	//STRING 		 "Hello, 世界"...
	//) 		 ")"
	//COMMENT 		 "// support"...
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//COMMENT 		 "// result "...
	//ENDL 		 "\n"
	//IDENTIFIER 		 "result"
	//is 		 "is"
	//IDENTIFIER 		 "square"
	//( 		 "("
	//INT 		 "2"
	//) 		 ")"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "result"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//EOF 		 EOF

}

func ExampleLexComment() {

	src := `
	/*
		multiline comment
	*/
	extern puts(string)

	extern fib_c(int, int/* test */, int) int

	// pointer support
	extern foo(@int) @

	square is func(int x) int { // comment
		return x*x
	}

	main is func() int {
		puts("Hello, 世界") // support UTF-8 encoding

		// result = square
		result is square(2)
		/*if result neq 4 {
			println("FAILURE: \"not 4\"")
		} else {
			println("SUCCESS")
		}*/

		fib_c(4, 1, 0)

		return result
	}
	`

	runTest("ExampleLexComment", src)

	// output:
	//
	//ENDL 		 "\n"
	//COMMENT 		 "/*\n\t\tmulti"...
	//ENDL 		 "\n"
	//extern 		 "extern"
	//IDENTIFIER 		 "puts"
	//( 		 "("
	//IDENTIFIER 		 "string"
	//) 		 ")"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//extern 		 "extern"
	//IDENTIFIER 		 "fib_c"
	//( 		 "("
	//IDENTIFIER 		 "int"
	//, 		 ","
	//IDENTIFIER 		 "int"
	//COMMENT 		 "/* test */"
	//, 		 ","
	//IDENTIFIER 		 "int"
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//COMMENT 		 "// pointer"...
	//ENDL 		 "\n"
	//extern 		 "extern"
	//IDENTIFIER 		 "foo"
	//( 		 "("
	//IDENTIFIER 		 "@int"
	//) 		 ")"
	//IDENTIFIER 		 "@"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "square"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//IDENTIFIER 		 "int"
	//IDENTIFIER 		 "x"
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//COMMENT 		 "// comment"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "x"
	//* 		 "*"
	//IDENTIFIER 		 "x"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "main"
	//is 		 "is"
	//func 		 "func"
	//( 		 "("
	//) 		 ")"
	//IDENTIFIER 		 "int"
	//{ 		 "{"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "puts"
	//( 		 "("
	//STRING 		 "Hello, 世界"...
	//) 		 ")"
	//COMMENT 		 "// support"...
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//COMMENT 		 "// result "...
	//ENDL 		 "\n"
	//IDENTIFIER 		 "result"
	//is 		 "is"
	//IDENTIFIER 		 "square"
	//( 		 "("
	//INT 		 "2"
	//) 		 ")"
	//ENDL 		 "\n"
	//COMMENT 		 "/*if resul"...
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//IDENTIFIER 		 "fib_c"
	//( 		 "("
	//INT 		 "4"
	//, 		 ","
	//INT 		 "1"
	//, 		 ","
	//INT 		 "0"
	//) 		 ")"
	//ENDL 		 "\n"
	//ENDL 		 "\n"
	//return 		 "return"
	//IDENTIFIER 		 "result"
	//ENDL 		 "\n"
	//} 		 "}"
	//ENDL 		 "\n"
	//EOF 		 EOF

}
