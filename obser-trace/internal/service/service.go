package service

import (
	"context"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Service struct {
	*qmgo.Database
	driver neo4j.DriverWithContext
}

var hopEventCollection *qmgo.Collection
var pathEventCollection *qmgo.Collection
var spanCollection *qmgo.Collection
var timeCollection *qmgo.Collection // time in milliseconds  in db: 1739721982
var pathIdCollection *qmgo.Collection

func NewService(db *qmgo.Database, driver neo4j.DriverWithContext) *Service {
	s := &Service{db, driver}

	hopEventCollection = s.Collection("hop_event")
	pathEventCollection = s.Collection("path_event")
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
