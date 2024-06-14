package entity

type Color struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
}

func (c *Color) GetHex() string {
	return c.Hex
}

func (c *Color) GetName() string {
	return c.Name
}

func (c *Color) Value() interface{} {
	return c
}
