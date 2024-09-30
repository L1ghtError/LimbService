package amqpclient

import (
	"fmt"
	"light-backend/config"
	"net"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var Conn *amqp.Connection

func Connect() error {
	uri := fmt.Sprintf("amqp://%s:%s@%s:%s", config.Config("AMQP_USER"), config.Config("AMQP_PASSWD"),
		config.Config("AMQP_HOST"), config.Config("AMQP_PORT"))

	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	c, err := amqp.DialConfig(uri, amqp.Config{Dial: dialer.Dial})
	if err != nil {
		return err
	}
	Conn = c
	return nil
}
