package service

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"kuroko.com/processor/internal/types"
)

func (s *Service) UpdateDataStatistic(ctx context.Context) error {
	yesterday := time.Now().AddDate(0, 0, -1).Local().Format("20060102")

	sd := []types.StatisticDone{}
	statisticDoneCollection.Find(ctx, bson.M{"_id": yesterday}).All(&sd)
	if len(sd) > 0 {
		return nil
	}

	entries := []types.HttpLogEntry{}
	httpLogEntryCollection.Find(ctx, bson.M{"start_time_date": yesterday}).Sort("start_time").All(&entries)

	svcList := []types.ServiceObject{}
	svcObjectCollection.Find(ctx, bson.M{}).All(&svcList)
	svcStatistic := []types.ServiceStatisticObject{}
	for _, so := range svcList {
		svcStatistic = append(svcStatistic, types.ServiceStatisticObject{
			ID:          so.ServiceName + "*" + yesterday,
			Date:        yesterday,
			ServiceName: so.ServiceName,
			Statistic:   map[int]int64{}})
	}
	uriList := []types.URIObject{}
	uriObjectCollection.Find(ctx, bson.M{}).All(&uriList)
	uriStatistic := []types.URIStatisticObject{}
	for _, u := range uriList {
		if u.Method == "GET" {
			continue
		}
		uriStatistic = append(uriStatistic, types.URIStatisticObject{
			ID:          u.ServiceName + "*" + u.URIPath + "*" + u.Method + "*" + yesterday,
			Date:        yesterday,
			ServiceName: u.ServiceName,
			URIPath:     u.URIPath,
			Method:      u.Method,
			Statistic:   map[int]int64{}})
	}
	entriesAlert := []types.HttpLogEntry{}
	for _, hle := range entries {
		if hle.Method == "GET" {
			entriesAlert = append(entriesAlert, hle)
			continue
		}

		// !TODO
		hour := time.Unix(hle.StartTime/1000, 0).Hour()
		serviceName := hle.ServiceName
		uriPath := hle.URIPath
		method := hle.Method
		date := hle.StartTimeDate

		for _, sso := range svcStatistic {
			if sso.ID == serviceName+"*"+date {
				sso.Statistic[hour] += 1
			}
		}

		for _, uso := range uriStatistic {
			if uso.ID == serviceName+"*"+uriPath+"*"+method+"*"+date {
				uso.Statistic[hour] += 1
			}
		}
	}
	go s.UpdateDataAlertGet(ctx, entriesAlert)

	_, err := serviceStatisticObjectCollection.InsertMany(ctx, svcStatistic)
	if err != nil {
		return err
	}
	_, err = uriStatisticObjectCollection.InsertMany(ctx, uriStatistic)
	if err != nil {
		return err
	}
	_, err = statisticDoneCollection.UpsertId(ctx, yesterday, types.StatisticDone{Date: yesterday})
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateDataAlertGet(ctx context.Context, http_logs []types.HttpLogEntry) error {
	mapTime := make(map[string]int)
	for _, hle := range http_logs {
		key := hle.ServiceName + "*" + hle.URIPath + "*" + hle.UserId + "*" + hle.Referer
		time := mapTime[key]
		if time == 0 {
			mapTime[key] = int(hle.StartTime)
		} else {
			if hle.StartTime/1000-int64(time) < 30 { // goi cung 1 api trong 30s
				id := hle.ServiceName + "*" + hle.URIPath + "*" + hle.Referer
				_, err := alertGetCollection.UpsertId(ctx, id, types.AlertGetObject{
					ID:          id,
					URIPath:     hle.URIPath,
					Referer:     hle.Referer,
					ServiceName: hle.ServiceName,
					Entry:       hle,
				})
				if err != nil {
					return err
				}
			}
			mapTime[key] = int(hle.StartTime)
		}
	}
	return nil
}
