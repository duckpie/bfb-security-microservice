ENV=local
PROTOV=0.3.2

.PHONY: run
run:
	sudo docker-compose -f docker-compose.$(ENV).yml build
	sudo docker-compose -f docker-compose.$(ENV).yml up


.PHONY: build
build:
	sudo docker-compose -f docker-compose.$(ENV).yml build


.PHONY: test
test:
	sudo docker-compose -f docker-compose.test.yml build
	sudo docker-compose -f docker-compose.test.yml up --remove-orphans


.PHONY: protoup
protoup:
	go get -u github.com/wrs-news/golang-proto@v$(PROTOV)


.PHONY: count
count:
	find . -name tests -prune -o -type f -name '*.go' | xargs wc -l


.DEFAULT_GOAL := run