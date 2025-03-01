COMPOSE = docker compose
SERVICE = app

build:
	$(COMPOSE) build --no-cache

up:
	$(COMPOSE) up -d

down:
	$(COMPOSE) down

rebuild: down build up

logs:
	$(COMPOSE) logs -f $(SERVICE)

shell:
	$(COMPOSE) exec $(SERVICE) sh

clean:
	$(COMPOSE) down --rmi all --volumes --remove-orphans

# Команды для локального PostgreSQL (ZKT)
zkt-up:
	$(COMPOSE) -f docker-compose.attendance-db.local.yml up -d
	@echo "PostgreSQL запущен"

zkt-down:
	$(COMPOSE) -f docker-compose.attendance-db.local.yml down -v
	@echo "PostgreSQL остановлен"

zkt-init:
	@echo "Инициализация базы данных..."
	@./init-attendance-trigger.sh

zkt-logs:
	$(COMPOSE) -f docker-compose.attendance-db.local.yml logs -f --tail=50