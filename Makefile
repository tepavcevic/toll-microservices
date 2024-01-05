obu:
	@go build -o bin/obu obu/main.go
	@./bin/obu

receiver:
	@go build -o bin/receiver ./obu_data_receiver
	@./bin/receiver

calculator:
	@go build -o bin/calculator ./distance_calculator
	@./bin/calculator

aggregator:
	@go build -o bin/aggregator ./aggregator
	@./bin/aggregator

proto:
	protoc --go_out=. --go_opt=paths=source_relative types/ptypes.proto

.PHONY: obu aggregator