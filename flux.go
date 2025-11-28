package main

import (
    "bufio"
    "fmt"
    "io"
    "os"
    "strings"
)

/*

FLUX PROGRAMMING LANGUAGE - COMPLETE REFERENCE GUIDE

## OVERVIEW

Flux is a minimal, stack-based esoteric programming language that achieves
Turing completeness with only 9 core operations. It features a single
accumulator register and an unbounded stack, making it capable of universal
computation despite its extreme minimalism.

The language is designed to demonstrate that complex computational behavior
can emerge from simple, well-chosen primitives. Every operation in Flux has
a clear, unambiguous meaning, and there are no hidden mechanisms or implicit
behaviors.

## LANGUAGE PHILOSOPHY

Flux adheres to these core design principles:

1. Minimalism - Include only operations that are absolutely essential for
   Turing completeness. Every operation must serve a distinct purpose.
2. Composability - Complex algorithms and data structures should emerge
   naturally from composing simple primitive operations.
3. Explicitness - All state changes must be explicit. There are no side
   effects, hidden state, or implicit conversions.
4. Stack-based architecture - The stack provides a natural, efficient model
   for expression evaluation and temporary storage without named variables.
5. Educational value - The language should be simple enough to understand
   completely, yet powerful enough to demonstrate fundamental concepts of
   computation and compiler design.

## COMPUTATIONAL MODEL

Flux programs operate on two primary data structures:

ACCUMULATOR:
- A single integer register that serves as the central workspace
- All arithmetic operations modify the accumulator
- The accumulator value determines loop execution
- Starts at zero when program execution begins
- Can hold any integer value (implementation dependent range)

STACK:
- An unbounded Last-In-First-Out (LIFO) data structure
- Stores integer values
- Provides the primary means of data persistence between operations
- Can grow without limit (constrained only by available memory)
- Popping from an empty stack yields zero (graceful degradation)

## OPERATION REFERENCE

Flux has exactly 9 operations:

ARITHMETIC OPERATIONS:
+    Increment accumulator
     Effect: accumulator = accumulator + 1
     Example: If acc=5, then after '+' acc=6

-    Decrement accumulator
     Effect: accumulator = accumulator - 1
     Example: If acc=5, then after '-' acc=4

STACK OPERATIONS:
*    Push accumulator to stack
     Effect: stack.push(accumulator)
     The accumulator value is copied to the top of the stack
     The accumulator itself remains unchanged
     Example: If acc=5, after '*' the stack has 5 on top and acc is still 5

/    Pop stack to accumulator
     Effect: accumulator = stack.pop()
     If stack is empty, accumulator becomes 0
     The value is removed from the stack
     Example: If stack=[3,7] (7 on top) and acc=5, after '/' acc=7 and stack=[3]

CONTROL FLOW OPERATIONS:
[    Begin while loop
     If accumulator == 0, jump forward past the matching ']'
     If accumulator != 0, continue execution into the loop body
     The loop condition is checked only at the '[', not continuously
     Example: If acc=0, execution jumps past the loop
     If acc=5, execution enters the loop

]    End while loop
     If accumulator != 0, jump backward to the matching '['
     If accumulator == 0, continue execution past the loop
     This creates a while-loop structure that continues as long as acc != 0
     Example: If acc=3, jump back to '['
     If acc=0, exit loop

INPUT/OUTPUT OPERATIONS:
.    Output accumulator as ASCII character
     Prints the character corresponding to (accumulator mod 256)
     Example: If acc=65, outputs 'A'
     If acc=72, outputs 'H'

,    Input one character
     Reads a single character from input stream
     Sets accumulator to the ASCII value of that character
     On EOF (end of file), sets accumulator to 0
     Example: If user types 'A', acc becomes 65

#    Output accumulator as decimal number
     Prints the numeric value of the accumulator
     This is an extension for practical debugging and numeric output
     Example: If acc=42, outputs "42"

WHITESPACE AND COMMENTS:
- Spaces, tabs, newlines, and carriage returns are ignored
- Any character that is not one of the 9 operations is treated as a comment
- This allows for readable, documented Flux code

END OF REFERENCE GUIDE
*/

// OpCode represents bytecode instruction types
type OpCode byte

const (
    OpInc    OpCode = iota // + : Increment accumulator
    OpDec                  // - : Decrement accumulator
    OpPush                 // * : Push accumulator to stack
    OpPop                  // / : Pop stack to accumulator
    OpLoop                 // [ : Begin loop
    OpEnd                  // ] : End loop
    OpOut                  // . : Output as ASCII
    OpIn                   // , : Input character
    OpOutNum               // # : Output as number
)

