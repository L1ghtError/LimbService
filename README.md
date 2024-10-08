## Limb Backend
### Image procesing service written on Go
Works with [Frontend](https://github.com/L1ghtError/LimbAppWeb). Workers source code recently not available.
### How to build
native:
```bash
$ go build .
$ .\light-backend
```
Docker:
```bash
** in future updates**
```
Docker-compose:
```bash
** in future updates**
```

>[!NOTE]
>**Do not forget to rename exmple.env to config.env and fill it with right values**

> **Tech stack:**
> - [Fiber](https://github.com/gofiber/fiber) as web-framework
> - [amqp091-go](https://github.com/rabbitmq/amqp091-go) for communication with workers
> - [mongo-driver **Version 1**](https://github.com/mongodb/mongo-go-driver) for communication with MongoDb
> - [OAuth2 ](https://github.com/golang/oauth2) for Google OAuth
> - All smaller dependencies can be found in [go.mod](https://github.com/L1ghtError/LimbService/blob/main/go.mod)
