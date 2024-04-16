run: postgres
	CompileDaemon -command=./the-blog-api

postgres:
	docker compose up -d

doc:
	swag init --parseDependency

.PHONY: postgres createdb run doc
