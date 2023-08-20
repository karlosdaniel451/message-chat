package broker

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

var (
	brokerHost = os.Getenv("BROKER_HOST")
	brokerPort = os.Getenv("BROKER_PORT")
)

var conn *nats.Conn

// Connect to a NATS server and assign the obtained connection to `conn`
func Connect() error {
	url := fmt.Sprintf("%s:%s", brokerHost, brokerPort)

	var err error
	conn, err = nats.Connect(url)

	return err
}
