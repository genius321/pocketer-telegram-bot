.PHONY:

build: 
	docker build -t pocketer-telegram-bot:v0.1 .

run:
	docker run --name pocketer-bot --env-file .env pocketer-telegram-bot:v0.1

clean:
	docker stop pocketer-bot
	docker rm pocketer-bot

all: build run