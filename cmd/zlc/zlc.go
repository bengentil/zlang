// Copyright 2013 Benjamin Gentil. All rights reserved.
// license can be found in the LICENSE file (MIT License)
package main

import (
	"fmt"
	"github.com/axw/gollvm/llvm"
	"github.com/bengentil/zlang"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime/pprof"
	"strings"
)

const (
	ResetCode  = "\033[0m" // reset color
	RedCode    = "\033[31m"
	GreenCode  = "\033[32m"
	YellowCode = "\033[33m"
)

type FileType int

const (
	begin_source_files FileType = iota
	FILE_TYPE_C
	FILE_TYPE_ZL
	FILE_TYPE_LL
	FILE_TYPE_BC
	FILE_TYPE_S
	end_source_files

	FILE_TYPE_O
)

var FileExtensions = [...]string{
	FILE_TYPE_C:  ".c",
	FILE_TYPE_ZL: ".zl",
	FILE_TYPE_LL: ".ll",
	FILE_TYPE_BC: ".bc",
	FILE_TYPE_S:  ".s",

	FILE_TYPE_O: ".o",
}

type InputFile struct {
	Name string
	Type FileType
}

func (i *InputFile) String() string {
	return i.Name
}

type Options struct {
	Verbose    bool   `short:"v" long:"verbose" description:"verbose output" default:"false"`
	Debug      bool   `long:"debug" description:"debug output" default:"false"`
	Cpuprofile string `long:"cpuprofile" description:"file to output pprof data" optional:"yes"`
	EmitLLVM   bool   `long:"emit-llvm" description:"print llvm IR and exit" optional:"yes"`
	EmitAST    bool   `long:"emit-ast" description:"print AST JSON representation and exit" optional:"yes"`
	LlcBin     string `long:"llc-bin" description:"llc command to execute" default:"llc"`
	ClangBin   string `long:"clang-bin" description:"clang command to execute" default:"clang"`
	NoLink     bool   `short:"c" description:"don't link (produce object file)" default:"false"`
	Force      bool   `short:"f" description:"force compilation" default:"false"`
	Output     string `short:"o" long:"output" description:"output file" default:"a.out"`
}

var options Options
var inputFiles []*InputFile
var parser = flags.NewParser(&options, flags.Default)

func verbose(format string, v ...interface{}) {
	if options.Verbose {
		fmt.Fprintf(os.Stderr, YellowCode+format+ResetCode+"\n", v...)
	}
}