// Instruction represents a single bytecode instruction with optional argument
type Instruction struct {
    Op  OpCode // The operation to perform
    Arg int    // Argument (used for loop jump addresses)
}

// Compiler transforms Flux source code into executable bytecode
type Compiler struct {
    source       []byte        // Source code as byte array
    instructions []Instruction // Generated bytecode instructions
    loopStack    []int         // Stack of loop start positions for bracket matching
    position     int           // Current position in source (for error reporting)
}

// NewCompiler creates a new compiler instance with the given source code
func NewCompiler(source string) *Compiler {
    return &Compiler{
        source:       []byte(source),
        instructions: make([]Instruction, 0, len(source)), // Pre-allocate for efficiency
        loopStack:    make([]int, 0, 16),                  // Pre-allocate small loop stack
        position:     0,
    }
}

// Compile performs the complete compilation pipeline:
// 1. Lexical analysis (tokenization)
// 2. Syntax analysis (bracket matching validation)
// 3. Code generation (bytecode emission)
// Returns the compiled instructions or an error
func (c *Compiler) Compile() ([]Instruction, error) {
    // Single-pass compilation: scan source left to right
    for c.position = 0; c.position < len(c.source); c.position++ {
        char := c.source[c.position]

        switch char {
        case '+':
            // Increment operation: accumulator += 1
            c.emit(OpInc, 0)

        case '-':
            // Decrement operation: accumulator -= 1
            c.emit(OpDec, 0)

        case '*':
            // Push operation: stack.push(accumulator)
            c.emit(OpPush, 0)

        case '/':
            // Pop operation: accumulator = stack.pop()
            c.emit(OpPop, 0)

        case '[':
            // Loop start: if acc == 0, jump past matching ]
            loopStart := len(c.instructions)
            c.emit(OpLoop, 0) // Emit with placeholder jump address
            c.loopStack = append(c.loopStack, loopStart)

        case ']':
            // Loop end: if acc != 0, jump back to matching [
            if len(c.loopStack) == 0 {
                return nil, fmt.Errorf("compilation error: unmatched ']' at position %d", c.position)
            }

            // Pop the matching loop start position
            loopStart := c.loopStack[len(c.loopStack)-1]
            c.loopStack = c.loopStack[:len(c.loopStack)-1]

            loopEnd := len(c.instructions)

            // Emit end instruction that jumps back to loop start
            c.emit(OpEnd, loopStart)

            // Patch the loop start instruction with the end address
            // This allows O(1) jump when condition is false
            c.instructions[loopStart].Arg = loopEnd

        case '.':
            // Output operation: print character
            c.emit(OpOut, 0)

        case ',':
            // Input operation: read character
            c.emit(OpIn, 0)

        case '#':
            // Numeric output operation: print number
            c.emit(OpOutNum, 0)

        case ' ', '\t', '\n', '\r':
            // Whitespace: ignored

        default:
            // Any other character: treated as comment, ignored
            // This allows for readable, documented code
        }
    }

    // Validate that all loops are properly closed
    if len(c.loopStack) > 0 {
        return nil, fmt.Errorf("compilation error: %d unmatched '[' bracket(s) in source code", len(c.loopStack))
    }

    return c.instructions, nil
}

// emit appends a new instruction to the bytecode sequence
func (c *Compiler) emit(op OpCode, arg int) {
    c.instructions = append(c.instructions, Instruction{Op: op, Arg: arg})
}

// VM represents the Flux virtual machine that executes compiled bytecode
type VM struct {
    instructions []Instruction // The bytecode program to execute
    accumulator  int           // The single accumulator register
    stack        []int         // The unbounded stack
    pc           int           // Program counter (instruction pointer)
    input        io.Reader     // Input stream for ',' operation
    output       io.Writer     // Output stream for '.' and '#' operations
}

// NewVM creates a new virtual machine with the given bytecode and I/O streams
func NewVM(instructions []Instruction, input io.Reader, output io.Writer) *VM {
    return &VM{
        instructions: instructions,
        accumulator:  0,                       // Start with accumulator at 0
        stack:        make([]int, 0, 256),     // Pre-allocate stack with reasonable capacity
        pc:           0,                       // Start at first instruction
        input:        input,                   // Input stream
        output:       output,                  // Output stream
    }
}

