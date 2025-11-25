#FLUX PROGRAMMING LANGUAGE - COMPLETE REFERENCE GUIDE

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


'+'    Increment accumulator

     Effect: accumulator = accumulator + 1
     Example: If acc=5, then after '+' acc=6

     

'-'    Decrement accumulator

     Effect: accumulator = accumulator - 1
     Example: If acc=5, then after '-' acc=4

     

STACK OPERATIONS:


'*'    Push accumulator to stack

     Effect: stack.push(accumulator)
     The accumulator value is copied to the top of the stack
     The accumulator itself remains unchanged
     Example: If acc=5, after '*' the stack has 5 on top and acc is still 5

     

'/'    Pop stack to accumulator

     Effect: accumulator = stack.pop()
     If stack is empty, accumulator becomes 0
     The value is removed from the stack
     Example: If stack=[3,7] (7 on top) and acc=5, after '/' acc=7 and stack=[3]

     

CONTROL FLOW OPERATIONS:


'['    Begin while loop

     If accumulator == 0, jump forward past the matching ']'
     If accumulator != 0, continue execution into the loop body
     The loop condition is checked only at the '[', not continuously
     Example: If acc=0, execution jumps past the loop
     If acc=5, execution enters the loop

     

']'    End while loop

     If accumulator != 0, jump backward to the matching '['
     If accumulator == 0, continue execution past the loop
     This creates a while-loop structure that continues as long as acc != 0
     Example: If acc=3, jump back to '['
     If acc=0, exit loop

     

INPUT/OUTPUT OPERATIONS:


'.'    Output accumulator as ASCII character

     Prints the character corresponding to (accumulator mod 256)
     Example: If acc=65, outputs 'A'
     If acc=72, outputs 'H'

     

','    Input one character

     Reads a single character from input stream
     Sets accumulator to the ASCII value of that character
     On EOF (end of file), sets accumulator to 0
     Example: If user types 'A', acc becomes 65

     

'#'    Output accumulator as decimal number

     Prints the numeric value of the accumulator
     This is an extension for practical debugging and numeric output
     Example: If acc=42, outputs "42"

     

WHITESPACE AND COMMENTS:


- Spaces, tabs, newlines, and carriage returns are ignored
- Any character that is not one of the 9 operations is treated as a comment
- This allows for readable, documented Flux code


END OF REFERENCE GUIDE
