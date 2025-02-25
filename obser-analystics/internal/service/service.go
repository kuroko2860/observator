package service

import (
	"github.com/qiniu/qmgo"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Service struct {
	*qmgo.Database
	driver neo4j.DriverWithContext
}

var spanCollection *qmgo.Collection
var hopCollection *qmgo.Collection
var pathCollection *qmgo.Collection
var operationCollection *qmgo.Collection
var hopEventCollection *qmgo.Collection
var pathEventCollection *qmgo.Collection
var httpLogEntryCollection *qmgo.Collection
var alertGetCollection *qmgo.Collection

var uriObjectCollection *qmgo.Collection
var svcObjectCollection *qmgo.Collection

// var statisticDoneCollection *qmgo.Collection
var serviceStatisticObjectCollection *qmgo.Collection
var uriStatisticObjectCollection *qmgo.Collection

func NewService(db *qmgo.Database, driver neo4j.DriverWithContext) *Service {
	s := &Service{db, driver}

	alertGetCollection = s.Collection("alert_get")
	hopCollection = s.Collection("hop")
	pathCollection = s.Collection("path")
	operationCollection = s.Collection("operation")
	hopEventCollection = s.Collection("hop_event")
	pathEventCollection = s.Collection("path_event")
	httpLogEntryCollection = s.Collection("http_log_entry")
	spanCollection = s.Collection("span")
	uriObjectCollection = s.Collection("uri_object")
	svcObjectCollection = s.Collection("svc_object")
	// statisticDoneCollection = s.Collection("statistic_done")
	serviceStatisticObjectCollection = s.Collection("svc_statistic_object")
	uriStatisticObjectCollection = s.Collection("uri_statistic_object")

	return s
}
