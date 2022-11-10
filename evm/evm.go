package evm

import (
	"encoding/hex"
	"errors"
	"github.com/emirpasic/gods/stacks/arraystack"
	log "github.com/sirupsen/logrus"
	"math/big"
)

const (
	PUSH1  byte = 96  // Hex=60
	PUSH2  byte = 97  // Hex=61
	PUSH3  byte = 98  // Hex=62
	PUSH32 byte = 127 // Hex=7f

	MSTORE  byte = 82 // Hex=52
	MSTORE8 byte = 83 // Hex=53

	ADD  byte = 1  // Hex=01
	MUL  byte = 2  // Hex=02
	SDIV byte = 5  // Hex=05
	EXP  byte = 10 // Hex=0a
)

var (
	InvalidInstruction  = errors.New("invalid instruction found")
	StackMaxSizeReached = errors.New("can't PUSH element, stack is already full")
)

var stack *arraystack.Stack
var memory []byte
var gasCost int // stores total gas cost

var MOD *big.Int

// byte array of length 32, useful in memory expansion
var memory32 = make([]byte, 32)

// executePUSH is a helper function which executes PUSH1, PUSH2, PUSH3 and PUSH32.
// Throws error when stack is full
func executePUSH(inputArr []byte, start int, n int) (int, error) {

	log.Info("(Success)Found the command PUSH", n)
	log.Info("Stack Size=", stack.Size())

	pushValues(inputArr, stack, start, n)
	if stack.Size() > 1024 {
		// return error, max stack size in EVM is 1024
		return start, StackMaxSizeReached
	}
	gasCost += 3
	start += n
	return start, nil
}

// executeMSTORE takes in offset and value by popping the stack,
// and then stores the 32-byte value
// in the memory, at the location specified by offset
func executeMSTORE() error {
	log.Info("(Success)Found the MSTORE command.")
	gasCost += 3 // base cost

	offset := new(big.Int)

	stackElement1, err := popByte(stack)
	if err != nil {
		return err
	}

	offset.SetBytes(stackElement1)
	offsetInt := int(offset.Int64()) // converting offset(big int format) to int format

	// Handle Resizing of the inputArr, based on offset value
	for offsetInt+32 > len(memory) {
		memory = append(memory, memory32...)
	}

	log.Info("Memory expansion done")

	stackElement2, err := popByte(stack)
	if err != nil {
		return err
	}

	for j := 0; j < 32; j++ {
		memory[offsetInt+j] = stackElement2[j]
	}
	log.Info("MSTORE Execution Done")

	return nil
}

// executeMSTORE8 takes in offset and value by popping the stack,
// and then stores the 1 byte (8-bit) value (the least significant value of the popped stack)
// in the memory, at the location specified by offset
func executeMSTORE8() error {
	log.Info("(Success)Found the MSTORE8 command.")
	gasCost += 3 // base cost

	offset := new(big.Int)
	stackElement1, err := popByte(stack)
	if err != nil {
		return err
	}

	offset.SetBytes(stackElement1)

	offsetInt := int(offset.Int64()) // converting offset(big int format) to int format

	// Resizing of the inputArr, based on offset value
	for offsetInt+1 > len(memory) {
		memory = append(memory, memory32...)
	}

	stackElement2, err := popByte(stack)
	if err != nil {
		return err
	}
	memory[offsetInt] = stackElement2[31]

	return nil
}

// popTwoElements Pops top two elements from the stack
func popTwoElements() (*big.Int, *big.Int, error) {
	a := new(big.Int)
	stackElementA, err := popByte(stack)
	if err != nil {
		return nil, nil, err
	}
	a.SetBytes(stackElementA)

	b := new(big.Int)
	stackElementB, err := popByte(stack)
	if err != nil {
		return nil, nil, err
	}
	b.SetBytes(stackElementB)

	return a, b, nil
}

