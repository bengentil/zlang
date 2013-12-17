// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	//"os"
)

func runTestParser(name, src string) (*Parser, *NodeBlock) {
	EnableDebug()
	p := NewParser(name, src)
	root, err := p.Parse()
	if err != nil {
		fmt.Println("error:", err)
		return nil, nil
	}

	fmt.Printf("%v", root)

	return p, root
}

func runTestCodeGen(name, src string) *llvm.GenericValue {
	p, root := runTestParser(name, src)
	_, err := root.CodeGen(&p.Module, &p.Builder)
	if err != nil {
		fmt.Println("error:", err)
		return nil
	}
	err = llvm.VerifyModule(p.Module, llvm.ReturnStatusAction)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	DebugDumpMod(&p.Module)

	llvm.LinkInJIT()
	llvm.InitializeNativeTarget()

	engine, err := llvm.NewJITCompiler(p.Module, 2)
	if err != nil {
		fmt.Println(err)
		return nil
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
	return &result
}

func ExampleParseOperator() {
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
	runTestParser("ExampleParseOperator", src)
	//output:
	//
	//{"__type":"NodeBlock","statements":[{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"main"}, "type":{"__type":"NodeIdentifier","name":"int"}, "args":[]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"eq", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":4}},"body":{"__type":"NodeBlock","statements":[],"depth":3},"elif":[],"else":null},{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"neq", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":5}},"body":{"__type":"NodeBlock","statements":[],"depth":3},"elif":[],"else":null},{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"lt", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":5}},"body":{"__type":"NodeBlock","statements":[],"depth":3},"elif":[],"else":null},{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"le", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":5}},"body":{"__type":"NodeBlock","statements":[],"depth":3},"elif":[],"else":null},{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"gt", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":5}},"body":{"__type":"NodeBlock","statements":[],"depth":3},"elif":[],"else":null},{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"ge", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":5}},"body":{"__type":"NodeBlock","statements":[],"depth":3},"elif":[],"else":null},{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"and", "lhs":{"__type":"NodeBinOperator","op":"or", "lhs":{"__type":"NodeBinOperator","op":"ge", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":5}}, "rhs":{"__type":"NodeBinOperator","op":"le", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":5}}}, "rhs":{"__type":"NodeBinOperator","op":"eq", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":4}}},"body":{"__type":"NodeBlock","statements":[],"depth":3},"elif":[],"else":null},{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"not", "lhs":null, "rhs":{"__type":"NodeInteger","value":0}},"body":{"__type":"NodeBlock","statements":[],"depth":3},"elif":[],"else":null}],"depth":2}}],"depth":1}
}

func ExampleParseBool() {

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

	runTestParser("ExampleParseBool", src)
	//output:
	//
	//{"__type":"NodeBlock","statements":[{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"foo"}, "type":{"__type":"NodeIdentifier","name":"bool"}, "args":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"i"},"type":{"__type":"NodeIdentifier","name":"int64"},"assign_expr":null}]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeIf","condition":{"__type":"NodeBinOperator","op":"gt", "lhs":{"__type":"NodeIdentifier","name":"i"}, "rhs":{"__type":"NodeInteger","value":1}},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeReturn","value":{"__type":"NodeBool","value":true}}],"depth":3},"elif":[],"else":null},{"__type":"NodeReturn","value":{"__type":"NodeBool","value":false}}],"depth":2}},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"main"}, "type":{"__type":"NodeIdentifier","name":"int"}, "args":[]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeIf","condition":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"foo"}, "args":[]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeReturn","value":{"__type":"NodeInteger","value":1}}],"depth":3},"elif":[],"else":null},{"__type":"NodeReturn","value":{"__type":"NodeInteger","value":0}}],"depth":2}}],"depth":1}

}

func ExampleParseInteger() {

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

	runTestParser("ExampleParseInteger", src)
	//output:
	//
	//{"__type":"NodeBlock","statements":[{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"foo"}, "type":{"__type":"NodeIdentifier","name":""}, "args":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"i"},"type":{"__type":"NodeIdentifier","name":"int64"},"assign_expr":null}]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeReturn","value":{"__type":"NodeIdentifier","name":"i"}}],"depth":2}},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"square"}, "type":{"__type":"NodeIdentifier","name":"int32"}, "args":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"x"},"type":{"__type":"NodeIdentifier","name":"int32"},"assign_expr":null}]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeReturn","value":{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeIdentifier","name":"x"}, "rhs":{"__type":"NodeIdentifier","name":"x"}}}],"depth":2}},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"main"}, "type":{"__type":"NodeIdentifier","name":"int"}, "args":[]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"result"},"type":null,"assign_expr":{"__type":"NodeAssignement","lhs":{"__type":"NodeIdentifier","name":"result"}, "rhs":{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeBinOperator","op":"+", "lhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeInteger","value":2}]}, "rhs":{"__type":"NodeInteger","value":1}}, "rhs":{"__type":"NodeInteger","value":2}}}},{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"result"},"type":null,"assign_expr":{"__type":"NodeAssignement","lhs":{"__type":"NodeIdentifier","name":"result"}, "rhs":{"__type":"NodeBinOperator","op":"+", "lhs":{"__type":"NodeInteger","value":1}, "rhs":{"__type":"NodeBinOperator","op":"/", "lhs":{"__type":"NodeInteger","value":20}, "rhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeInteger","value":2}]}}}}},{"__type":"NodeReturn","value":{"__type":"NodeIdentifier","name":"result"}}],"depth":2}}],"depth":1}

}

func ExampleParseFloat() {

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
		return 23.3
	}
	`

	runTestParser("ExampleParseFloat", src)
	//output:
	//
	//{"__type":"NodeBlock","statements":[{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"foo"}, "type":{"__type":"NodeIdentifier","name":""}, "args":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"i"},"type":{"__type":"NodeIdentifier","name":"float64"},"assign_expr":null}]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeReturn","value":{"__type":"NodeIdentifier","name":"i"}}],"depth":2}},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"square"}, "type":{"__type":"NodeIdentifier","name":"float32"}, "args":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"x"},"type":{"__type":"NodeIdentifier","name":"float32"},"assign_expr":null}]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeReturn","value":{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeIdentifier","name":"x"}, "rhs":{"__type":"NodeIdentifier","name":"x"}}}],"depth":2}},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"main"}, "type":{"__type":"NodeIdentifier","name":"float"}, "args":[]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"result"},"type":null,"assign_expr":{"__type":"NodeAssignement","lhs":{"__type":"NodeIdentifier","name":"result"}, "rhs":{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeBinOperator","op":"+", "lhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeFloat","value":2}]}, "rhs":{"__type":"NodeFloat","value":1}}, "rhs":{"__type":"NodeFloat","value":2}}}},{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"result"},"type":null,"assign_expr":{"__type":"NodeAssignement","lhs":{"__type":"NodeIdentifier","name":"result"}, "rhs":{"__type":"NodeBinOperator","op":"+", "lhs":{"__type":"NodeFloat","value":1}, "rhs":{"__type":"NodeBinOperator","op":"/", "lhs":{"__type":"NodeFloat","value":20}, "rhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeFloat","value":2}]}}}}},{"__type":"NodeReturn","value":{"__type":"NodeFloat","value":23.3}}],"depth":2}}],"depth":1}

}

func ExampleParseFunction() {

	src := `
	extern puts(string)
	extern fib_c(int, @int, int) int
	extern foo(@int) @

	square is func(int x) @int { // comment
		return x*x
	}

	main is func() int {
		puts("Hello, 世界") // support UTF-8 encoding

		printi((2+square(2))*3) // 18
		printi(square(2)*3) // 12
		printi(2+square(2)*3) // 14
		printi(square(2)/2+4*2) // 10
		printi(fib_c(square(fib_c(2, 0, 1)), 0, 1)) // 5
		printi(fib(square(fib_c(2, 0, 1)), 0, 1)) // 5

		return result
	}
	`

	runTestParser("ExampleParseFunction", src)
	//output:
	//
	//{"__type":"NodeBlock","statements":[{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"puts"}, "type":{"__type":"NodeIdentifier","name":""}, "args":[{"__type":"NodeVariable","name":null,"type":{"__type":"NodeIdentifier","name":"string"},"assign_expr":null}]},{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"fib_c"}, "type":{"__type":"NodeIdentifier","name":"int"}, "args":[{"__type":"NodeVariable","name":null,"type":{"__type":"NodeIdentifier","name":"int"},"assign_expr":null},{"__type":"NodeVariable","name":null,"type":{"__type":"NodeIdentifier","name":"@int"},"assign_expr":null},{"__type":"NodeVariable","name":null,"type":{"__type":"NodeIdentifier","name":"int"},"assign_expr":null}]},{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"foo"}, "type":{"__type":"NodeIdentifier","name":"@"}, "args":[{"__type":"NodeVariable","name":null,"type":{"__type":"NodeIdentifier","name":"@int"},"assign_expr":null}]},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"square"}, "type":{"__type":"NodeIdentifier","name":"@int"}, "args":[{"__type":"NodeVariable","name":{"__type":"NodeIdentifier","name":"x"},"type":{"__type":"NodeIdentifier","name":"int"},"assign_expr":null}]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeReturn","value":{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeIdentifier","name":"x"}, "rhs":{"__type":"NodeIdentifier","name":"x"}}}],"depth":2}},{"__type":"NodeFunction","proto":{"__type":"NodePrototype","name":{"__type":"NodeIdentifier","name":"main"}, "type":{"__type":"NodeIdentifier","name":"int"}, "args":[]},"body":{"__type":"NodeBlock","statements":[{"__type":"NodeExpression","expression":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"puts"}, "args":[{"__type":"NodeString","name":"Hello, 世界"}]}},{"__type":"NodeExpression","expression":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"printi"}, "args":[{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeBinOperator","op":"+", "lhs":{"__type":"NodeInteger","value":2}, "rhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeInteger","value":2}]}}, "rhs":{"__type":"NodeInteger","value":3}}]}},{"__type":"NodeExpression","expression":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"printi"}, "args":[{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeInteger","value":2}]}, "rhs":{"__type":"NodeInteger","value":3}}]}},{"__type":"NodeExpression","expression":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"printi"}, "args":[{"__type":"NodeBinOperator","op":"+", "lhs":{"__type":"NodeInteger","value":2}, "rhs":{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeInteger","value":2}]}, "rhs":{"__type":"NodeInteger","value":3}}}]}},{"__type":"NodeExpression","expression":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"printi"}, "args":[{"__type":"NodeBinOperator","op":"+", "lhs":{"__type":"NodeBinOperator","op":"/", "lhs":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeInteger","value":2}]}, "rhs":{"__type":"NodeInteger","value":2}}, "rhs":{"__type":"NodeBinOperator","op":"*", "lhs":{"__type":"NodeInteger","value":4}, "rhs":{"__type":"NodeInteger","value":2}}}]}},{"__type":"NodeExpression","expression":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"printi"}, "args":[{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"fib_c"}, "args":[{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"fib_c"}, "args":[{"__type":"NodeInteger","value":2},{"__type":"NodeInteger","value":0},{"__type":"NodeInteger","value":1}]}]},{"__type":"NodeInteger","value":0},{"__type":"NodeInteger","value":1}]}]}},{"__type":"NodeExpression","expression":{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"printi"}, "args":[{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"fib"}, "args":[{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"square"}, "args":[{"__type":"NodeCall","name":{"__type":"NodeIdentifier","name":"fib_c"}, "args":[{"__type":"NodeInteger","value":2},{"__type":"NodeInteger","value":0},{"__type":"NodeInteger","value":1}]}]},{"__type":"NodeInteger","value":0},{"__type":"NodeInteger","value":1}]}]}},{"__type":"NodeReturn","value":{"__type":"NodeIdentifier","name":"result"}}],"depth":2}}],"depth":1}

}

/*
var src = `
	extern puts(string)

	extern printi(int)
	extern fib_c(int, int, int) int

	// pointer support
	//extern foo(@int) @

	square is func(int x) int {
		return x*x
	}

	fib is func(int n, int fn1, int fn) int {
		if n eq 0 {
			return fn
		} else {
			return fib(n-1, fn, fn + fn1)
		}

		return 0
	}

	main is func() int {
		puts("Hello, 世界") // support UTF-8 encoding

		// result = square
		result is square(2)
		if result neq 4 {
			puts("FAILURE: \"not 4\"")
			printi(result)
		} else {
			puts("SUCCESS")
		}

		result is square(4)/square(2)

		printi((2+square(2))*3) // 18
		printi(square(2)*3) // 12
		printi(2+square(2)*3) // 14
		printi(square(2)/2+4*2) // 10
		printi(fib_c(square(fib_c(2, 0, 1)), 0, 1)) // 5
		printi(fib(square(fib_c(2, 0, 1)), 0, 1)) // 5
		return result
	}


	`

/*
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
}*/
