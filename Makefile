up:
	docker compose up

upd:
	docker compose up -d

build:
	go build -o ./tmp/main .

run:
	air -c ./.air.toml
