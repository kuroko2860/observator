package service

import (
	"github.com/qiniu/qmgo"
)

type Service struct {
	*qmgo.Database
}

var httpLogEntryCollection *qmgo.Collection
var alertGetCollection *qmgo.Collection
var statisticDoneCollection *qmgo.Collection
var serviceStatisticObjectCollection *qmgo.Collection
var uriStatisticObjectCollection *qmgo.Collection
var uriObjectCollection *qmgo.Collection
var svcObjectCollection *qmgo.Collection

func NewService(db *qmgo.Database) *Service {
	s := &Service{db}

	httpLogEntryCollection = s.Collection("http_log_entry")
	alertGetCollection = s.Collection("alert_get")
	statisticDoneCollection = s.Collection("statistic_done")
	serviceStatisticObjectCollection = s.Collection("service_statistic_object")
	uriStatisticObjectCollection = s.Collection("uri_statistic_object")
	uriObjectCollection = s.Collection("uri_object")
	svcObjectCollection = s.Collection("service_object")
	return s
}
