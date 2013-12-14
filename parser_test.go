// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
)

func ExampleParser() {

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
		/*if result neq 4 {
			println("FAILURE: \"not 4\"")
		} else {
			println("SUCCESS")
		}*/

		return result
	}


	`

	EnableDebug()
	p := NewParser("ExampleParser", src)
	root, err := p.Parse()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("%v\n\n", root)

	code, err := root.CodeGen(&p.Module, &p.Builder)
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	DebugDump(code)

	// output:
	//
}
