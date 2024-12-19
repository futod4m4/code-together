.PHONY:

run:
	go run ./cmd/api/main.go $(config)

#======================================================
# Go migrate psql

force:
	migrate -database postgres://fedor:postgres@localhost:5432/code_together?sslmode=disable -path migrations force 1

version:
	migrate -database postgres://fedor:postgres@localhost:5432/code_together?sslmode=disable -path migrations version

migrate_up:
	migrate -database postgres://fedor:postgres@localhost:5432/code_together?sslmode=disable -path migrations up 1

migrate_down:
	migrate -database postgres://fedor:postgres@localhost:5432/code_together?sslmode=disable -path migrations down 1