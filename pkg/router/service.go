package router

type Service struct {
}

type Servicer interface {
	Getter
	Poster
	Putter
	Patcher
	Deleter
	Header
	Optioner
}

type Getter interface {
	Get(...any) HandlerFunc
}

type Poster interface {
	Post(...any) HandlerFunc
}

type Putter interface {
	Put(...any) HandlerFunc
}

type Patcher interface {
	Patch(...any) HandlerFunc
}

type Deleter interface {
	Delete(...any) HandlerFunc
}

type Header interface {
	Head(...any) HandlerFunc
}

type Optioner interface {
	Options(...any) HandlerFunc
}
