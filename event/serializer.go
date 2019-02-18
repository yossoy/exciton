package event

import "log"

type fncRelayEventToNative func(e *Event)
type fncRelayEventWithResultToNative func(e *Event, callback ResponceCallback)

type command int

const (
	relay command = iota
	relayWithResult
	emit
	emitWithResult
	quit
)

type cmd struct {
	command command
	// name         string
	e                 *Event
	respCallback      ResponceCallback
	quitChan          chan bool
	handler           Handler
	handlerWithResult HandlerWithResult
	//	unmatchedHandler UnmatchedHandler
	// argument interface{}
	// targetRouter     *Router
	// resultChan       chan interface{}
}

type EventSerializer struct {
	handler           fncRelayEventToNative
	handlerWithResult fncRelayEventWithResultToNative
	eventChan         chan cmd
}

var serializerSingleton *EventSerializer

func NewSerializer(handler fncRelayEventToNative, HandlerWithResult fncRelayEventWithResultToNative) *EventSerializer {
	s := &EventSerializer{
		handler:           handler,
		handlerWithResult: HandlerWithResult,
		eventChan:         make(chan cmd),
	}
	serializerSingleton = s
	go s.run()
	return s
}

func (s *EventSerializer) run() {
	for {
		c := <-s.eventChan
		log.Printf("Serializer received: %v", c)
		switch c.command {
		case quit:
			close(s.eventChan)
			c.quitChan <- true
			return
		case relay:
			s.handler(c.e)
		case relayWithResult:
			s.handlerWithResult(c.e, c.respCallback)
		case emit:
			log.Printf("emit: %v", c.e)
			c.handler(c.e)
			if c.respCallback != nil {
				c.respCallback(NewValueResult(nil))
			}
		case emitWithResult:
			log.Printf("emitWithResult: %v", c.e)
			c.handlerWithResult(c.e, c.respCallback)
		}
	}
}

func (s *EventSerializer) RelayEvent(e *Event) {
	c := cmd{
		command: relay,
		e:       e,
	}
	log.Printf("RelayEvent: %v", c)
	s.eventChan <- c
}

func (s *EventSerializer) RelayEventWithResult(e *Event, callback ResponceCallback) {
	c := cmd{
		command:      relayWithResult,
		e:            e,
		respCallback: callback,
	}
	log.Printf("RelayEventWithResult: %v", c)
	s.eventChan <- c
}

func (s *EventSerializer) Stop() {
	qc := make(chan bool)
	c := cmd{
		command:  quit,
		quitChan: qc,
	}
	s.eventChan <- c
	<-qc
	close(qc)
}

func emitSerializedEvent(handler Handler, e *Event, callback ResponceCallback) {
	log.Printf("emitSerializedEvent: %v", e)
	c := cmd{
		command:      emit,
		e:            e,
		handler:      handler,
		respCallback: callback,
	}
	serializerSingleton.eventChan <- c
}
func emitSerializedEventWithResult(handler HandlerWithResult, e *Event, callback ResponceCallback) {
	log.Printf("emitSerializedEventWithResult: %v", e)
	c := cmd{
		command:           emitWithResult,
		e:                 e,
		handlerWithResult: handler,
		respCallback:      callback,
	}
	serializerSingleton.eventChan <- c
}
