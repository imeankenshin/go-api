build:
	go build -o ./bin/serve .

run:
	air -c ./.air.toml

sql:
	sqlite3 todos.db
