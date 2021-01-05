package model

// Kratos hello kratos.
type Kratos struct {
	Hello string
}

type Article struct {
	ID      int64 `json:",string"`
	Content string
	Author  string
}
