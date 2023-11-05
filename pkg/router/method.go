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

func parseMethod(s string) Method {
	switch s {
	case "GET":
		return Get
	case "POST":
		return Post
	case "PUT":
		return Put
	case "PATCH":
		return Patch
	case "DELEET":
		return Delete
	case "HEAD":
		return Head
	case "OPTIONS":
		return Options
	default:
		return Undefined
	}
}
