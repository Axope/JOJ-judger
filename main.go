package main

import (
	"github.com/Axope/JOJ-Judger/common/log"
	"github.com/Axope/JOJ-Judger/configs"
	"github.com/Axope/JOJ-Judger/internal/dao"
	"github.com/Axope/JOJ-Judger/internal/judger"
	"github.com/Axope/JOJ-Judger/internal/middleware/rabbitmq"
	pb "github.com/Axope/JOJ-Judger/protocol"
	"google.golang.org/protobuf/proto"
)

func main() {
	// load configs
	configs.InitConfigs()

	// log
	log.InitLogger()
	defer log.Logger.Sync()
	log.Logger.Info("Golbal log init success")

	// InitJudger()

	// mongoDB
	if err := dao.InitMongo(); err != nil {
		log.Logger.Error("mongoDB init failed", log.Any("err", err))
		return
	} else {
		log.Logger.Info("mongoDB init success")
	}
	// rabbitMQ
	msgs, err := rabbitmq.InitMQ()
	if err != nil {
		log.Logger.Error("RabbitMQ init failed", log.Any("err", err))
		return
	} else {
		log.Logger.Info("RabbitMQ init success")
	}

	for msg := range msgs {
		judgeReq := &pb.Judge{}
		// TODO: Unmarshal error, redeliver the message
		if err := proto.Unmarshal(msg.Body, judgeReq); err != nil {
			log.Logger.Error("json.Unmarshal", log.Any("err", err))
			break
		}
		// TODO: judge error and result
		status, err := judger.JudgeFromProtobuf(judgeReq)
		if err != nil {
			log.Logger.Error("judger internal error", log.Any("err", err))
			continue
		}

		if err := handleJudgeResult(judgeReq, status); err != nil {
			log.Logger.Error("handleJudgeResult error", log.Any("err", err))
		} else {
			log.Logger.Debug("handleJudgeResult success")
		}
		if err = msg.Ack(false); err != nil {
			panic(err)
		}
	}
}

func handleJudgeResult(judgeReq *pb.Judge, status pb.StatusSet) error {
	log.Logger.Sugar().Debugf("update# judgeReq = %v, status = %v", judgeReq, status)
	judgeResultReq := &pb.JudgeResult{
		Sid:             judgeReq.Sid,
		Pid:             judgeReq.Pid,
		Uid:             judgeReq.Uid,
		Cid:             judgeReq.Cid,
		Status:          status,
		SubmitTimestamp: judgeReq.SubmitTimestamp,
	}
	msg, err := proto.Marshal(judgeResultReq)
	if err != nil {
		return err
	}
	if err := rabbitmq.SendMsgByProtobuf(msg); err != nil {
		return err
	}
	log.Logger.Info("send judgeResult request success",
		log.Any("judgeResultReq", judgeResultReq))
	return nil
}
