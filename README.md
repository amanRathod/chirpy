# chirpy

run build: `go build -o out && ./out`

migration: `goose -dir ./sql/schema postgres "postgres://local_user:123451@localhost:5432/chirpy" up`
