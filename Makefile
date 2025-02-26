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