DB_USER=root
DB_PASS=password
DATABASE=web_api
MYSQL = mysql --user=$(DB_USER) --password=$(DB_PASS)

initdb:
	$(MYSQL) --execute "CREATE DATABASE IF NOT EXISTS $(DATABASE);"
	mysql $(DATABASE) < raw.sql

cleandb:
	mysql $(DATABASE) < clean.sql

run:
	go run main.go



