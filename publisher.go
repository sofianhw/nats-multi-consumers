// Copyright 2022‑2023 The NATS Authors
// Licensed under the Apache License, Version 2.0.

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
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

	// Make sure the stream exists (idempotent if it already does).
	_, _ = js.CreateStream(ctx, jetstream.StreamConfig{
		Name:     "TEST_STREAM",
		Subjects: []string{"FOO.*"},
	})

	endlessPublish(ctx, js, nc)
}

func endlessPublish(ctx context.Context, js jetstream.JetStream, nc *nats.Conn) {
	for i := 0; ; i++ {
	// for i := 0; i < 5; i++ {
		time.Sleep(500 * time.Millisecond)

		if nc.Status() != nats.CONNECTED {
			continue
		}
		_, err := js.Publish(ctx, "FOO.TEST1", []byte(fmt.Sprintf("msg %d", i)))
		if err != nil {
			fmt.Println("publish error:", err)
		}
	}
}