# LOANS REST API

This repo contains solutions of "Loans" REST API on Golang

## How 2RUN

1. Clone project

```
git clone git@github.com:rabdavinci/loan.git .
```

2. Run microservice

```
$ go run cmd/main.go
```

3. Get all loans

```
GET localhost:9090

```

4. Create loan

```
POST localhost:9090
BODY {"product":"Смартфон","phone":"+998995881375","month":12,"price":1000}
```

## TODO

1. Finish transactions
2. Add tests
3. Move storage to Database, use PG for example
4. Use config package, env for keeping params
5. Dockerize
