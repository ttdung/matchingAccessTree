package common

import (
	"github.com/ttdung/matchingAccessTree/attributes"
	"github.com/ttdung/matchingAccessTree/policytree"
)

func EvaluatePolicyTree(attributeStr string, policyTreeStr string) (bool, error) {
	attribute, err := attributes.NewAttributeFromString(attributeStr)
	if err != nil {
		return false, err
	}

	policyTree := policytree.NewPolicyTree(policyTreeStr)
	return policyTree.EvaluatePolicyTree(*attribute)
}
