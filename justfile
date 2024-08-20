build:
  go build -o api -v

test-all:
  go test ./tests/ -v

test TEST:
  go test ./tests/{{TEST}} -v

migrate: build
  ./api migrate

dev:
  air
