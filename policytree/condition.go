package policytree

import (
	"fmt"
	"github.com/ttdung/matchingAccessTree/attributes"
	"strings"
)

type ConditionType int

const (
	FieldOnly ConditionType = iota + 1
	Comparation
	NumberRange
	DateRange
)

type Condition struct {
	Content       string
	conditionType ConditionType
}

var comparison_characters = []string{"<", ">", "<=", ">=", "="}

var range_characters = []string{"-", "in"}

func NewCondition(s string) *Condition {
	s = strings.TrimSpace(s)
	return &Condition{
		Content: s,
	}
}

// Return Value, isError
func (c *Condition) SolveConditionString(attribute attributes.Attribute) (bool, error, bool) {
	if c.Content == "" {
		return false, nil, false
	}

	//fmt.Println("Solving condition: ", c.Content)

	conditionType, err := c.findConditionType()
	if err != nil {
		return false, err, false
	}
	c.conditionType = conditionType

	//fmt.Println("\tCondition Type: ", c.conditionType)

	var value bool

	switch conditionType {
	case FieldOnly:
		value, err = c.solveFieldOnly(attribute)
	case Comparation:
		value, err = c.solveComperation(attribute)
	case NumberRange:
		value, err = c.solveNumberRange(attribute)
	case DateRange:
		value, err = c.solveDateRange(attribute)
	default:
		value, err = false, fmt.Errorf("CAN'T MATCH CONDITION TYPE")
	}

	return value, err, true
}

func (c *Condition) solveFieldOnly(attr attributes.Attribute) (bool, error) {
	return attr.ContainsAttribute(c.Content), nil
}

func compare(attr attributes.Attribute, fieldStr string, valueStr *string, compFunction func(int, int) bool) (bool, error) {
	fieldStr = strings.TrimSpace(fieldStr)
	*valueStr = strings.TrimSpace((*valueStr))

	fieldValue, isContain := attr.GetAttribute(fieldStr)
	if !isContain {
		return false, nil
	}

	value, err := attributes.ParseValue(*valueStr)
	if err != nil {
		return false, err
	}

	//fmt.Printf("\tField Value: %d, Value: %d\n", fieldValue, value)
	return compFunction(fieldValue, value), nil
}

func (c *Condition) solveComperation(attr attributes.Attribute) (bool, error) {
	if strings.Contains(c.Content, "<=") {
		splitPatterns := strings.Split(c.Content, "<=")
		return compare(attr, splitPatterns[0], &splitPatterns[1], func(a int, b int) bool { return a <= b })

	} else if strings.Contains(c.Content, ">=") {
		splitPatterns := strings.Split(c.Content, ">=")
		return compare(attr, splitPatterns[0], &splitPatterns[1], func(a int, b int) bool { return a >= b })

	} else if strings.Contains(c.Content, "<") {
		splitPatterns := strings.Split(c.Content, "<")
		return compare(attr, splitPatterns[0], &splitPatterns[1], func(a int, b int) bool { return a < b })

	} else if strings.Contains(c.Content, ">") {
		splitPatterns := strings.Split(c.Content, ">")
		return compare(attr, splitPatterns[0], &splitPatterns[1], func(a int, b int) bool { return a > b })

	} else if strings.Contains(c.Content, "=") {
		splitPatterns := strings.Split(c.Content, "=")
		return compare(attr, splitPatterns[0], &splitPatterns[1], func(a int, b int) bool { return a == b })
	}

	return false, fmt.Errorf("CAN'T MATCH COMPARATION TYPE")
}

