# L0. test task
![Static Badge](https://img.shields.io/badge/go-1.25-blue?color=%2300ADD8)
![Static Badge](https://img.shields.io/badge/redis-8.2.1-blue?color=%23FF4438)
![Static Badge](https://img.shields.io/badge/postgresql-17.5-blue?color=%234169E1)
![Static Badge](https://img.shields.io/badge/kafka-7.3.0-blue?color=%23231F20)
## task
Необходимо разработать демонстрационный сервис с простейшим интерфейсом, отображающий данные о заказе.
Данное задание предполагает создание небольшого микросервиса на Go с использованием базы данных и очереди сообщений. 
Сервис будет получать данные заказов из очереди (Kafka), сохранять их в базу данных (PostgreSQL) и кэшировать в памяти для быстрого доступа.

### run
1. copy code
```bash
git clone https://github.com/Killazius/L0.git && cd L0
```
2. make .env file
```bash
cp .env.example .env
```
```env
POSTGRES_HOST="postgres"
POSTGRES_PORT="5432"
POSTGRES_USER="myuser"
POSTGRES_PASSWORD="mypassword"
POSTGRES_DB="mydb"
POSTGRES_SSL_MODE="disable"

REDIS_ADDR="redis:6379"
REDIS_PASSWORD="redis"
REDIS_DB=0
```
3. run service
```bash
make docker OR docker compose up -d
```

### api endpoints
```
GET /order/{order_uid} - данные по заказу
GET / - веб-интерфейс
GET /swagger/ - документация swagger
```

### kafka producer
для того, чтобы отправить сообщения в кафку, написан скрипт, запустить его можно путем команды:
```bash
make produce
```
`COUNT=(x). default COUNT = 3`


