package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/nats-io/nats.go"
)

type HttpRequestLog struct {
	ServiceName string `json:"service_name" bson:"service_name"`
	Method      string `json:"method" bson:"method"`
	URI         string `json:"uri" bson:"uri"`
	URIPath     string `json:"uri_path" bson:"uri_path"`
	Protocol    string `json:"protocol" bson:"protocol"`
	Host        string `json:"host" bson:"host"`
	RemoteIP    string `json:"remote_ip" bson:"remote_ip"`
	RequestId   string `json:"request_id" bson:"request_id"`
	// TraceId      string `json:"trace_id" bson:"trace_id"`
	// PathId       uint32 `json:"path_id" bson:"path_id"`
	UserId       string `json:"user_id" bson:"user_id"`
	Referer      string `json:"referer" bson:"referer"`
	UserAgent    string `json:"user_agent" bson:"user_agent"`
	StartTime    int64  `json:"start_time" bson:"start_time"`
	Duration     int64  `json:"duration" bson:"duration"`
	ResquestSize string `json:"resquest_size" bson:"resquest_size"`
	ResponseSize int64  `json:"response_size" bson:"response_size"`
	StatusCode   int    `json:"status_code" bson:"status_code"`
	ErrorMessage string `json:"error_message" bson:"error_message"`
}

var services = []string{}
var SERVICES_NUMS = 7

// var NATS_SERVER = "nats://nats-server:4222"
var NATS_SERVER = "nats://localhost:4222"

func main() {
	for i := 1; i <= SERVICES_NUMS; i++ {
		services = append(services, fmt.Sprintf("service-%d", i))
	}
	nc, err := nats.Connect(NATS_SERVER)
	if err != nil {
		fmt.Println("Failed to connect to NATS", err)
		return
	}
	fmt.Println("Connected to NATS")
	for {
		httpLog := generateHttpLog()
		httpLogJson, err := json.Marshal(httpLog)
		if err != nil {
			fmt.Printf("Error when json httplog: %v\n", err)
			continue
		}
		nc.Publish("http-log", httpLogJson)
		fmt.Println("Published a message")
		time.Sleep(1 * time.Millisecond)
	}
}

func generateHttpLog() *HttpRequestLog {

	methods := []string{"GET", "POST", "PUT", "DELETE", "PATCH"}
	protocols := []string{"HTTP/1.1", "HTTP/2"}
	statusCodes := []int{200, 201, 400, 401, 403, 404, 500}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.107 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_0 like Mac OS X) AppleWebKit/537.36 (KHTML, like Gecko) Version/14.0 Mobile/15E148 Safari/537.36",
	}
	randomDuration := rand.Int63n(5000) + 50 // Between 50ms and 5s
	startTime := time.Now().UnixNano() / int64(time.Millisecond)

	return &HttpRequestLog{
		ServiceName: services[rand.Intn(SERVICES_NUMS)],
		Method:      methods[rand.Intn(len(methods))],
		URI:         "/api/v1/" + faker.Word(),
		URIPath:     "/api/v1/" + faker.Word(),
		Protocol:    protocols[rand.Intn(len(protocols))],
		Host:        faker.DomainName(),
		RemoteIP:    faker.IPv4(),
		RequestId:   faker.UUIDHyphenated(),
		// TraceId:      faker.UUIDHyphenated(),
		// PathId:       uint32(rand.Intn(1000)),
		UserId:    faker.UUIDDigit(),
		Referer:   "https://" + faker.DomainName() + "/ref",
		UserAgent: userAgents[rand.Intn(len(userAgents))],

		StartTime:    startTime,
		Duration:     randomDuration,
		ResquestSize: strconv.Itoa(rand.Intn(5000) + 100),
		ResponseSize: rand.Int63n(5000) + 100,
		StatusCode:   statusCodes[rand.Intn(len(statusCodes))],
		ErrorMessage: "",
	}
}
