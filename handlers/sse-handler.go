package handlers

import (
	"bufio"
	"context"
	"fmt"
	"light-backend/amqpclient"
	"light-backend/model"
	"light-backend/validation"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/valyala/fasthttp"
)

// TODO: Need more advanced tehnique to estimate process time
func EnhanceImage(c *fiber.Ctx) error {
	// TODO: move header definitions to middleware

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	myValidator := validation.XValidator{Validator: validator.New()}
	body := new(model.MUpscaleImage)
	if err := c.BodyParser(body); err != nil {
		return err
	}
	if errs := myValidator.Validate(body); len(errs) > 0 && errs[0].Error {
		err := validation.GenerateErrorResp(&errs)
		return err
	}

	c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		// AMQP
		// TODO: move all amqp-handling stuff to approptiate worker-service
		ch, err := amqpclient.NewChannel()
		if err != nil {
			fmt.Printf("SSE AMQP channel err:%s\n", err.Error())
			return
		}
		defer ch.Close()
		queueName := "UpscaleImage"
		q, err := ch.QueueDeclare(
			"",    // name
			false, // durable
			true,  // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		if err != nil {
			fmt.Printf("SSE AMQP DEC QUEUE err:%s\n", err.Error())
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
			fmt.Printf("SSE AMQP CONSUME err:%s\n", err.Error())
			return
		}
		timeN := time.Now()
		ctx, _ := context.WithDeadline(context.Background(), timeN.Add(30*time.Second))
		// Convert input to Big-endian input for AMQP
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
			fmt.Printf("SSE AMQP PUBLISH err:%s\n", err.Error())
			return
		}
		// Communication Loop
		for msg := range msgs {
			if corrId == msg.CorrelationId {
				body := string(msg.Body)
				if body == "end" {
					fmt.Print("eeeend!!!\n")
					break
				}
				parts := strings.Split(body, ":")
				if len(parts) != 2 {
					fmt.Printf("SSE AMQP received mailformed %v\n", parts)
					break
				}
				fmt.Printf("got %s\n", body)
				_ = parts[0] // Progress
				workerEstimation := parts[1]

				duration := time.Since(timeN) * 2
				weDuration, err := time.ParseDuration(workerEstimation)
				if err != nil {
					fmt.Printf("SSE amqp write err:%s\n", err.Error())
					return
				}
				total := duration + weDuration
				totalMs := strconv.FormatInt(total.Milliseconds(), 10) + "ms"

				fmt.Fprintf(w, "{\"estimation\":\"%s\"}", totalMs)
				fmt.Printf("Sent %v\n", totalMs)
				if err := w.Flush(); err != nil {
					fmt.Printf("SSE write err:%s\n", err.Error())
					break
				}
			}
		}
	}))
	return nil
}
