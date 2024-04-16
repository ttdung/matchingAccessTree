package main

import (
	"fmt"
	"github.com/ttdung/matchingAccessTree/common"
)

func main() {
	attributeLists := "|Date = March 11, 2010|Manager|IT|Experience=10|"
	policyTreeString := "((Manager) AND IT) and ((Experience <= 10 or Experience > 100) and Date = March 10 -12, 2010)"

	fmt.Println(common.EvaluatePolicyTree(attributeLists, policyTreeString))
}
