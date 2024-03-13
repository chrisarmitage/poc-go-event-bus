package eventbus

import "sync"

type DataEvent struct {
	Data  any
	Topic TopicName
}

type DataChannel chan DataEvent

type DataChannels []DataChannel

type TopicName string

type EventBus struct {
	subscribers map[TopicName]DataChannels
	rm          sync.RWMutex
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

func (eb *EventBus) Publish(topic TopicName, data any) {
	eb.rm.RLock()
	defer eb.rm.RUnlock()

	if chans, found := eb.subscribers[topic]; found {
		// this is done because the slices refer to same array even though they are passed by value
		// thus we are creating a new slice with our elements thus preserve locking correctly.
		channels := append(DataChannels{}, chans...)
		go func(data DataEvent, dataChannelSlices DataChannels) {
			for _, ch := range dataChannelSlices {
				ch <- data
			}
		}(DataEvent{Data: data, Topic: topic}, channels)
	}
}
