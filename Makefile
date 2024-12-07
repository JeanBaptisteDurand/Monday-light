re:
	docker-compose down && docker-compose up --build -d

rev:
	docker-compose down -v && docker-compose up --build -d