// executeADD adds the top 2 elements of the stack
// then pops the top 2 elements
// then push sum to the top of the stack
func executeADD() error {
	log.Info("(Success)Found the ADD command.")
	gasCost += 3

	a, b, err := popTwoElements()
	if err != nil {
		return err
	}

	sum := new(big.Int)
	sum.Add(a, b)
	sum.Mod(sum, MOD)
	stack.Push(bigIntToByteArr(sum))

	return nil
}

// executeMUL multiplies (MOD multiply) the top 2 elements of the stack
// then pops the top 2 elements
// then push sum to the top of the stack
func executeMUL() error {
	log.Info("(Success)Found the MUL command.")
	gasCost += 5

	a, b, err := popTwoElements()
	if err != nil {
		return err
	}

	prod := new(big.Int)
	prod.Mul(a, b)
	prod.Mod(prod, MOD)
	stack.Push(bigIntToByteArr(prod))

	return nil
}

// executeSDIV performs MOD division the top 2 elements of the stack
// then pops the top 2 elements
// then push sum to the top of the stack
func executeSDIV() error {
	log.Info("(Success)Found the SDIV command.")
	gasCost += 5

	a, b, err := popTwoElements()
	if err != nil {
		return err
	}

	if b.Cmp(big.NewInt(0)) == 0 {
		stack.Push(bigIntToByteArr(big.NewInt(0)))
	} else {
		quotient := new(big.Int)
		quotient.Div(a, b)
		stack.Push(bigIntToByteArr(quotient))
	}
	return nil
}

// executeSDIV performs MOD exponentiation the top 2 elements of the stack
// then pops the top 2 elements
// then push sum to the top of the stack
func executeEXP() error {
	log.Info("(Success)Found the EXP command.")
	gasCost += 10 // base cost

	a := new(big.Int)
	stackElement1, err := popByte(stack)
	if err != nil {
		return err
	}
	a.SetBytes(stackElement1)

	b := new(big.Int)
	stackElement2, err := popByte(stack)
	if err != nil {
		return err
	}
	b.SetBytes(stackElement2)

	byteLength := len(stackElement2)
	gasCost += 50 * byteLength // variable gas fees

	// calculating a^b mod 2^256
	result := a.Exp(a, b, MOD)

	stack.Push(bigIntToByteArr(result))

	return nil
}

// initialize sets the stack, memory and gasCost to their default values
func initialize() {
	stack = arraystack.New() // create a stack of elements, where each element is a byte array (of length 32)
	memory = nil
	gasCost = 0 // total gas cost

	var base, exp = big.NewInt(2), big.NewInt(256)
	MOD = base.Exp(base, exp, nil) // MOD = 2^256
}

// RunByteCode takes in input bytecode, executes it, and returns the hashed value of memory
// and total gas consumed
func RunByteCode(input string) (string, int, error) {

	// input string to byte array conversion
	inputArr, err := hex.DecodeString(input)
	if err != nil {
		return "", 0, err
	}

	log.Info("Output: ", inputArr)
	initialize()

	// Iterate through the bytecode
	for i := 0; i < len(inputArr); i++ {

		if inputArr[i] == PUSH1 {
			i, err = executePUSH(inputArr, i, 1)
		} else if inputArr[i] == PUSH2 {
			i, err = executePUSH(inputArr, i, 2)
		} else if inputArr[i] == PUSH3 {
			i, err = executePUSH(inputArr, i, 3)
		} else if inputArr[i] == PUSH32 {
			i, err = executePUSH(inputArr, i, 32)
		} else if inputArr[i] == MSTORE {
			err = executeMSTORE()
		} else if inputArr[i] == MSTORE8 {
			err = executeMSTORE8()
		} else if inputArr[i] == ADD {
			err = executeADD()
		} else if inputArr[i] == MUL {
			err = executeMUL()
		} else if inputArr[i] == SDIV {
			err = executeSDIV()
		} else if inputArr[i] == EXP {
			err = executeEXP()
		} else {
			log.Warn("Invalid Instruction Found: ", inputArr[i])
			return "", gasCost, InvalidInstruction
		}

		if err != nil {
			return "", gasCost, err
		}
	}

	gasCost += calcMemGasCost(memory)
	return calcHash(memory), gasCost, err
}
