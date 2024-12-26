package policytree

import (
	"fmt"
	"github.com/ttdung/matchingAccessTree/attributes"
	"github.com/ttdung/matchingAccessTree/utils"
	"unicode/utf8"
)

var and_characters = []string{"and", "AND"}

var or_characters = []string{"or", "OR"}

type PolicyTree struct {
	Content string
}

func NewPolicyTree(s string) *PolicyTree {
	sExtend := "(" + s + ")"
	return &PolicyTree{
		Content: sExtend,
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

func solveConditionString(conditionString string, attribute attributes.Attribute) (bool, bool, error) {
	conditionObj := NewCondition(conditionString)

	val, is_push, err := conditionObj.SolveConditionString(attribute)
	if err != nil {
		return false, false, err
	}

	return val, is_push, nil
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

func handleExpression(operationStack *utils.Stack, valueStack *utils.Stack) error {
	for operationStack.Length() > 0 && operationStack.Peek() != '(' {
		operation := operationStack.Pop().(string)
		value, err := tryToSolveOperation(operation, valueStack)
		if err != nil {
			return err
		} else {
			valueStack.Push(value)
		}
	}

	if operationStack.Length() > 0 {
		operationStack.Pop()
	}

	return nil
}

func solveWithEndWord(operationStack *utils.Stack, valueStack *utils.Stack, buffer string, conditionString string, attributes attributes.Attribute) (string, string, error) {
	if isOperationString(buffer) { // it is AND or OR
		// fmt.Println("\t", buffer, " is an operation")
		value, is_push, err := solveConditionString(conditionString, attributes)
		if err != nil {
			return "", "", nil
		}

		if is_push {
			valueStack.Push(value)
		}

		operationStack.Push(buffer)

		conditionString = ""
		buffer = ""
	} else { // it is not AND or OR
		conditionString = conditionString + " " + buffer
		buffer = ""
	}

	return buffer, conditionString, nil
}

func solveWithEndCondition(valueStack *utils.Stack, conditionString string, attributes attributes.Attribute) (string, error) {
	value, ok, err := solveConditionString(conditionString, attributes)
	if err != nil {
		return "", err
	}

	if ok {
		valueStack.Push(value)
	}

	return "", nil
}

func (pTree *PolicyTree) EvaluatePolicyTree(attributes attributes.Attribute) (bool, error) {
	valueStack := utils.NewStack()
	operationStack := utils.NewStack()

	s := pTree.Content

	// for idx, runeValue := range s {
	// 	fmt.Printf("%#U starts at %d\n", runeValue, idx)
	// }

	var conditionString string
	var buffer string
	var err error

	//fmt.Println("\nUsing DecodeRuneInString")
	for i, w := 0, 0; i < len(s); i += w {
		runeValue, width := utf8.DecodeRuneInString(s[i:])
		w = width

		if runeValue == '(' {
			operationStack.Push('(')
		} else if runeValue == ')' {
			buffer, conditionString, err = solveWithEndWord(operationStack, valueStack, buffer, conditionString, attributes)
			if err != nil {
				return false, err
			}

			conditionString, err = solveWithEndCondition(valueStack, conditionString, attributes)
			if err != nil {
				return false, err
			}

			err = handleExpression(operationStack, valueStack)
			if err != nil {
				return false, err
			}
		} else if runeValue == ' ' {
			buffer, conditionString, err = solveWithEndWord(operationStack, valueStack, buffer, conditionString, attributes)
			if err != nil {
				return false, err
			}
		} else { // it is a character
			buffer = buffer + string(runeValue)
		}
	}

	err = handleExpression(operationStack, valueStack)
	if err != nil {
		return false, err
	}

	return valueStack.Pop().(bool), nil
}
