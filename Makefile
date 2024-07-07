run:
	docker-compose --env-file .env up --build

test:
	docker-compose -f tests/docker-compose.yaml --env-file .env up --build --abort-on-container-exit --exit-code-from postman