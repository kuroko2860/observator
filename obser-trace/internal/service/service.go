package service

import (
	"context"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	*qmgo.Database
}

var hopEventCollection *qmgo.Collection
var hopCollection *qmgo.Collection
var operationCollection *qmgo.Collection
var pathEventCollection *qmgo.Collection
var spanCollection *qmgo.Collection
var timeCollection *qmgo.Collection // time in milliseconds  in db: 1739721982
var pathIdCollection *qmgo.Collection
var pathCollection *qmgo.Collection

func NewService(db *qmgo.Database) *Service {
	s := &Service{db}

	hopEventCollection = s.Collection("hop_event")
	hopCollection = s.Collection("hop")
	operationCollection = s.Collection("operation")
	pathEventCollection = s.Collection("path_event")
	pathCollection = s.Collection("path")
	timeCollection = s.Collection("time")
	spanCollection = s.Collection("span")
	pathIdCollection = s.Collection("path_id")

	pathIds = make(map[uint32]bool)
	var pathIdArr []struct {
		ID uint32 `json:"_id" bson:"_id"`
	}
	pathIdCollection.Find(context.Background(), bson.M{}).All(&pathIdArr)
	for _, pathId := range pathIdArr {
		pathIds[pathId.ID] = true
	}
	return s
}
