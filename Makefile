.PHONY: build clean docker-build docker-run docker-run-prod docker-stop docker-logs

CONTAINER_NAME := syowatchdog

##@ Build

build: ## Build application binary.
	go build -o bin/syowatchdog .

clean: ## Remove application binary.
	rm -f bin/syowatchdog

##@ Image build

docker-build: ## Build application Docker image.
	docker build -t syowatchdog:latest .

docker-run: ## Run application Docker container.
	docker run -d --rm --name $(CONTAINER_NAME) syowatchdog:latest

docker-run-prod: ## Run application Docker container, with sensible configuration.
	docker run -d --name $(CONTAINER_NAME) \
		--env-file $(PWD)/.env \
		-v $(PWD)/data:/app/data \
		-v $(PWD)/configs:/app/configs \
		syowatchdog:latest

docker-stop: ## Stop application Docker container.
	docker stop $(CONTAINER_NAME)

docker-logs: ## Print application Docker container logs.
	docker logs -f $(CONTAINER_NAME)

docker-restart: docker-stop docker-run ## Restart application Docker container.

docker-clean: ## Clean up everything.
	-docker stop $(CONTAINER_NAME)
	-docker rm $(CONTAINER_NAME)
	docker rmi syowatchdog:latest

##@ Help

.DEFAULT_GOAL := help
.PHONY: help
help: ## Show this help screen.
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-25s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
