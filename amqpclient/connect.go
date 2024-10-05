package amqpclient

import (
	"fmt"
	"light-backend/config"
	"net"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn *amqp.Connection
	mtx  sync.Mutex
)

func Init() error {
	c, err := Connect()
	conn = c
	return err
}

func Connect() (*amqp.Connection, error) {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s", config.Config("AMQP_USER"), config.Config("AMQP_PASSWD"),
		config.Config("AMQP_HOST"), config.Config("AMQP_PORT"))

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	c, err := amqp.DialConfig(uri, amqp.Config{Dial: dialer.Dial})

	if err != nil {
		return nil, err
	}
	return c, nil
}

// Before we open the channel heartbeats do not work
// An use of a channel is not an option eather, because then not thread-safe
// So to ensure that connection wasnt interrupted we try to reconnect
func getChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	var ch *amqp.Channel
	ch, err := conn.Channel()
	if err != nil {
		mtx.Lock()
		conn, err = Connect()
		mtx.Unlock()
		if err != nil {
			return nil, err
		}
		ch, err = conn.Channel()
		if err != nil {
			return nil, err
		}
	}
	return ch, nil
}

func NewChannel() (*amqp.Channel, error) {
	ch, err := getChannel(conn)
	return ch, err
}
