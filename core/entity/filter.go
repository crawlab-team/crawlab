package entity

import (
	"github.com/crawlab-team/crawlab/core/interfaces"
	"reflect"
)

type Condition struct {
	Key   string      `json:"key"`
	Op    string      `json:"op"`
	Value interface{} `json:"value"`
}

func (c *Condition) GetKey() (key string) {
	return c.Key
}

func (c *Condition) SetKey(key string) {
	c.Key = key
}

func (c *Condition) GetOp() (op string) {
	return c.Op
}

func (c *Condition) SetOp(op string) {
	c.Op = op
}

func (c *Condition) GetValue() (value interface{}) {
	return c.Value
}

func (c *Condition) SetValue(value interface{}) {
	c.Value = value
}

type Filter struct {
	IsOr       bool         `form:"is_or" url:"is_or"`
	Conditions []*Condition `json:"conditions"`
}

func (f *Filter) GetIsOr() (isOr bool) {
	return f.IsOr
}

func (f *Filter) SetIsOr(isOr bool) {
	f.IsOr = isOr
}

func (f *Filter) GetConditions() (conditions []interfaces.FilterCondition) {
	for _, c := range f.Conditions {
		conditions = append(conditions, c)
	}
	return conditions
}

func (f *Filter) SetConditions(conditions []interfaces.FilterCondition) {
	f.Conditions = make([]*Condition, len(conditions))
	for _, c := range conditions {
		f.Conditions = append(f.Conditions, c.(*Condition))
	}
}

func (f *Filter) IsNil() (ok bool) {
	val := reflect.ValueOf(f)
	return val.IsNil()
}
