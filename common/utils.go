package common

import (
	"github.com/ttdung/matchingAccessTree/attributes"
	"github.com/ttdung/matchingAccessTree/policytree"
)

func EvaluatePolicyTree(attributeStr string, policyTreeStr string) bool {
	attribute := attributes.NewAttributeFromString(attributeStr)
	policyTree := policytree.NewPolicyTree(policyTreeStr)
	return policyTree.EvaluatePolicyTree(*attribute)
}
