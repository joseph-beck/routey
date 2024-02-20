package websockets

type Config struct {
	Origins []string

	ReadBufferSize  int
	WriteBufferSize int

	RecoverHandler func(*Conn)
}

func Default() Config {
	return Config{
		Origins: []string{},

		ReadBufferSize: 1024,
		WriteBufferSize: 1024,

		RecoverHandler: defaultRecover,
	}
}
