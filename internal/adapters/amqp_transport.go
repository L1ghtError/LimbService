package adapters

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"light-backend/internal/amqpclient"
	"light-backend/internal/domain/media"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	amqp "github.com/rabbitmq/amqp091-go"
)

const processImageQueue = "ProcessImage"
const getAppInfoQueue = "GetAppInfo"

type AmqpTransport struct {
	collection mongo.Collection
}

type AppInfoTask = media.AppInfoTask
type ImageTask = media.ImageTask
type ImageTaskResult = media.ImageTaskResult

func NewAmqpTransport() *AmqpTransport {
	return &AmqpTransport{}
}

func (t *AmqpTransport) ProcessImage(task ImageTask, w *bufio.Writer) error {

	ch, err := amqpclient.NewChannel()
	if err != nil {
		fmt.Printf("[AmqpTransport] channel err:%s\n", err.Error())
		return nil
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Printf("[AmqpTransport] QueueDeclare err:%s\n", err.Error())
		return nil
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("[AmqpTransport] Consume err:%s\n", err.Error())
		return nil
	}

	timeN := time.Now()
	ctx, cancel := context.WithDeadline(context.Background(), timeN.Add(30*time.Second))
	defer cancel()

	paylaod, err := json.Marshal(task)

	err = ch.PublishWithContext(ctx,
		"",                // exchange
		processImageQueue, // routing key
		false,             // mandatory
		false,             // immediate
		amqp.Publishing{
			ContentType: "application/octet-stream",
			ReplyTo:     q.Name,
			Body:        paylaod,
		})
	if err != nil {
		fmt.Printf("[AmqpTransport] PublishWithContext err:%s\n", err.Error())
		return nil
	}

	for msg := range msgs {
		var resp ImageTaskResult
		err = json.Unmarshal(msg.Body, &resp)
		if err != nil {
			fmt.Printf("[AmqpTransport] Unmarshal err:%s\n", err.Error())
			return nil
		}

		if resp.Status == media.StatusDone {
			break
		} else if resp.Status == media.StatusFail {
			break
		} else if resp.Status == media.StatusProgress {
			progress := resp.Message

			fmt.Fprintf(w, "{\"progress\":\"%s\"}", progress)
			if err := w.Flush(); err != nil {
				fmt.Printf("[AmqpTransport] StatusProgress err:%s\n", err.Error())
				break
			}
		}

	}
	return nil
}

func (t *AmqpTransport) GetWorkersInfo() (AppInfoTask, error) {
	ch, err := amqpclient.NewChannel()
	if err != nil {
		fmt.Printf("[AmqpTransport] channel err:%s\n", err.Error())
		return AppInfoTask{}, err
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Printf("[AmqpTransport] QueueDeclare err:%s\n", err.Error())
		return AppInfoTask{}, err
	}
	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Printf("[AmqpTransport] Consume err:%s\n", err.Error())
		return AppInfoTask{}, err
	}

	timeN := time.Now()
	ctx, cancel := context.WithDeadline(context.Background(), timeN.Add(30*time.Second))
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",              // exchange
		getAppInfoQueue, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ReplyTo: q.Name,
		})

	if err != nil {
		fmt.Printf("[AmqpTransport] PublishWithContext err:%s\n", err.Error())
		return AppInfoTask{}, err
	}
	var resp AppInfoTask

	msg, ok := <-msgs
	if !ok {
		return resp, errors.New("Internal error")
	}

	err = json.Unmarshal(msg.Body, &resp)
	if err != nil {
		fmt.Printf("[AmqpTransport] Unmarshal err:%s\n", err.Error())
		return AppInfoTask{}, nil
	}

	return resp, nil
}
