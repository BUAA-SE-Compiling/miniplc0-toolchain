package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gitlab.lazymio.cn/mio/miniplc0/vm"
)

var debugHelp = `Simple miniplc0 debugger.
You can use the abbreviation of a command.
[H]elp -- Show this message.
[N]ext -- Run a single instruction.
[L]ist n -- List n instructions.
[S]tack n -- Show n stack elemets.
[I]formation -- Show current information.
[Q]uit -- Quit the debugger.
`

// [R]estart -- Restart the debugging. Not impleneted

type cmd int32

const (
	cNext cmd = iota
	cList
	cStack
	cInformation
	cRestart
	cHelp
	cQuit
)

type debuggerCommand struct {
	C cmd
	X int32
}

func newCommandFromString(line string) *debuggerCommand {
	ln := strings.TrimSpace(line)
	tokens := strings.Split(ln, " ")
	if len(tokens) == 0 {
		return nil
	}
	if len(tokens) == 1 {
		switch strings.ToLower(tokens[0]) {
		case "h":
			fallthrough
		case "help":
			return &debuggerCommand{C: cHelp}
		case "q":
			fallthrough
		case "quit":
			return &debuggerCommand{C: cQuit}
		case "r":
			fallthrough
		case "restart":
			//return &debuggerCommand{C: cRestart}
			return nil
		case "i":
			fallthrough
		case "info":
			fallthrough
		case "infomation":
			return &debuggerCommand{C: cInformation}
		case "l":
			fallthrough
		case "list":
			return &debuggerCommand{C: cList, X: 10}
		case "s":
			fallthrough
		case "stack":
			return &debuggerCommand{C: cStack, X: 20}
		case "n":
			fallthrough
		case "next":
			return &debuggerCommand{C: cNext}
		}
		return nil
	}
	if len(tokens) == 2 {
		x, err := strconv.ParseInt(tokens[1], 10, 32)
		if err != nil {
			return nil
		}
		switch strings.ToLower(tokens[0]) {
		case "l":
			fallthrough
		case "list":
			return &debuggerCommand{C: cList, X: int32(x)}
		case "s":
			fallthrough
		case "stack":
			return &debuggerCommand{C: cStack, X: int32(x)}
		}
		return nil
	}
	return nil
}

// Debug debugs the epf file.
func Debug(file *os.File) {
	epf, err := vm.NewEPFv1FromFile(file)
	if err != nil {
		panic(err)
	}
	v := vm.NewVMDefault(len(epf.GetInstructions()), epf.GetEntry())
	v.Load(epf.GetInstructions())
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print(debugHelp)
	for true {
		fmt.Print(">")
		scanner.Scan()
		cmd := newCommandFromString(scanner.Text())
		quit := false
		if cmd == nil {
			fmt.Println("Wrong format.\nType 'help' to see more.")
			continue
		}
		switch cmd.C {
		case cHelp:
			fmt.Print(debugHelp)
		case cQuit:
			quit = true
			break
		case cRestart:
			quit = true
			break
		case cNext:
			if err = v.RunSingle(); err != nil {
				if err == vm.ErrIllegalInstruction {
					fmt.Println("The program has stopped.")
				} else {
					// Should have better user experience.
					fmt.Println(err)
					fmt.Println(*v)
					os.Exit(0)
				}
			}
			fmt.Printf("Next instruction: %v\n", *(v.NextInstruction()))
		case cInformation:
			fmt.Printf("IP=%v SP=%v\nInstructions[IP]:%v\n", v.IP, v.SP, *(v.NextInstruction()))
			stackvalue := v.GetStackTop()
			if stackvalue != nil {
				fmt.Printf("Stack[SP-1]=%v\n", *stackvalue)
			} else {
				fmt.Println("Stack[SP-1]=[Invalid]")
			}
		case cList:
			fmt.Println(v.InstructionGraph(cmd.X))
		case cStack:
			fmt.Println(v.StackGraph(cmd.X))
		}
		if quit {
			break
		}
	}
}
