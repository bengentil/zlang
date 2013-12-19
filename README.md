Zlang
=====

Implementation of a basic Language, designed for readability


Example:
--------

``` go
	/*
		1st Hello world script in zlang
	*/

	// extern keyword to allow binding of C functions
	// C equivalent: extern void println(char*)
	extern println(string)

	// pointer support
	// C equivalent: extern void* foo(int*)
	extern foo(@int) @

	square is func(int x) int {
		return x*x
	}

	main is func() int {
		println("Hello, 世界") // support UTF-8 encoding

		// result = square(2)
		result is square(2)
		if result neq 4 {
			println("FAILURE: \"not 4\"")
		} else {
			println("SUCCESS")
		}

		return 0
	}

```

Running:
-------------------------

``` bash

# with compiler
$ zlc -o hello hello_world.zl 
$ ./hello
Hello, 世界
SUCCESS


# with interpreter
$ cat hello_world.zl | zlang
Hello, 世界
SUCCESS


```

Debugging Bytecode
------------------

``` bash
# Dump bytecode
$ zlc --emit-llvm hello_world.zl 
	; ModuleID = 'hello_world.zl'

	@.str = private unnamed_addr constant [20 x i8] c"Hello World, \E4\B8\96\E7\95\8C\00", align 1
	@.str1 = private unnamed_addr constant [17 x i8] c"FAILURE: \22not 4\22\00", align 1
	@.str2 = private unnamed_addr constant [8 x i8] c"SUCCESS\00", align 1

	declare void @println(i8*)
	declare void* @foo(i32*)

	define i32 @square(i32 %x) nounwind uwtable {
	  %1 = alloca i32, align 4
	  store i32 %x, i32* %1, align 4
	  %2 = load i32* %1, align 4
	  %3 = load i32* %1, align 4
	  %4 = mul nsw i32 %2, %3
	  ret i32 %4
	}

	define i32 @main() nounwind uwtable {
	  %1 = alloca i32, align 4
	  %result = alloca i32, align 4
	  store i32 0, i32* %1
	  call void @println(i8* getelementptr inbounds ([20 x i8]* @.str, i32 0, i32 0))
	  %2 = call i32 @square(i32 2)
	  store i32 %2, i32* %result, align 4
	  %3 = load i32* %result, align 4
	  %4 = icmp eq i32 %3, 4
	  br i1 %4, label %5, label %6

	; <label>:5                                       ; preds = %0
	  call void @println(i8* getelementptr inbounds ([17 x i8]* @.str1, i32 0, i32 0))
	  br label %7

	; <label>:6                                       ; preds = %0
	  call void @println(i8* getelementptr inbounds ([8 x i8]* @.str2, i32 0, i32 0))
	  br label %7

	; <label>:7                                       ; preds = %6, %5
	  ret i32 0
	}

```

Status
------

Zlang is at **early alpha stage** and is for educational purposes only (**don't implement real-world application with Zlang**)

- [x] Lexer
- [x] Parser
- [ ] JIT Interpreter (In progress)
- [x] Compiler
- [ ] Compiler with optimizer pass

Language features:
- [x] Functions
- [x] Extern functions (C bindings)
- [x] Variable assignation
- [x] Mutable Variables
- [-] Pointers (To be tested)
- [x] Conditions (if/else)
- [x] Loops (while)
- [x] Bitwise operations (and, or, xor, nand, nor, lshift, rshift)

Types:
- [x] Boolean (bool)
- [x] String (string)
- [ ] Byte (byte)
- [x] Integer (int32, int64)
- [x] Float (float32, float64)
- [ ] Arrays
- [ ] Structs

