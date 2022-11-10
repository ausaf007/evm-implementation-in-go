package evm

import (
	"errors"
	"github.com/emirpasic/gods/stacks/arraystack"
	log "github.com/sirupsen/logrus"
)

var (
	EmptyStack       = errors.New("stack is empty. cannot pop")
	ConversionFailed = errors.New("unable to convert stack element to byte array")
)

// pushValues pushes
// i+1, i+2, i+3 ..... i+n values from byteArr to localStack
// Suitable to be used with PUSH1, PUSH2, PUSH3, ..... PUSH32
func pushValues(byteArr []byte, localStack *arraystack.Stack, start int, n int) {

	stackElement := make([]byte, 32)

	for j := 1; j <= n; j++ {
		stackElement[32-n+j-1] = byteArr[start+j]
	}
	log.Info("Stack element from PUSH Values=", stackElement)
	localStack.Push(stackElement)
}

// popByte takes input of a stack, pops the element, converts it to []byte datatype, and returns the []byte
func popByte(localStack *arraystack.Stack) ([]byte, error) {
	element, ok := localStack.Pop()
	if !ok {
		return nil, EmptyStack
	}

	elementByteArr, ok := element.([]byte)
	if !ok {
		return nil, ConversionFailed
	}

	return elementByteArr, nil
}
