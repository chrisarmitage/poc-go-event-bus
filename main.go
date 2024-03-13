package main

import (
	"fmt"
	"math/rand"
	"time"

	ebs "github.com/chrisarmitage/poc-go-event-bus/internal/event_bus"
)

// func publishTo(eb *ebs.EventBus, topic ebs.TopicName, data string) {
// 	for {
// 		//eb.Publish(topic, data)
// 		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
// 	}
// }

// func printDataEvent(ch string, data ebs.Event) {
// 	fmt.Printf("Channel: %s; Topic: %s; DataEvent: %v\n", ch, data.Topic, data.Data)
// }

func printSendMailEvent(ch string, data ebs.Event) {
	if data.Type() == "UserRegistered" {
		event := data.(ebs.UserRegisteredEvent)

		fmt.Printf("Event: %s, At: %s, Email: %s\n", event.Type(), event.OccurredAt.Format(time.RFC3339), event.Email)
	}
	// fmt.Printf("Channel: %s; Topic: %s; DataEvent: %v\n", ch, data.Topic, data.Data)
}

func main() {
	eb := ebs.NewEventBus()

	// ch1 := make(chan ebs.Event)
	// ch2 := make(chan ebs.Event)
	// ch3 := make(chan ebs.Event)
	// eb.Subscribe("topic1", ch1)
	// eb.Subscribe("topic2", ch2)
	// eb.Subscribe("topic2", ch3)
	// go publishTo(eb, "topic1", "Hi topic 1")
	// go publishTo(eb, "topic2", "Welcome to topic 2")

	chSendMail := make(chan ebs.Event)
	eb.Subscribe("user_registered", chSendMail)

	go func() {
		for {
			e := ebs.UserRegisteredEvent{
				OccurredAt: time.Now(),
				Email:      fmt.Sprintf("user_%04d@example.com", rand.Intn(10000)),
			}
			eb.Publish("user_registered", e)
			time.Sleep(time.Duration(rand.Intn(1000)*3) * time.Millisecond)
		}
	}()
	for {
		select {
		// case d := <-ch1:
		// 	go printDataEvent("ch1", d)
		// case d := <-ch2:
		// 	go printDataEvent("ch2", d)
		// case d := <-ch3:
		// 	go printDataEvent("ch3", d)
		case d := <-chSendMail:
			go printSendMailEvent("chSendMail", d)
		}
	}
}
