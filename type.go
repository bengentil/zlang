// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package zlang

import (
	"github.com/axw/gollvm/llvm"
)

const (
	POINTER_CHAR = '@'
	TYPE_BOOL    = "bool"
	TYPE_BYTE    = "byte"
	TYPE_INT     = "int"
	TYPE_INT32   = "int32"
	TYPE_INT64   = "int64"
	TYPE_FLOAT   = "float"
	TYPE_FLOAT32 = "float32"
	TYPE_FLOAT64 = "float64"
	TYPE_STRING  = "string"
	TYPE_STRUCT  = "struct"
)

var types = map[string]llvm.Type{
	TYPE_BOOL:    llvm.Int1Type(),
	TYPE_BYTE:    llvm.Int8Type(),
	TYPE_INT:     llvm.Int32Type(),
	TYPE_INT32:   llvm.Int32Type(),
	TYPE_INT64:   llvm.Int64Type(),
	TYPE_FLOAT:   llvm.FloatType(),
	TYPE_FLOAT32: llvm.FloatType(),
	TYPE_FLOAT64: llvm.DoubleType(),
	TYPE_STRING:  llvm.PointerType(llvm.Int8Type(), 0),
}

func RegisterType(name string, t llvm.Type) {
	types[name] = t
}

func LLVMType(name string) llvm.Type {
	// no type (in function declaration for instance)
	if name == "" {
		return llvm.VoidType()
	}

	// pointer
	if name[0] == POINTER_CHAR {
		return llvm.PointerType(LLVMType(name[1:]), 0)
	}

	if t, ok := types[name]; ok {
		return t
	}

	// if none has been resolved
	return llvm.VoidType()
}

// TODO: return zlang type instead of llvm type
func ZlangType(t llvm.Type) string {
	return ""
}
