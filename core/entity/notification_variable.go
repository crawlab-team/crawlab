package entity

import "fmt"

type NotificationVariable struct {
	Category string `json:"category"`
	Name     string `json:"name"`
}

func (v *NotificationVariable) GetKey() string {
	return fmt.Sprintf("${%s:%s}", v.Category, v.Name)
}
