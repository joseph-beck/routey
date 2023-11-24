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

func (m Method) String() string {
	switch m {
	case Get:
		return "GET"
	case Post:
		return "POST"
	case Put:
		return "PUT"
	case Patch:
		return "PATCH"
	case Delete:
		return "DELETE"
	case Head:
		return "HEAD"
	case Options:
		return "OPTIONS"
	default:
		return "UNDEFINED"
	}
}

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
	case "DELETE":
		return Delete
	case "HEAD":
		return Head
	case "OPTIONS":
		return Options
	default:
		return Undefined
	}
}
