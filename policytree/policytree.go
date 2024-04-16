package policytree

import (
	"fmt"
	"github.com/ttdung/matchingAccessTree/attributes"
	"github.com/ttdung/matchingAccessTree/utils"
	"log"
	"os"
	"unicode/utf8"
)

var and_characters = []string{"and", "AND"}

var or_characters = []string{"or", "OR"}

type PolicyTree struct {
	Content string
}

func NewPolicyTree(s string) *PolicyTree {
	return &PolicyTree{
		Content: s,
	}
}

func isAndOperation(s string) bool {
	for _, value := range and_characters {
		if s == value {
			return true
		}
	}

	return false
}

func isOrOperation(s string) bool {
	for _, value := range or_characters {
		if s == value {
			return true
		}
	}
	return false
}

func isOperationString(s string) bool {
	return isAndOperation(s) || isOrOperation(s)
}

func solveConditionString(conditionString string, attribute attributes.Attribute) (bool, bool) {
	conditionObj := NewCondition(conditionString)

	var val bool
	var err error
	var isPush bool

	val, err, isPush = conditionObj.SolveConditionString(attribute)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return val, isPush
}

func tryToSolveOperation(operation string, valueStack *utils.Stack) (bool, error) {
	if valueStack.Length() < 2 {
		return false, fmt.Errorf("NOT ENOUGH VALUES IN STACK TO SOLVE OPERATION %s", operation)
	}

	value1 := valueStack.Pop().(bool)
	value2 := valueStack.Pop().(bool)
	value := false

	if isAndOperation(operation) {
		value = value1 && value2
	} else {
		value = value1 || value2
	}

	return value, nil
}

func handleExpression(operationStack *utils.Stack, valueStack *utils.Stack) {
	for operationStack.Length() > 0 && operationStack.Peek() != '(' {
		operation := operationStack.Pop().(string)
		value, error := tryToSolveOperation(operation, valueStack)
		if error != nil {
			fmt.Println(error)
			os.Exit(1)
		} else {
			valueStack.Push(value)
		}
	}

	if operationStack.Length() > 0 {
		operationStack.Pop()
	}
}

func solveWithEndWord(operationStack *utils.Stack, valueStack *utils.Stack, buffer string, conditionString string, attributes attributes.Attribute) (string, string) {
	if isOperationString(buffer) { // it is AND or OR
		// fmt.Println("\t", buffer, " is an operation")
		value, isPush := solveConditionString(conditionString, attributes)

		if isPush {
			valueStack.Push(value)
		}

		operationStack.Push(buffer)

		conditionString = ""
		buffer = ""
	} else { // it is not AND or OR
		conditionString = conditionString + " " + buffer
		buffer = ""
	}

	return buffer, conditionString
}

func solveWithEndCondition(valueStack *utils.Stack, conditionString string, attributes attributes.Attribute) string {
	value, ok := solveConditionString(conditionString, attributes)
	if ok {
		valueStack.Push(value)
	}

	conditionString = ""

	return conditionString
}

func (pTree *PolicyTree) EvaluatePolicyTree(attributes attributes.Attribute) bool {
	valueStack := utils.NewStack()
	operationStack := utils.NewStack()

	s := pTree.Content

	// for idx, runeValue := range s {
	// 	fmt.Printf("%#U starts at %d\n", runeValue, idx)
	// }

	conditionString := ""
	buffer := ""

	//fmt.Println("\nUsing DecodeRuneInString")
	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		w = width
		// fmt.Printf("%#U starts at %d\n", runeValue, i)

		if runeValue == '(' {
			//fmt.Printf("( at position %d\n", i)
			operationStack.Push('(')

		} else if runeValue == ')' {
			//fmt.Printf(") at position %d, buffer = %s, conditionString = %s\n", i, buffer, conditionString)

			buffer, conditionString = solveWithEndWord(operationStack, valueStack, buffer, conditionString, attributes)
			conditionString = solveWithEndCondition(valueStack, conditionString, attributes)

			handleExpression(operationStack, valueStack)

			//fmt.Println("\tvalueStack: ", valueStack)
			//fmt.Println("\toperationStack: ", operationStack)
		} else if runeValue == ' ' {
			//fmt.Printf("Space, buffer = %s, condition string = %s\n", buffer, conditionString)
			buffer, conditionString = solveWithEndWord(operationStack, valueStack, buffer, conditionString, attributes)
			//fmt.Println("\tvalueStack: ", valueStack)
			//fmt.Println("\toperationStack: ", operationStack)
		} else { // it is a character
			buffer = buffer + string(runeValue)
		}
	}

	handleExpression(operationStack, valueStack)
	return valueStack.Pop().(bool)
}
