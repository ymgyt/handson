docker-build:
	docker build -f Dockerfile_app -t go-app .
	docker build -f Dockerfile_nginx -t my-nginx .
docker-run:
	docker container run --rm -it --net=host --name my-nginx my-nginx
docker-compose-up:
	docker-compose up

.PHONY: docker-build
.PHONY: docker-run
.PHONY: docker-compose-up
