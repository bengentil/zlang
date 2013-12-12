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

		// result = square
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
--------

``` bash
$ cat hello_world.zl | zlang
Hello, 世界
SUCCESS
```

Status
------

Zlang is at **early alpha stage** and is for educational purposes only (**don't implement real-world application with Zlang**)

- [x] Lexer
- [ ] Parser (In progress)
- [ ] JIT Interpreter (In progress)
- [ ] Compiler