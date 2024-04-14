run: postgres
	CompileDaemon -command=./the-blog-api

createdb:postgres
	docker exec -it postgres createdb --username=root --owner=root blog-api

postgres:
	docker compose up -d

.PHONY: postgres createdb run
