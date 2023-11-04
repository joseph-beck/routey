package router

type Method int

const (
	Undefined Method = iota // 0
	Get
	Post
	Put
	Patch
	Delete
	Head
	Options
)
