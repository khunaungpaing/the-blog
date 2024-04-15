run: postgres
	CompileDaemon -command=./the-blog-api

postgres:
	docker compose up -d

.PHONY: postgres createdb run
