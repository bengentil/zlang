// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import "github.com/axw/gollvm/llvm"

func init() {
	// store token associations in maps
	keywords = make(map[string]Token)
	for i := keyword_begin + 1; i < keyword_end; i++ {
		keywords[tokens[i]] = i
	}
	operators = make(map[string]Token)
	for i := operator_begin + 1; i < operator_end; i++ {
		operators[tokens[i]] = i
	}
	delimiters = make(map[string]Token)
	for i := delimiter_begin + 1; i < delimiter_end; i++ {
		delimiters[tokens[i]] = i
	}

	// construct context map for code generation
	contextVariable = make(map[ContextValue]*llvm.Value)
}
