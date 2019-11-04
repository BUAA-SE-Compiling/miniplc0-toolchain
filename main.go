package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/BUAA-SE-Compiling/vm"

	flag "github.com/spf13/pflag"
)

// Run runs the epf file.
func Run(file *os.File) {
	epf, err := vm.NewEPFv1FromFile(file)
	if err != nil {
		panic(err)
	}
	vm := vm.NewVMDefault(len(epf.GetInstructions()), epf.GetEntry())
	vm.Load(epf.GetInstructions())
	if err := vm.Run(); err != nil {
		fmt.Println(vm)
		panic(err)
	}
}

// Interprete interprets the input file directly.
func Interprete(file *os.File) {
	scanner := bufio.NewScanner(file)
	instructions := []vm.Instruction{}
	linenum := 0
	for scanner.Scan() {
		linenum++
		single := vm.ParseInstruction(scanner.Text())
		if single == nil {
			fmt.Fprintf(os.Stderr, "Line %v: Bad instruction", linenum)
		}
		instructions = append(instructions, *single)
	}
	vm := vm.NewVMDefault(len(instructions), 0)
	vm.Load(instructions)
	if err := vm.Run(); err != nil {
		fmt.Println(vm)
		panic(err)
	}
}

// Decompile decompiles the epf file.
func Decompile(file *os.File) {
	epf, err := vm.NewEPFv1FromFile(file)
	if err != nil {
		panic(err)
	}
	for _, i := range epf.GetInstructions() {
		fmt.Println(i)
	}
	return
}

// Assemble assembles the text file to an epf file.
func Assemble(in *os.File, out *os.File) {
	scanner := bufio.NewScanner(in)
	instructions := []vm.Instruction{}
	linenum := 0
	for scanner.Scan() {
		linenum++
		single := vm.ParseInstruction(scanner.Text())
		if single == nil {
			fmt.Fprintf(os.Stderr, "Line %v: Bad instruction", linenum)
		}
		instructions = append(instructions, *single)
	}
	epf := vm.NewEPFv1FromInstructions(instructions, 0)
	epf.WriteFile(out)
	return
}

// We don't know whehter a flag is unset or set with an empty string if its default value is an empty string.
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func main() {
	var input string
	var output string
	var decompile bool
	var run bool
	var debug bool
	var help bool
	var assemble bool
	var interprete bool
	flag.CommandLine.Init("Default", flag.ContinueOnError)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "A vm implementation for mini plc0.\n")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.StringVarP(&input, "input", "i", "-", "The input file. The default is os.Stdin.")
	flag.StringVarP(&output, "output", "o", "", "The output file.")
	flag.BoolVarP(&decompile, "decompile", "D", false, "Decompile without running.")
	flag.BoolVarP(&run, "run", "R", false, "Run the file.")
	flag.BoolVarP(&interprete, "interprete", "I", false, "Interprete the file.")
	flag.BoolVarP(&debug, "debug", "d", false, "Debug the file.")
	flag.BoolVarP(&help, "help", "h", false, "Show this message.")
	flag.BoolVarP(&assemble, "assemble", "A", false, "Assemble a text file to an EPFv1 file.")
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		fmt.Println(err)
		fmt.Fprintf(os.Stderr, "Run with --help for details.\n")
		os.Exit(2)
	}
	if help || (decompile && !isFlagPassed("output")) {
		flag.Usage()
		os.Exit(2)
	}
	if !debug && !run && !decompile && !assemble && !interprete {
		fmt.Fprintf(os.Stderr, "You must choose to decomple, run, assemble, interprete or debug.\n")
		fmt.Fprintf(os.Stderr, "Run with --help for details.\n")
		os.Exit(2)
	}
	var file *os.File
	var err error
	if input != "-" {
		file, err = os.Open(input)
		if err != nil {
			panic(err)
		}
		defer file.Close()
	} else {
		file = os.Stdin
	}
	if debug {
		Debug(file)
	} else if decompile {
		Decompile(file)
	} else if run {
		Run(file)
	} else if assemble {
		out, err := os.OpenFile(output, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0660)
		if err != nil {
			panic(err)
		}
		defer out.Close()
		Assemble(file, out)
	} else if interprete {
		Interprete(file)
	}
	os.Exit(0)
}
