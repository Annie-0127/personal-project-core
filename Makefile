build:
	docker-compose build
start:
	docker-compose up -d
ssh-app:
	docker exec -it personal-project-core-app bash

generate:
	wire ./services/wire

swagger:
	redocly bundle ./api/open-api.yaml -o bundled.yaml

gen-auth:
	oapi-codegen --config ./api/auth/config.yaml ./api/auth/api.yaml
gen-users:
	oapi-codegen --config ./api/user/config.yaml ./api/user/api.yaml