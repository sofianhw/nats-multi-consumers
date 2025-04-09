# nats-multi-consumers
Testing nats multi consumers in golang

# Reference
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
