// Copyright 2022‑2023 The NATS Authors
// Licensed under the Apache License, Version 2.0.

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)
	defer cancel()

	nc, err := nats.Connect("nats://127.0.0.1:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Drain()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}

	// Ensure the stream exists (safe if it already does).
	s, _ := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "TEST_STREAM",
		Subjects: []string{"FOO.*"},
	})

	// Create (or update) a durable consumer.
	cons, err := s.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:   "TestConsumerParallelConsume",
		AckPolicy: jetstream.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Launch five concurrent consumers.
	for i := 0; i < 10; i++ {
		cc, err := cons.Consume(
			func(id int) jetstream.MessageHandler {
				return func(msg jetstream.Msg) {
					time.Sleep(3 * time.Second)
					fmt.Printf("%s: %s\n", time.Now().Format("2006-01-02 15:04:05 MST"), msg.Data())
					msg.Ack()
				}
			}(i),
			jetstream.ConsumeErrHandler(func(_ jetstream.ConsumeContext, err error) {
				fmt.Println("consume error:", err)
			}),
		)
		if err != nil {
			log.Fatal(err)
		}
		defer cc.Stop()
	}

	// Wait for Ctrl‑C / SIGTERM.
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}