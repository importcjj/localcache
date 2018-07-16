package singlecache

const (
	EvtKindTimeout  = "TIMEOUT"
	EvtKindComplete = "COMPLETE"
)

type (
	event struct {
		Key  RequestKey
		Kind string
	}

	getEvent struct {
		Key     RequestKey
		request chan *requestGroup
	}
)
