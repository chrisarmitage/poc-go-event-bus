package main

import (
	"fmt"
	"math/rand"
	"time"

	ebs "github.com/chrisarmitage/poc-go-event-bus/internal/event_bus"
)

func publishTo(eb *ebs.EventBus, topic ebs.TopicName, data string) {
	for {
		eb.Publish(topic, data)
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	}
}

func printDataEvent(ch string, data ebs.DataEvent) {
	fmt.Printf("Channel: %s; Topic: %s; DataEvent: %v\n", ch, data.Topic, data.Data)
}

func main() {
	eb := ebs.NewEventBus()

	ch1 := make(chan ebs.DataEvent)
	ch2 := make(chan ebs.DataEvent)
	ch3 := make(chan ebs.DataEvent)
	eb.Subscribe("topic1", ch1)
	eb.Subscribe("topic2", ch2)
	eb.Subscribe("topic2", ch3)
	go publishTo(eb, "topic1", "Hi topic 1")
	go publishTo(eb, "topic2", "Welcome to topic 2")
	for {
		select {
		case d := <-ch1:
			go printDataEvent("ch1", d)
		case d := <-ch2:
			go printDataEvent("ch2", d)
		case d := <-ch3:
			go printDataEvent("ch3", d)
		}
	}
}
