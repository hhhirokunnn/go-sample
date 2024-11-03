https://go.dev/doc/tutorial/database-access

```bash
DBUSER=root go run main.go request_handler.go repository.go
```

```bash
DBUSER=root go run migration.go
```

```bash
curl https://localhost:8080/albums

curl http://localhost:8080/albums \
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"id": "4","title": "The Modern Sound of Betty Carter","artist": "Betty Carter","price": 49.99}'
```

```bash
mysql -uroot -h127.0.0.1
```