func stderr(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

func (f *InputFile) SetType(t FileType) {
	f.Name = f.ShortName() + FileExtensions[t]
	f.Type = t
}

func (f *InputFile) ShortName() string {
	end := len(f.Name) - len(FileExtensions[f.Type])
	if end > 0 {
		return f.Name[:end]
	} else {
		return f.Name
	}
}

func llc(in *InputFile) error {
	verbose("[llc (%s) %v]", options.LlcBin, in)
	out, err := exec.Command(options.LlcBin, in.Name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("llc error: %v\nOut:%s\n", err, out)
	}
	in.SetType(FILE_TYPE_S)
	return nil
}

func clang_c(in *InputFile) error {
	verbose("[clang_c (%s) %v]", options.ClangBin, in)
	out, err := exec.Command(options.ClangBin, "-c", in.Name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("clang error: %v\nOut:%s\n", err, out)
	}
	in.SetType(FILE_TYPE_O)
	return nil
}

func clang_link(in []*InputFile) error {
	verbose("[clang_link (%s) %v]", options.ClangBin, in)
	args := []string{"-o", options.Output}

	for _, f := range in {
		args = append(args, f.Name)
	}

	out, err := exec.Command(options.ClangBin, args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("clang error: %v\nOut:%s\n", err, out)
	}
	return nil
}

func zlc(in *InputFile) error {
	/*
		f, err := os.Open(in.Name)
		if err != nil {
			return err
		}
		defer f.Close()
	*/

	if options.Debug {
		zlang.EnableDebug()
	}

	verbose("[zlc %s]", in.Name)

	source, err := ioutil.ReadFile(in.Name)
	if err != nil {
		return err
	}

	verbose("%s: parsing", in.Name)
	parser := zlang.NewParser(in.Name, string(source))
	root, err := parser.Parse()
	if err != nil {
		return err
	}

	if options.EmitAST {
		fmt.Printf("--> %s:\n%v\n\n", in.Name, root)
		in.SetType(FILE_TYPE_BC)
		return nil
	}

	verbose("%s: bytecode generation", in.Name)
	_, err = root.CodeGen(&parser.Module, &parser.Builder)
	if err != nil {
		return err
	}

	in.SetType(FILE_TYPE_BC)

	if options.EmitLLVM {
		fmt.Printf("--> %s:\n", in.Name)
		parser.Module.Dump()
		return nil
	}

	if !options.Force {
		verbose("--> bytecode verification")
		err = llvm.VerifyModule(parser.Module, llvm.ReturnStatusAction)
		if err != nil {
			return err
		}
	}

	out, err := os.Create(in.Name)
	defer out.Close()
	if err != nil {
		return err
	}

	verbose("%s: writing bytecode to file", in.Name)
	return llvm.WriteBitcodeToFile(parser.Module, out)
}

func (f *InputFile) IsSourceFile() bool {
	return f.Type > begin_source_files && f.Type < end_source_files
}

func (f *InputFile) IsObjectFile() bool {
	return f.Type == FILE_TYPE_O
}

func (f *InputFile) Compile() error {
	switch f.Type {
	case FILE_TYPE_ZL:
		return zlc(f)
	case FILE_TYPE_LL, FILE_TYPE_BC:
		return llc(f)
	case FILE_TYPE_C, FILE_TYPE_S:
		return clang_c(f)
	}

	return nil
}

func CompileFiles() error {
Loop:
	for {
		i := 0
		for _, f := range inputFiles {
			// Compile all source files
			// except if EmitAST|EmitLLVM are enabled, in this case only ZL
			if f.IsSourceFile() && (!(options.EmitAST || options.EmitLLVM) || f.Type == FILE_TYPE_ZL) {
				verbose("Compiling %s ...", f.Name)
				err := f.Compile()
				if err != nil {
					return err
				}
				i++
			}
		}

		// no more file to compile
		if i == 0 {
			break Loop
		}
	}

	return nil
}

func LinkFiles() error {
	for _, f := range inputFiles {
		if !f.IsObjectFile() {
			return fmt.Errorf("File %s isn't an object file", f.Name)
		}
	}
	return clang_link(inputFiles)
}

func handleFile(filename string) (*InputFile, error) {
	for i, ext := range FileExtensions {
		if ext != "" && strings.HasSuffix(filename, ext) {
			return &InputFile{filename, FileType(i)}, nil
		}
	}
	return nil, fmt.Errorf("%s file type not supported", filename)
}

func main() {
	args, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if options.Cpuprofile != "" {
		f, err := os.Create(options.Cpuprofile)
		if err != nil {
			stderr("%v", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if len(args) < 1 {
		stderr("error: no input file")
		os.Exit(1)
	}

	for _, file := range args {
		f, err := handleFile(file)
		if err != nil {
			stderr("error: %v", err)
			os.Exit(1)
		}
		inputFiles = append(inputFiles, f)
	}

	verbose("# Compilation Phase")
	err = CompileFiles()
	if err != nil {
		stderr("compile error: %v", err)
		os.Exit(1)
	}

	if !options.NoLink && !options.EmitAST && !options.EmitLLVM {
		verbose("# Link Phase")
		err = LinkFiles()
		if err != nil {
			stderr("compile error: %v", err)
			os.Exit(1)
		}
	}
}
