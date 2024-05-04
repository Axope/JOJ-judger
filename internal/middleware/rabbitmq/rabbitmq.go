package rabbitmq

import (
	"context"
	"fmt"

	"github.com/Axope/JOJ-Judger/common/log"
	"github.com/Axope/JOJ-Judger/configs"
	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	judgeRecvConn  *amqp.Connection
	judgeRecvCh    *amqp.Channel
	judgeRecvQueue amqp.Queue
)

func InitRecvQ(cfg configs.RabbitmqConfig) (<-chan amqp.Delivery, error) {
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)
	log.Logger.Debug(url)

	var err error
	if judgeRecvConn, err = amqp.Dial(url); err != nil {
		return nil, err
	}
	if judgeRecvCh, err = judgeRecvConn.Channel(); err != nil {
		return nil, err
	}
	if judgeRecvQueue, err = judgeRecvCh.QueueDeclare(
		"publisher",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}
	if err = judgeRecvCh.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	); err != nil {
		return nil, err
	}

	msgs, err := judgeRecvCh.Consume(
		judgeRecvQueue.Name,
		"",
		false, // auto ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return msgs, nil
}

var (
	judgeResultSendConn  *amqp.Connection
	judgeResultSendCh    *amqp.Channel
	judgeResultSendQueue amqp.Queue
)

func InitSendQ(cfg configs.RabbitmqConfig) error {
	username := cfg.Username
	password := cfg.Password
	host := cfg.Host
	port := cfg.Port

	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", username, password, host, port)

	var err error
	if judgeResultSendConn, err = amqp.Dial(url); err != nil {
		return err
	}
	if judgeResultSendCh, err = judgeResultSendConn.Channel(); err != nil {
		return err
	}
	if judgeResultSendQueue, err = judgeResultSendCh.QueueDeclare(
		"JudgeResponseQueue",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}
	return nil
}
func SendMsgByProtobuf(msg []byte) error {
	if err := judgeResultSendCh.PublishWithContext(
		context.TODO(),
		"",
		judgeResultSendQueue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/x-protobuf",
			Body:         msg,
		}); err != nil {
		return err
	}
	return nil
}

func InitMQ() (<-chan amqp.Delivery, error) {
	cfg := configs.GetRBTConfig()
	if err := InitSendQ(cfg); err != nil {
		return nil, err
	}
	return InitRecvQ(cfg)
}
