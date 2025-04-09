package main

import (
    "fmt"
    "log"
    "time"

    "github.com/nats-io/nats.go"
)

func main() {
    nc, err := nats.Connect(nats.DefaultURL)
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Drain()

    js, err := nc.JetStream()
    if err != nil {
        log.Fatal(err)
    }

    // Ensure stream exists
    _, err = js.AddStream(&nats.StreamConfig{
        Name:     "IMAGE_TASKS",
        Subjects: []string{"image.task"},
    })
    if err != nil && err != nats.ErrStreamNameAlreadyInUse {
        log.Fatal(err)
    }

    i := 0
    for {
        i++
        msg := fmt.Sprintf("Task-%03d", i)
        _, err := js.Publish("image.task", []byte(msg))
        if err != nil {
            log.Printf("Failed to publish %s: %v", msg, err)
        } else {
            log.Printf("[%s] Published: %s", time.Now().Format("15:04:05"), msg)
        }
        time.Sleep(500 * time.Millisecond)
    }
}