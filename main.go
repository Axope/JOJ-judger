package main

import (
	"context"

	"github.com/Axope/JOJ-Judger/common/log"
	"github.com/Axope/JOJ-Judger/configs"
	"github.com/Axope/JOJ-Judger/internal/dao"
	"github.com/Axope/JOJ-Judger/internal/judger"
	"github.com/Axope/JOJ-Judger/internal/middleware/rabbitmq"
	"github.com/Axope/JOJ-Judger/internal/model"
	pb "github.com/Axope/JOJ-Judger/protocol"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// received and serve
	// for msg := range msgs {
	// 	var req request.JudgeRequest
	// 	// TODO: Unmarshal error, redeliver the message
	// 	if err := json.Unmarshal(msg.Body, &req); err != nil {
	// 		log.Logger.Error("json.Unmarshal", log.Any("err", err))
	// 		break
	// 	}
	// 	// TODO: judge error and result
	// 	status, err := judger.Judge(req)
	// 	if err != nil {
	// 		log.Logger.Error("judger internal error", log.Any("err", err))
	// 		continue
	// 	}

	// 	if err := updateJudgeResult(req.SID, status); err != nil {
	// 		log.Logger.Error("updateJudgeResult error", log.Any("err", err))
	// 	} else {
	// 		log.Logger.Debug("updateJudgeResult success")
	// 	}
	// 	if err = msg.Ack(false); err != nil {
	// 		panic(err)
	// 	}
	// }
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

		if err := updateJudgeResult(judgeReq.Sid, status); err != nil {
			log.Logger.Error("updateJudgeResult error", log.Any("err", err))
		} else {
			log.Logger.Debug("updateJudgeResult success")
		}
		if err = msg.Ack(false); err != nil {
			panic(err)
		}
	}
}

func updateJudgeResult(sid string, status model.StatusSet) error {
	log.Logger.Sugar().Debugf("update# sid = %v, status = %v", sid, status)
	realSID, err := primitive.ObjectIDFromHex(sid)
	if err != nil {
		return err
	}
	filter := bson.D{{Key: "_id", Value: realSID}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "status", Value: status}}}}
	_, err = dao.GetSubmissionColl().UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}
