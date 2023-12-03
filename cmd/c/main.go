package main

import (
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func main() {
	var opts server.Options
	if err := opts.ProcessConfigFile("./etc/n3-c1.conf"); err != nil {
		panic(err)
	}

	server, err := server.NewServer(&opts)
	if err != nil {
		panic(err)
	}

	server.ConfigureLogger()
	server.Start()
	if !server.ReadyForConnections(time.Second * 10) {
		panic("NATS server didn't start")
	}

	client, err := nats.Connect("", nats.InProcessServer(server))
	if err != nil {
		panic(err)
	}

	js, err := client.JetStream()
	if err != nil {
		panic(err)
	}

	kv, err := js.CreateKeyValue(&nats.KeyValueConfig{
		Bucket: "test",
	})
	if err != nil {
		panic(err)
	}

	for k, v := range map[string]string{
		"foo": "bar",
		"baz": "qux",
	} {
		if _, err := kv.PutString(k, v); err != nil {
			panic(err)
		}
	}

	keys, err := kv.Keys()
	if err != nil {
		panic(err)
	}

	fmt.Println("Found keys:", keys)

	server.WaitForShutdown()
}
