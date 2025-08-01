APP_NAME = "gwid-core"

MAIN_FILE = "./cmd/server/main.go"

TMP_DIR = "tmp"

dev: tidy
	@echo "Starting development server with hot-reloading..."
	air -c .air.toml

dev-docker:
	docker compose up

down:
	docker compose down

build:
	docker compose build

logs:
	docker compose logs -f

# logs-db:
# 	docker compose logs -f postgres

clean:
	docker compose down -v --rmi all --remove-orphans
	docker system prune -f

tidy:
	@echo "Tidying modules..."
	go mod tidy
	@echo "Tidy complete!"

restart-gwid:
	docker compose restart gwid

shell:
	docker compose exec gwid bash

# db-shell:
# 	docker compose exec postgres psql -U user -d gwid

ps:
	docker compose ps
