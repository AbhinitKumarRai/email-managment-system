all:
	docker-compose up --build -d

kill:
	docker-compose down 
