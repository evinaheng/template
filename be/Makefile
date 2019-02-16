run: build start

build:
	@echo ">> Building API..."
	@go build -o api-temp ./cmd/api
	@echo ">> Building Cron..."
	@go build -o cron-temp ./cmd/cron
	@echo ">> Finished"

start-api: server-api.PID

server-api.PID:
	@./api & echo $$! > $@;

stop-api: server-api.PID
	kill `cat $<` && rm $<

start-cron: server-cron.PID

server-cron.PID:
	@./cron & echo $$! > $@;

stop-cron: server-cron.PID
	kill `cat $<` && rm $< 
	
start: start-api start-cron
stop: stop-api stop-cron