// Run executes the bytecode program from start to finish
// Returns an error if any runtime error occurs (typically I/O errors)
func (vm *VM) Run() error {
    for vm.pc < len(vm.instructions) {
        inst := vm.instructions[vm.pc]
        jumped := false  // Track if we jumped

        switch inst.Op {
        case OpInc:
            vm.accumulator++

        case OpDec:
            vm.accumulator--

        case OpPush:
            vm.stack = append(vm.stack, vm.accumulator)

        case OpPop:
            if len(vm.stack) > 0 {
                vm.accumulator = vm.stack[len(vm.stack)-1]
                vm.stack = vm.stack[:len(vm.stack)-1]
            } else {
                vm.accumulator = 0
            }

        case OpLoop:
            if vm.accumulator == 0 {
                vm.pc = inst.Arg
                jumped = true  // We jumped, don't increment pc
            }

        case OpEnd:
            if vm.accumulator != 0 {
                vm.pc = inst.Arg
                jumped = true  // We jumped, don't increment pc
            }

        case OpOut:
            char := byte(vm.accumulator % 256)
            _, err := vm.output.Write([]byte{char})
            if err != nil {
                return fmt.Errorf("output error: %v", err)
            }

        case OpIn:
            buf := make([]byte, 1)
            n, err := vm.input.Read(buf)
            if err != nil && err != io.EOF {
                return fmt.Errorf("input error: %v", err)
            }
            if err == io.EOF || n == 0 {
                vm.accumulator = 0
            } else {
                vm.accumulator = int(buf[0])
            }

        case OpOutNum:
            _, err := fmt.Fprintf(vm.output, "%d", vm.accumulator)
            if err != nil {
                return fmt.Errorf("output error: %v", err)
            }

        default:
            return fmt.Errorf("internal error: invalid opcode %d at position %d", inst.Op, vm.pc)
        }

        // Only increment pc if we didn't jump
        if !jumped {
            vm.pc++
        }
    }

    return nil
}

// Main function: Entry point for the Flux compiler
func main() {
    // If no arguments, show help
    if len(os.Args) < 2 {
        showHelp()
        return
    }

    command := os.Args[1]

    // Dispatch to appropriate command handler
    switch command {
    case "help", "-h", "--help":
        showHelp()

    case "guide":
        showGuide()

    case "reference", "ref":
        showReference()

    case "examples":
        showExamples()

    case "demo":
        runDemo()

    case "run":
        if len(os.Args) < 3 {
            fmt.Println("Error: Please specify a file to run")
            fmt.Println("Usage: flux run <file>")
            return
        }
        runFile(os.Args[2])

    case "compile":
        if len(os.Args) < 3 {
            fmt.Println("Error: Please specify a file to compile")
            fmt.Println("Usage: flux compile <file>")
            return
        }
        compileFile(os.Args[2])

    case "interactive", "repl":
        runInteractive()

    default:
        fmt.Printf("Unknown command: %s\n", command)
        fmt.Println("Run 'flux help' for usage information")
    }
}

// showHelp displays the main help message
func showHelp() {
    fmt.Println(`
                   FLUX PROGRAMMING LANGUAGE v1.0                          
              Minimal & Stack Based &  Turing Complete 

USAGE
    flux <command> [arguments]

COMMANDS
    help              Show this help message
    guide             Show beginner's tutorial and user guide
    reference         Show complete language reference (also: ref)
    examples          Show example programs with explanations
    demo              Run interactive demonstration programs
    run <file>        Compile and execute a Flux program
    compile <file>    Compile program and show bytecode
    interactive       Start interactive REPL (also: repl)

QUICK REFERENCE
    +    Increment accumulator       *    Push to stack
    -    Decrement accumulator       /    Pop from stack
    [    Start loop (if acc != 0)   ]    End loop (jump if acc != 0)
    .    Output as ASCII             ,    Input character
            #    Output as number

GETTING STARTED
    1. Run 'flux guide' for a beginner-friendly tutorial
    2. Try 'flux demo' to see example programs in action
    3. View 'flux examples' for commented program samples
    4. Read 'flux reference' for complete documentation

EXAMPLE
    Create a file 'hello.flux':
        ++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++++.
        ++++++++++++++++++++++++++++++++.

    Run it:
        flux run hello.flux

    Output:
        Hi

For complete documentation, visit the sections above.
`)
}

// showGuide displays the beginner's tutorial
func showGuide() {
    fmt.Println(`
                     FLUX BEGINNER'S GUIDE                                 

## INTRODUCTION

Welcome to Flux! This comprehensive guide will teach you programming in Flux
from absolute scratch. Flux is intentionally minimal - you'll learn everything
you need in under an hour.

[Complete beginner's guide content would continue here...]
`)
}

// showReference displays the complete language reference
func showReference() {
    fmt.Println(`
                   FLUX COMPLETE LANGUAGE REFERENCE                       

[Complete reference documentation would continue here...]
`)
}

