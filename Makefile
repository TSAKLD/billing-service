db:
	docker run -v avito:/var/lib/postgresql/data/ -p "5432:5432" -e POSTGRES_PASSWORD=asdbnm321 -e POSTGRES_USER=kr -e POSTGRES_DB=avito -d postgres:14.2
up:
	go run main.go