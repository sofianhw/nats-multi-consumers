package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName  = "TEST_STREAM"
	subject     = "FOO.TEST1"
	durableName = "JOBS" // All consumer processes bind to the same durable.
	batch       = 100
	maxWait     = 2 * time.Second
)

func main() {
	ctx := context.Background()

	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure stream and durable consumer exist. (This is idempotent.)
	strm, _ := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: []string{"FOO.*"},
	})

	cons, err := strm.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       durableName,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       30 * time.Second,
		MaxAckPending: 10_000,
	})
	if err != nil {
		log.Fatal(err)
	}

	for {
		msgs, err := cons.Fetch(2, jetstream.FetchMaxWait(1*time.Second))
		if err != nil {
			fmt.Println(err)
		}
		for msg := range msgs.Messages() {
			time.Sleep(3 * time.Second)
			fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05 MST"), msg.Data())
			msg.Ack()
		}
		if msgs.Error() != nil {
			fmt.Println("Error fetching messages: ", err)
		}
	}
}