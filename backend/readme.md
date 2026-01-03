

### Menjalankan Migrate UP
```shell
migrate -path backend/db/migrations -database "postgresql://username:password@host:port/namadb?sslmode=disable" up
```
### Migrate Down
```shell
migrate -path backend/db/migrations -database "postgresql://username:password@host:port/namadb?sslmode=disable" down -all
````

### Menjalankan Unit test
```shell
go test -v ./internal/modules/user/...
```