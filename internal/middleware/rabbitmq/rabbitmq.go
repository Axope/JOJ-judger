package rabbitmq

import (
	"fmt"

	"github.com/Axope/JOJ-Judger/common/log"
	"github.com/Axope/JOJ-Judger/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
	msgs  <-chan amqp.Delivery
)

func InitMQ() (<-chan amqp.Delivery, error) {
	cfg := configs.GetRBTConfig()
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	log.Logger.Debug(url)

	var err error
	if conn, err = amqp.Dial(url); err != nil {
		return nil, err
	}
	if ch, err = conn.Channel(); err != nil {
		return nil, err
	}
	if queue, err = ch.QueueDeclare(
		"judger",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	if err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, err
	}

	if msgs, err = ch.Consume(
		queue.Name,
		"",
		false, // auto ack
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	return msgs, nil
}
