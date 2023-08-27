package broker

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

var (
	brokerHost string
	brokerPort string
)

var conn *nats.Conn

func GetConn() *nats.Conn {
	return conn
}

// Connect to a NATS server and assign the obtained connection to `conn`
func Connect() error {
	brokerHost = os.Getenv("BROKER_HOST")
	brokerPort = os.Getenv("BROKER_PORT")

	url := fmt.Sprintf("nats://%s:%s", brokerHost, brokerPort)

	var err error
	conn, err = nats.Connect(url)

	return err
}
