DB_URL=postgres://postgres:1234@localhost:5432/bookmark_manager?sslmode=disable

migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrateup:
	migrate -path db/migrations -database "$(DB_URL)" -verbose up

migratedown:
	migrate -path db/migrations -database "$(DB_URL)" -verbose down

run:
	go run cmd/api/main.go
