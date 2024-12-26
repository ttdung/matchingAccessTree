package attributes

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// Attribute is a struct that holds a list of attributes
type Attribute struct {
	List map[string]int `json "list"`
}

func (a *Attribute) ToString() string {
	b := new(bytes.Buffer)

	for key, value := range a.List {
		fmt.Fprintf(b, "%s = %d\n", key, value)
	}
	return b.String()
}

func NewAttribute() *Attribute {
	return &Attribute{}
}

func ParseValue(s string) (int, error) {
	s = strings.TrimSpace(s)

	// check if s is an int
	var valueInt int
	_, err := fmt.Sscanf(s, "%d", &valueInt)
	if err == nil {
		return valueInt, nil
	}

	// check if s is a date
	date := NewDateFromString(s)
	if date != nil {
		return date.DateToInt(), nil
	}

	return 0, fmt.Errorf("CAN'T PARSE VALUE, VALUE %s IS NOT FORMATTED CORRECTLY", s)
}

func ParseDateRange(s string) (int, int, error) {
	var month string
	var dayLeft int
	var dayRight int
	var year int

	s = strings.TrimSpace(s)
	// TODO: space is not important
	// _, err := fmt.Sscanf(s, "%s %d - %d, %d", &month, &dayLeft, &dayRight, &year);

	// if err != nil {
	// 	return 0, 0, fmt.Errorf("CAN'T PARSE DATE RANGE, VALUE %s IS NOT FORMATTED CORRECTLY", s);
	// }

	words := seperateWordsFromString(s)
	if len(words) != 4 {
		return 0, 0, fmt.Errorf("CAN'T PARSE DATE RANGE, VALUE %s IS NOT FORMATTED CORRECTLY", s)
	}

	month = words[0]
	dayLeft, err := strconv.Atoi(words[1])
	if err != nil {
		return 0, 0, fmt.Errorf("CAN'T PARSE DATE RANGE, VALUE %s IS NOT FORMATTED CORRECTLY", s)
	}

	dayRight, err = strconv.Atoi(words[2])
	if err != nil {
		return 0, 0, fmt.Errorf("CAN'T PARSE DATE RANGE, VALUE %s IS NOT FORMATTED CORRECTLY", s)
	}

	year, err = strconv.Atoi(words[3])
	if err != nil {
		return 0, 0, fmt.Errorf("CAN'T PARSE DATE RANGE, VALUE %s IS NOT FORMATTED CORRECTLY", s)
	}

	date1 := NewDate(dayLeft, month, year)
	date2 := NewDate(dayRight, month, year)

	return date1.DateToInt(), date2.DateToInt(), nil
}

func parseAttribute(s string) (string, int, error) {
	splitSlice := strings.Split(s, "=")
	field := splitSlice[0]
	field = strings.TrimSpace(field)
	//fmt.Println("field = ", field);

	value := 0
	if len(splitSlice) > 1 {
		valueString := splitSlice[1]
		val, err := ParseValue(valueString)
		if err != nil {
			return "", 0 , err
		}

		value = val
	}

	//fmt.Println("value = ", value);
	return field, value, nil
}

func NewAttributeFromString(s string) (*Attribute, error) {
	splitSlice := strings.Split(s, "|")

	list := make(map[string]int)
	for _, v := range splitSlice {
		v = strings.TrimSpace(v)
		if len(v) == 0 {
			continue
		}

		field, value, err := parseAttribute(v)
		if err != nil {
			return nil, err
		}

		list[field] = value
	}
	return &Attribute{
		List: list,
	}, nil
}

func (a *Attribute) SetAttribute(lable string, value int) {
	a.List[lable] = value
}

func (a *Attribute) ContainsAttribute(lable string) bool {
	_, ok := a.List[lable]
	return ok
}

func (a *Attribute) RemoveAttribute(lable string) {
	delete(a.List, lable)
}

func (a Attribute) GetAttribute(lable string) (int, bool) {
	val, isContain := a.List[lable]
	return val, isContain
}
