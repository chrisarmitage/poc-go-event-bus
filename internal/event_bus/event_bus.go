package eventbus

import (
	"sync"
	"time"
)

type DataChannel chan Event

type DataChannels []DataChannel

type TopicName string

type EventBus struct {
	subscribers map[TopicName]DataChannels
	rm          sync.RWMutex
}

type Event interface {
	Type() string
}

type UserRegisteredEvent struct {
	OccurredAt time.Time
	Email      string
}

func (u UserRegisteredEvent) Type() string {
	return "UserRegistered"
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: map[TopicName]DataChannels{},
	}
}

func (eb *EventBus) Subscribe(topic TopicName, ch DataChannel) {
	eb.rm.Lock()
	defer eb.rm.Unlock()

	if prev, found := eb.subscribers[topic]; found {
		eb.subscribers[topic] = append(prev, ch)
	} else {
		eb.subscribers[topic] = append([]DataChannel{}, ch)
	}
}

func (eb *EventBus) Publish(topic TopicName, event Event) {
	eb.rm.RLock()
	defer eb.rm.RUnlock()

	if chans, found := eb.subscribers[topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		channels := append(DataChannels{}, chans...)
		go func(event Event, dataChannelSlices DataChannels) {
			for _, ch := range dataChannelSlices {
				ch <- event
			}
		}(event, channels)
	}
}
