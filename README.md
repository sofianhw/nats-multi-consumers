# nats-multi-consumers
Testing nats multi consumers in golang

# Reference
## Alternative 1
[Pull Consumers in JetStream](https://natsbyexample.com/examples/jetstream/pull-consumer/go)

Run:
```sh
docker compose up -d
```

Go to browser:
http://localhost:31312

Split terminal

Terminal 1
```sh
go run publisher.go
```

Terminal 2
```sh
go run consumer.go
```

Terminal 3
```sh
go run consumer.go
```

## Alternative 2
[Queue Groups](https://github.com/nats-io/nats.go/tree/main?tab=readme-ov-file#queue-groups)

Run:
```sh
docker compose up -d
```

Go to browser:
http://localhost:31312

Split terminal
Terminal 1
```sh
go run publisher_old.go
```

Terminal 2
```sh
go run consumer_pr_old.go
```

Terminal 3
```sh
go run consumer_pr_old.go
```