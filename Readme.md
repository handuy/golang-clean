### CRUD APICRUD API

1. Khởi tạo MySQL

```bash
docker run -e MYSQL_ROOT_PASSWORD=123456 -e MYSQL_USER=duy -e MYSQL_PASSWORD=123456 -e MYSQL_DATABASE=todolist -d  -p 3306:3306 -v crud:/var/lib/mysql -v $PWD/init.sql:/docker-entrypoint-initdb.d/init.sql mysql:5
```

2. Chạy app

```go
go run main.go
`