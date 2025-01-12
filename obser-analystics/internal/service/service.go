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

func NewService(db *qmgo.Database, driver neo4j.DriverWithContext) *Service {
	s := &Service{db, driver}

	timeCollection = s.Collection("time")
	spanCollection = s.Collection("span")
	return s
}
