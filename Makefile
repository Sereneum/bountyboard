up:
	docker-compose up -d

down:
	docker-compose down

restart: down up

logs:
	docker-compose logs -f

psql:
	docker exec -it bounty-postgres psql -U $(POSTGRES_USER) -d $(POSTGRES_DB)

reset-db:
	docker-compose down -v
	@if [ "$(OS)" = "Windows_NT" ]; then \
		powershell.exe -Command "Remove-Item -Recurse -Force ./pgdata"; \
	else \
		rm -rf ./pgdata; \
	fi
	docker-compose up -d --build

help:
	@echo "Makefile commands:"
	@echo "  up        - Start all containers"
	@echo "  down      - Stop and remove containers"
	@echo "  restart   - Restart all containers"
	@echo "  logs      - Show logs"
	@echo "  psql      - Connect to PostgreSQL"
	@echo "  reset-db  - Fully reset DB and remove ./pgdata"

