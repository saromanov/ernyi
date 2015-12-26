package event

type Event interface {
	String() string
}

type EventStop struct {
}

func (evn *EventStop) String() string {
	return "stop"
}
