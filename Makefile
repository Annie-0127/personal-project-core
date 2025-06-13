build:
	docker-compose build
start:
	docker-compose up -d
ssh-app:
	docker exec -it personal-project-core-app bash

generate:
	wire ./services/wire