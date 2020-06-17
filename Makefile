run:

all: build docker deploy

build:
	env GOOS=linux  go build -v spiderweb

docker:
	docker build . -t spiderweb:0.26
	docker tag spiderweb:0.26 garybowers/spiderweb:0.26
	docker push garybowers/spiderweb:0.26
	
deploy:
	kubectl apply -f .

test:
