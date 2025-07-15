start:
	docker stop $$(docker ps -aq) \
	&& docker compose start
gen:
	go generate ./...
test:
	go test ./...