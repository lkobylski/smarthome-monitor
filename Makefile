CONTAINER_NAME = smarthome-monitor
CONFIG_FILE = config.yaml

all: build run

build:
	docker-compose build
run:
	docker-compose up -d

logs:
	docker logs -f $(CONTAINER_NAME)

stop:
	docker-compose down

.PHONY: all build run logs stop