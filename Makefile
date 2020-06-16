run:
	
build:
	env GOOS=linux  go build -v spiderweb

docker:
	docker build . -t spiderweb:0.23
	docker tag spiderweb:0.23 garybowers/spiderweb:0.23
	docker push garybowers/spiderweb:0.23
	
deploy:
	kubectl apply -f .

test:
