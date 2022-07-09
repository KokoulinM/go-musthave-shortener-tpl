# Postgres

Run server:

````
pg_ctl -D /opt/homebrew/var/postgres start
````

Start psql and open database postgres, which is the database postgres uses itself to store roles, permissions, and structure:
````
psql postgres
````
# Test

````
Check coverage
go test ./... -v -short -p 1 -cover
````
# Linting

````
golangci-lint run
````