package service

import (
	"context"

	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
)

type Service struct {
	*qmgo.Database
}

// http log
var httpLogEntryCollection *qmgo.Collection
var alertGetCollection *qmgo.Collection
var statisticDoneCollection *qmgo.Collection
var serviceStatisticObjectCollection *qmgo.Collection
var uriStatisticObjectCollection *qmgo.Collection
var uriObjectCollection *qmgo.Collection
var svcObjectCollection *qmgo.Collection

// trace
var hopEventCollection *qmgo.Collection
var hopCollection *qmgo.Collection
var operationCollection *qmgo.Collection
var pathEventCollection *qmgo.Collection
var spanCollection *qmgo.Collection
var pathIdCollection *qmgo.Collection
var pathCollection *qmgo.Collection

func NewService(db *qmgo.Database) *Service {
	s := &Service{db}

	s.init()

	httpLogEntryCollection = s.Collection("http_log_entry")
	alertGetCollection = s.Collection("alert_get")
	statisticDoneCollection = s.Collection("statistic_done")
	serviceStatisticObjectCollection = s.Collection("service_statistic_object")
	uriStatisticObjectCollection = s.Collection("uri_statistic_object")
	uriObjectCollection = s.Collection("uri_object")
	svcObjectCollection = s.Collection("service_object")

	hopEventCollection = s.Collection("hop_event")
	hopCollection = s.Collection("hop")
	operationCollection = s.Collection("operation")
	pathEventCollection = s.Collection("path_event")
	pathCollection = s.Collection("path")
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
