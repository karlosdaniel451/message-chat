package broker

import (
	"fmt"
	"os"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/apache/pulsar-client-go/pulsar/log"
)

var (
	brokerHost string
	brokerPort string
)

var client pulsar.Client

func GetClient() *pulsar.Client {
	return &client
}

// Connect to a NATS server and assign the obtained connection to `conn`
func Connect() error {
	brokerHost = os.Getenv("BROKER_HOST")
	brokerPort = os.Getenv("BROKER_PORT")

	url := fmt.Sprintf("http://%s:%s", brokerHost, brokerPort)

	var err error
	client, err = pulsar.NewClient(pulsar.ClientOptions{
		URL: url,
		Logger: log.DefaultNopLogger(),
	})

	return err
}
