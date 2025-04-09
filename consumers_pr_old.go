package main

import (
    "log"
    "runtime"
    "time"

    "github.com/nats-io/nats.go"
)

const (
    numWorkers = 5 // number of goroutines per consumer process
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

    log.Printf("Starting consumer with %d workers (goroutines)\n", numWorkers)

    for i := 0; i < numWorkers; i++ {
        workerID := i + 1

        go func(id int) {
            _, err := js.QueueSubscribe("image.task", "workers", func(m *nats.Msg) {
                task := string(m.Data)
                start := time.Now()
                log.Printf("[Worker-%d] [%s] Received: %s", id, start.Format("15:04:05"), task)

                // Simulate image generation time
                time.Sleep(10 * time.Second)

                end := time.Now()
                log.Printf("[Worker-%d] [%s] Completed: %s (Duration: %.1fs)", id, end.Format("15:04:05"), task, end.Sub(start).Seconds())

                m.Ack()
            }, nats.Durable("worker"), nats.ManualAck(), nats.AckWait(30*time.Second))
            if err != nil {
                log.Fatalf("Worker-%d failed to subscribe: %v", id, err)
            }

            runtime.Goexit() // keep goroutine alive after subscription
        }(workerID)
    }

    select {} // block main forever
}