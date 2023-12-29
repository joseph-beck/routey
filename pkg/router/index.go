package router

type State int

const (
	Healthy   State = iota // 0
	UnHealthy              // 1
	Aborted                // 2
	Restored               // 3
)