func (c *Condition) solveNumberRange(attr attributes.Attribute) (bool, error) {
	splits := strings.Split(c.Content, "in")
	for idx, pattern := range splits {
		splits[idx] = strings.TrimSpace(pattern)
	}

	field := splits[0]

	rangeStr := splits[1]
	rangeStr = strings.TrimSpace(rangeStr)
	rangeStr = strings.ReplaceAll(rangeStr, "(", "")
	rangeStr = strings.ReplaceAll(rangeStr, ")", "")

	rangeSplits := strings.Split(rangeStr, "-")
	for idx, pattern := range rangeSplits {
		rangeSplits[idx] = strings.TrimSpace(pattern)
	}

	//fmt.Printf("attribute %s: %s - %s\n", field, rangeSplits[0], rangeSplits[1])

	compLeft, err1 := compare(attr, field, &rangeSplits[0], func(a int, b int) bool { return a > b })
	compRight, err2 := compare(attr, field, &rangeSplits[1], func(a int, b int) bool { return a < b })

	if err1 != nil {
		return false, err1
	}

	if err2 != nil {
		return false, err2
	}

	// check rangeSplits[0] < rangeSplits[1]
	val1, _ := attributes.ParseValue(rangeSplits[0])
	val2, _ := attributes.ParseValue(rangeSplits[1])
	if val1 > val2 {
		return false, fmt.Errorf("LEFT VALUE IS BIGGER THAN RIGHT VALUE: %s", c.Content)
	}

	return compLeft && compRight, nil
}

func (c *Condition) solveDateRange(attr attributes.Attribute) (bool, error) {
	splits := strings.Split(c.Content, "=")
	for idx, pattern := range splits {
		splits[idx] = strings.TrimSpace(pattern)
	}

	field := splits[0]
	val, isContain := attr.GetAttribute(field)
	if !isContain {
		return false, nil
	}

	rangeStr := splits[1]
	rangeStr = strings.TrimSpace(rangeStr)

	val1, val2, err := attributes.ParseDateRange(rangeStr)
	if err != nil {
		return false, err
	}

	if val1 > val2 {
		return false, fmt.Errorf("LEFT VALUE IS BIGGER THAN RIGHT VALUE, %s", c.Content)
	}

	//fmt.Print("\tField Value: ", val, "\n")
	//fmt.Print("\tDate Range: ", val1, " - ", val2, "\n")

	return (val >= val1 && val <= val2), nil
}

func (c *Condition) findConditionType() (ConditionType, error) {
	if c.isOnlyField() {
		return ConditionType(FieldOnly), nil
	}
	if c.isNumberRange() {
		return ConditionType(NumberRange), nil
	}
	if c.isDateRange() {
		return ConditionType(DateRange), nil
	}
	if c.isComparation() {
		return ConditionType(Comparation), nil
	}

	return 0, fmt.Errorf("CONDITION WITH WRONG FORMAT")
}

func (c *Condition) isOnlyField() bool {

	return len(strings.Split(c.Content, " ")) == 1
}

func (c *Condition) isComparation() bool {
	for _, str := range comparison_characters {
		splitPatterns := strings.Split(c.Content, str)
		for idx, pattern := range splitPatterns {
			splitPatterns[idx] = strings.TrimSpace(pattern)
		}

		if len(splitPatterns) == 2 && len(strings.Split(splitPatterns[0], " ")) == 1 {
			return true
		}
	}

	return false
}

func (c *Condition) isNumberRange() bool {
	splits := strings.Split(c.Content, "in")
	for idx, pattern := range splits {
		splits[idx] = strings.TrimSpace(pattern)
	}

	if len(splits) != 2 {
		return false
	}

	// left of in is only one field
	if len(strings.Split(splits[0], " ")) != 1 {
		return false
	}

	// right is seperate 2 pattern with -
	if len(strings.Split(splits[1], "-")) != 2 {
		return false
	}

	return true
}

func (c *Condition) isDateRange() bool {
	splits := strings.Split(c.Content, "=")
	for idx, pattern := range splits {
		splits[idx] = strings.TrimSpace(pattern)
	}

	if len(splits) != 2 {
		return false
	}

	// left of in is only one field
	if len(strings.Split(splits[0], " ")) != 1 {
		return false
	}

	// right is seperate 2 pattern with -
	if len(strings.Split(splits[1], "-")) != 2 {
		return false
	}

	return true
}
