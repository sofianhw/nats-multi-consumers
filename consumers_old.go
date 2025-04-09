package main

import (
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

    // Each process uses a shared queue group ("workers")
    _, err = js.QueueSubscribe("image.task", "workers", func(m *nats.Msg) {
        task := string(m.Data)
        start := time.Now()
        log.Printf("[%s] Received: %s", start.Format("15:04:05"), task)

        // Simulate image generation
        time.Sleep(10 * time.Second)

        end := time.Now()
        log.Printf("[%s] Completed: %s (Duration: %.1fs)", end.Format("15:04:05"), task, end.Sub(start).Seconds())

        m.Ack()
    }, nats.Durable("worker"), nats.ManualAck(), nats.AckWait(30*time.Second))

    if err != nil {
        log.Fatal(err)
    }

    log.Println("Consumer started. Waiting for tasks...")
    select {}
}