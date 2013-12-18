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

func LLVMType(name string) llvm.Type {
	// no type (in function declaration for instance)
	if name == "" {
		return llvm.VoidType()
	}

	// pointer
	if name[0] == POINTER_CHAR {
		return llvm.PointerType(LLVMType(name[1:]), 0)
	}

	// basic type
	switch name {
	case TYPE_BOOL:
		return llvm.Int1Type()
	case TYPE_BYTE:
		return llvm.Int8Type()
	case TYPE_INT, TYPE_INT32:
		return llvm.Int32Type()
	case TYPE_INT64:
		return llvm.Int64Type()
	case TYPE_FLOAT, TYPE_FLOAT32:
		return llvm.FloatType()
	case TYPE_FLOAT64:
		return llvm.DoubleType()
	case TYPE_STRING:
		// i8* to have char*
		return llvm.PointerType(llvm.Int8Type(), 0)
	}

	// if none has been resolved
	return llvm.VoidType()
}

// TODO: return zlang type instead of llvm type
func ZlangType(t llvm.Type) string {
	return ""
}