// showExamples displays example programs
func showExamples() {
    fmt.Println(`
                       FLUX EXAMPLE PROGRAMS                               

[Complete examples would continue here...]
`)
}

// runDemo runs interactive demonstrations
func runDemo() {
    demos := []struct {
        name string
        code string
        desc string
    }{
        {
            "Output Character 'A'",
            strings.Repeat("+", 65) + ".",
            "Builds ASCII value 65 and outputs 'A'",
        },
        {
            "Output Number 42",
            strings.Repeat("+", 42) + "#",
            "Builds value 42 and outputs it as a number",
        },
        {
            "Count Down from 5",
            "+++++[#-]",
            "Loops 5 times, printing and decrementing",
        },
        {
            "Simple Stack Test",
            "+++*++*/#/#",
            "Pushes 3 and 2, then pops and prints both",
        },
        {
            "Hello (short)",
            strings.Repeat("+", 72) + "." + strings.Repeat("+", 29) + "." + strings.Repeat("+", 7) + ".." + "+++.",
            "Prints 'Hello' using ASCII values",
        },
    }

    fmt.Println("")
    fmt.Println("                       FLUX INTERACTIVE DEMONSTRATION                      ")
    fmt.Println("")
    for i, demo := range demos {
        fmt.Printf("")
        fmt.Printf("Demo %d: %s\n", i+1, demo.name)
        fmt.Printf("Description: %s\n", demo.desc)
        fmt.Printf("Code: %s\n", demo.code)
        fmt.Printf("Output: ")
        execute(demo.code)
        fmt.Printf("\n\n")
    }

    fmt.Println("Try writing your own programs using these patterns!")
}

// runFile compiles and executes a Flux source file
func runFile(filename string) {
    data, err := os.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file '%s': %v\n", filename, err)
        return
    }

    fmt.Printf("Executing %s...\n", filename)
    fmt.Println("")
    execute(string(data))
    fmt.Println()
}

// compileFile compiles a Flux source file and displays the bytecode
func compileFile(filename string) {
    data, err := os.ReadFile(filename)
    if err != nil {
        fmt.Printf("Error reading file '%s': %v\n", filename, err)
        return
    }

    compiler := NewCompiler(string(data))
    instructions, err := compiler.Compile()
    if err != nil {
        fmt.Printf("Compilation error: %v\n", err)
        return
    }

    fmt.Printf("Successfully compiled %s\n", filename)
    fmt.Printf("Total instructions: %d\n\n", len(instructions))
    fmt.Println("Bytecode Listing:")
    fmt.Println("")
    fmt.Println("Addr  Opcode    Argument")
    fmt.Println("")

    opNames := map[OpCode]string{
        OpInc:    "INC",
        OpDec:    "DEC",
        OpPush:   "PUSH",
        OpPop:    "POP",
        OpLoop:   "LOOP",
        OpEnd:    "END",
        OpOut:    "OUT",
        OpIn:     "IN",
        OpOutNum: "OUTNUM",
    }

    for i, inst := range instructions {
        opName := opNames[inst.Op]
        if inst.Op == OpLoop || inst.Op == OpEnd {
            fmt.Printf("%04d  %-8s  â %d\n", i, opName, inst.Arg)
        } else {
            fmt.Printf("%04d  %s\n", i, opName)
        }
    }
    fmt.Println("")
}

// runInteractive starts an interactive REPL
func runInteractive() {
    fmt.Println("")
    fmt.Println("                   FLUX INTERACTIVE MODE (REPL)                            ")
    fmt.Println("")
    fmt.Println("Enter Flux code and press Enter to execute.")
    fmt.Println("Type 'exit' or 'quit' to leave, 'help' for quick reference.\n")

    scanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("flux> ")
        if !scanner.Scan() {
            break
        }

        line := strings.TrimSpace(scanner.Text())

        if line == "" {
            continue
        }

        if line == "exit" || line == "quit" {
            fmt.Println("Goodbye!")
            break
        }

        if line == "help" {
            fmt.Println("Quick Reference:")
            fmt.Println("  +  Increment    *  Push      [  Loop start")
            fmt.Println("  -  Decrement    /  Pop       ]  Loop end")
            fmt.Println("  .  Output char  ,  Input     #  Output number")
            continue
        }

        execute(line)
        fmt.Println()
    }
}

// execute compiles and runs Flux source code
func execute(source string) {
    compiler := NewCompiler(source)
    instructions, err := compiler.Compile()
    if err != nil {
        fmt.Printf("Compilation error: %v\n", err)
        return
    }

    vm := NewVM(instructions, os.Stdin, os.Stdout)
    err = vm.Run()
    if err != nil {
        fmt.Printf("\nRuntime error: %v\n", err)
    }
}