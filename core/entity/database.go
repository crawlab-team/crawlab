package entity

type DatabaseMetadata struct {
	Databases []Database `json:"databases"`
}

type Database struct {
	Name   string          `json:"name"`
	Tables []DatabaseTable `json:"tables"`
}

type DatabaseTable struct {
	Name    string           `json:"name"`
	Columns []DatabaseColumn `json:"columns"`
	Indexes []DatabaseIndex  `json:"indexes"`
}

type DatabaseColumn struct {
	Name     string           `json:"name"`
	Type     string           `json:"type"`
	Null     bool             `json:"null,omitempty"`
	Key      string           `json:"key,omitempty"`
	Default  string           `json:"default,omitempty"`
	Extra    string           `json:"extra,omitempty"`
	Children []DatabaseColumn `json:"children,omitempty"`
}

type DatabaseIndex struct {
	Name    string                `json:"name"`
	Type    string                `json:"type,omitempty"`
	Columns []DatabaseIndexColumn `json:"columns"`
	Unique  bool                  `json:"unique"`
}

type DatabaseIndexColumn struct {
	Name  string `json:"name"`
	Order int    `json:"order"`
}
