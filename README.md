# Hotel reservation backend

## Project outline

- user -> book room from an hotel
- admins -> going to check reservation/bootings
- Authentication and authorization -> JWT tokens
- Hotels -> CRUD api -> JSON
- rooms -> CRUD api -> JSON
- Scripts -> database management -> seeding, migration

## Resources

### Mongodb driver

Documentaion

```
https://mongodb.com/docs/drivers/go/current/quick-start
```

Installing mongodb client

```
go get go.mongodb.org/mongo-driver/mongo
```

### gofiber

Documentaion

```
https://gofiber.io
```

Installing gofiber

```
go get github.com/gofiber/fiber/v2
```

## Docker

### Installing mongodb as a Docker container

```
docker run --name mongodb -d mongo:latest -p 27017:27017
```
