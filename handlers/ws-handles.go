package handlers

import (
	"context"
	"fmt"
	"light-backend/amqpclient"
	"light-backend/model"
	"light-backend/validation"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
)

// TODO: allow to enhance only for onwer
// Also do not allow operation on non existing files
// Deprecated - better use SSE
func EnhanceImageWs(c *websocket.Conn) {

	myValidator := validation.XValidator{Validator: validator.New()}
	body := new(model.MUpscaleImage)
	if err := c.ReadJSON(body); err != nil {
		fmt.Printf("Websock got err:%s\n", err.Error())
		return
	}
	if errs := myValidator.Validate(body); len(errs) > 0 && errs[0].Error {
		err := validation.GenerateErrorResp(&errs)
		fmt.Printf("Websock validation err:%s\n", err.Error())
		return
	}

	// AMQP
	ch, err := amqpclient.NewChannel()
	if err != nil {
		fmt.Printf("Websock AMQP channel err:%s\n", err.Error())
		return
	}
	defer ch.Close()
	queueName := "UpscaleImage"
	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		fmt.Printf("Websock AMQP channel err:%s\n", err.Error())
		return
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
		fmt.Printf("Websock AMQP consume err:%s\n", err.Error())
		return
	}
	timeN := time.Now()
	ctx, _ := context.WithDeadline(context.Background(), timeN.Add(30*time.Second))
	// Convert Websock input to Big-endian input for AMQP
	raw := body.Htoberaw()

	corrId := randomString(32)
	err = ch.PublishWithContext(ctx,
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType:   "application/octet-stream",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          raw,
		})
	if err != nil {
		fmt.Printf("Websock amqp write err:%s\n", err.Error())
		return
	}

	// Communication Loop
	for msg := range msgs {
		if corrId == msg.CorrelationId {
			body := string(msg.Body)
			if body == "end" {
				c.Close()
				break
			}
			parts := strings.Split(body, ":")
			if len(parts) != 2 {
				fmt.Printf("Websock AMQP received mailformed\n")
				break
			}
			_ = parts[0] // Progress
			workerEstimation := parts[1]

			duration := time.Since(timeN) * 2
			weDuration, err := time.ParseDuration(workerEstimation)
			if err != nil {
				fmt.Printf("Websock amqp write err:%s\n", err.Error())
				return
			}
			total := duration + weDuration
			totalMs := strconv.FormatInt(total.Milliseconds(), 10) + "ms"
			resp := fiber.Map{"estimation": totalMs}
			if err := c.WriteJSON(resp); err != nil {
				fmt.Printf("Websock write err:%s\n", err.Error())
				break
			}
		}
	}

}
