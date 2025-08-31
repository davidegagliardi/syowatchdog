.PHONY: build clean docker-build docker-run docker-run-prod docker-stop docker-logs

CONTAINER_NAME := syowatchdog

build:
	go build -o bin/syowatchdog .

clean:
	rm -f bin/syowatchdog

docker-build:
	docker build -t syowatchdog:latest .

docker-run:
	docker run -d --rm --name $(CONTAINER_NAME) syowatchdog:latest

docker-run-prod:
	docker run -d --name $(CONTAINER_NAME) \
		--env-file $(PWD)/.env \
		-v $(PWD)/data:/app/data \
		-v $(PWD)/configs:/app/configs \
		syowatchdog:latest

docker-stop:
	docker stop $(CONTAINER_NAME)

docker-logs:
	docker logs -f $(CONTAINER_NAME)

docker-restart: docker-stop docker-run

# Clean up everything
docker-clean:
	-docker stop $(CONTAINER_NAME)
	-docker rm $(CONTAINER_NAME)
	docker rmi syowatchdog:latest