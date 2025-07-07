# Makefile

.PHONY: up down restart logs db psql reset help

# Запуск всех контейнеров
up:
	docker-compose up -d

# Остановка и удаление контейнеров
down:
	docker-compose down

# Перезапуск всех контейнеров
restart: down up

# Просмотр логов всех сервисов
logs:
	docker-compose logs -f

# Подключение к psql внутри PostgreSQL-контейнера
psql:
	docker exec -it bounty-postgres psql -U postgres -d bounty

# Очистка volume
reset:
	docker-compose down -v

reset-db:
	docker-compose down -v
	rm -rf ./pgdata
	docker-compose up -d --build
	sleep 2
	docker exec -it bountyboard-db psql -U postgres -c "CREATE ROLE admin WITH LOGIN PASSWORD 'pirate-pwd';"
	docker exec -it bountyboard-db psql -U postgres -c "CREATE DATABASE db OWNER admin;"

super-reset:
	docker-compose down -v
	#docker volume rm bountyboard_pgdata
	docker-compose up --build

# Справка по командам
help:
	@echo "make up       # Запуск всех сервисов"
	@echo "make down     # Остановка и удаление сервисов"
	@echo "make restart  # Перезапуск"
	@echo "make logs     # Просмотр логов"
	@echo "make psql     # Войти в psql внутри контейнера"
	@echo "make reset    # Полная остановка + очистка volume (удаление данных)"
