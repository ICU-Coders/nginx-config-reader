package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	Transfer() string
}
type BaseAst struct {
}
type MultiValue struct {
	BaseAst
	Values []string
}
type PropertyExpression struct {
	BaseAst
	Key   string
	Value Node
}
type ObjectExpression struct {
	BaseAst
	Properties []*PropertyExpression
}

func (receiver MultiValue) GetValue() string {
	return fmt.Sprintf("%s", strings.Join(receiver.Values, ""))
}
func (receiver MultiValue) GetValueWithSep(sep string) string {
	return fmt.Sprintf("%s", strings.Join(receiver.Values, sep))
}
func (receiver MultiValue) Transfer() string {
	return fmt.Sprintf("%s", strings.Join(receiver.Values, " "))
}
func (receiver *BaseAst) Transfer() string {
	return ""
}
func (receiver PropertyExpression) Transfer() string {
	return fmt.Sprintf("%s : %s", receiver.Key, receiver.Value.Transfer())
}

func (receiver ObjectExpression) Get(key string) []Node {
	var objs []Node
	for _, property := range receiver.Properties {
		if property.Key == key {
			objs = append(objs, property.Value)
		}
	}
	return objs
}
func (receiver *ObjectExpression) GetMustObject(key string) []*ObjectExpression {
	var objs []*ObjectExpression
	for _, property := range receiver.Properties {
		if property.Key == key {
			objs = append(objs, property.Value.(*ObjectExpression))
		}
	}
	return objs
}
func (receiver *ObjectExpression) GetMustString(key string) []string {
	var objs []string
	for _, property := range receiver.Properties {
		if property.Key == key {
			objs = append(objs, property.Value.(*MultiValue).GetValue())
		}
	}
	return objs
}
func (receiver *ObjectExpression) GetFirst(key string) Node {
	for _, property := range receiver.Properties {
		if property.Key == key {
			return property.Value
		}
	}
	return nil
}
func (receiver *ObjectExpression) GetFirstMustObject(key string) *ObjectExpression {
	for _, property := range receiver.Properties {
		if property.Key == key {
			obj, ok := property.Value.(*ObjectExpression)
			if ok {
				return obj
			}
			return nil
		}
	}
	return nil
}
func (receiver *ObjectExpression) GetFirstMustString(key string) string {
	for _, property := range receiver.Properties {
		if property.Key == key {
			obj, ok := property.Value.(*MultiValue)
			if ok {
				return obj.GetValue()
			}
			return ""
		}
	}
	return ""
}
func (receiver ObjectExpression) Transfer() string {
	var maker []string
	maker = append(maker, "{")
	for _, i2 := range receiver.Properties {
		if len(i2.Key) == 0 {
			continue
		}
		maker = append(maker, i2.Transfer())
	}
	maker = append(maker, "}")
	return strings.Join(maker, "\n")
}

func (receiver *ObjectExpression) Append(node Node) {
	switch node.(type) {
	case *PropertyExpression:
		receiver.Properties = append(receiver.Properties, node.(*PropertyExpression))
	case *ObjectExpression:
		receiver.Properties = append(receiver.Properties, node.(*ObjectExpression).Properties...)
	}
}
