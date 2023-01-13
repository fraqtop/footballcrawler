ARGS=$(filter-out $@, $(MAKECMDGOALS))

up:
	docker-compose up -d

down:
	docker-compose down

ps:
	docker-compose ps

db-status:
	docker-compose run --rm migration-tool goose status

db-up:
	docker-compose run --rm migration-tool goose up

db-modify:
	docker-compose run --rm migration-tool goose create ${ARGS}

db-down:
	docker-compose run --rm migration-tool goose down

run-cli-seed:
	docker-compose run --rm compiler go run . seedCompetitions

run-cli-stats:
	docker-compose run --rm compiler go run . loadStats
