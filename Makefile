.PHONY: test

build:
	go build -mod=vendor -a -installsuffix cgo -ldflags '-extldflags "-static"' -o fin

run:
	chmod +x fin
	./fin -migrate && ./fin -serve

scrape:
	./fin -scrape

test:
	go test

coverage:
	goveralls -repotoken ${COVERALLS_TOKEN}

docker-build:
	docker build -t jonatasbaldin/fin:latest -t jonatasbaldin/fin:${TRAVIS_TAG} .

docker-push:
	echo "${DOCKER_PASSWORD}" | docker login -u "${DOCKER_USERNAME}" --password-stdin
	docker push jonatasbaldin/fin:latest
	docker push jonatasbaldin/fin:${TRAVIS_TAG}
