run-mock-zipkin:
	nodemon ./mock-zipkin/server.js

run-mock-nats:
	go run ./mock-nats/main.go

run-obser-analystics:
	cd ./obser-analystics/main
	go run main.go

run-obser-http-log:
	cd ./obser-http-log/main
	go run main.go

run-obser-trace:
	cd ./obser-trace/main
	go run main.go

run-frontend:
	cd ./frontend
	npm run dev

run-neo4j:
	docker run -d --name neo4j -p 7687:7687 -p 7474:7474 -v D:\Projects\kltn\neo4j:/data neo4j
run-nats:
	docker run -d --name nats -p 4222:4222 nats

stop-neo4j:
	docker stop neo4j

stop-nats:
	docker stop nats

start-neo4j:
	docker restart neo4j

start-nats:
	docker restart nats
run-mock-data:
	make run-mock-nats
	make run-mock-zipkin
	make start-neo4j
	make start-nats