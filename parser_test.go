// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	//"os"
)

var src = `
	/*
		1st Hello world script in zlang
	*/
	extern puts(string)

	extern printi(int)
	extern fib_c(int, int, int) int

	// pointer support
	extern foo(@int) @

	square is func(int x) int {
		return x*x
	}

	fib is func(int n, int fn1, int fn) int {
		if n eq 0 {
			return fn
		} else {
			return fib_c(n-1, fn, fn + fn1)
		}

		return 0
	}

	main is func() int {
		puts("Hello, 世界") // support UTF-8 encoding

		// result = square
		result is square(2)
		if result neq 4 {
			puts("FAILURE: \"not 4\"")
		} else {
			puts("SUCCESS")
		}

		printi(square(2)*3) // 12
		printi(2+square(2)*3) // 14
		printi(square(2)/2+4*2) // 10
		printi(fib_c(square(fib_c(2, 0, 1)), 0, 1)) // 5
		printi(fib(square(fib_c(2, 0, 1)), 0, 1)) // 5

		return result
	}


	`

func ExampleParser() {

	EnableDebug()
	p := NewParser("ExampleParser", src)
	root, err := p.Parse()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("%v", root)

	// output:
	//
	//{"__type":"NodeBlock","statements":[{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"println"}, "type":{"__type":"NodeIdentifier","name":""}, "args":[{"__type":"NodeVariable","name":null,"type":{"__type":"NodeIdentifier","name":"string"},"assign_expr":null}]},{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"foo"}, "type":{"__type":"NodeIdentifier","name":"@"}, "args":[{"__type":"NodeVariable","name":null,"type":{"__type":"NodeIdentifier","name":"@int"},"assign_expr":null}]},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"square"}, "type":{"__type":"NodeIdentifier","name":"int"}, "args":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"x"},"type":{"__type":"NodeIdentifier","name":"int"},"assign_expr":null}]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeReturn","value":{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeIdentifier","name":"x"}, "rhs":{"__type":"NodeIdentifier","name":"x"}}}],"depth":2}},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"main"}, "type":{"__type":"NodeIdentifier","name":"int"}, "args":[]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeExpression","expression":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"println"}, "args":[{"__type":"NodeString","name":"\"Hello, 世界\""}]}},{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"result"},"type":null,"assign_expr":{"__type":"NodeAssignement","lhs":{"__type":"NodeIdentifier","name":"result"}, "rhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeInteger","value":2}]}}},{"__type":"NodeReturn","value":{"__type":"NodeIdentifier","name":"result"}}],"depth":2}}],"depth":1}
}

func ExampleCodeGeneration() {

	EnableDebug()
	p := NewParser("ExampleCodeGeneration", src)
	root, err := p.Parse()
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	_, err = root.CodeGen(&p.Module, &p.Builder)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	DebugDumpMod(&p.Module)

	err = llvm.VerifyModule(p.Module, llvm.ReturnStatusAction)
	if err != nil {
		fmt.Println(err)
		return
	}

	llvm.LinkInJIT()
	llvm.InitializeNativeTarget()

	engine, err := llvm.NewJITCompiler(p.Module, 2)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer engine.Dispose()

	pass := llvm.NewPassManager()
	defer pass.Dispose()

	pass.Add(engine.TargetData())
	pass.AddConstantPropagationPass()
	pass.AddInstructionCombiningPass()
	pass.AddPromoteMemoryToRegisterPass()
	pass.AddGVNPass()
	pass.AddCFGSimplificationPass()
	pass.Run(p.Module)

	var args []llvm.GenericValue
	result := engine.RunFunction(p.Module.NamedFunction("main"), args)
	fmt.Printf("result=%v\n", result.Int(false))
	//DebugDumpMod(&p.Module)
	//llvm.WriteBitcodeToFile(p.Module, os.Stdout)

	// output:
	//
}
