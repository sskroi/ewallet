# Description
API электронного кошелька.
Endpoints:
* `/api/v1/wallet` POST - создание кошелька, по умолчанию на кошельке 100 у.е.
Ответ: статус 200 + json объект с полями "id" - строковый id кошелька, "balance" - число, баланс кошелька.
* `/api/v1/wallet/{walletId}/send` POST - перевод средств. {walletId} - id исходящего кошелька, запрос должен содержать json объект с полями "to" - строковый id входящего кошелька и "amount" - число, сумма перевода. 
Ответы: 1) статус 404, если исходящий кошелек не найден 2) статус 400 если входящий кошелек не найден или на исходящем недостаточно средств, 3) статус 200 если перевод прошёл успешно
* `/api/v1/wallet/{walletId}` GET - проверка баланса кошелька. {walletId} - id кошелька. Ответы: 1) статус 404, если кошелек не найден 2) статус 200 + json объект с полями "id" - строковый id кошелька, "balance" - число, баланс кошелька.
* `/api/v1/wallet/{walletId}/history` GET - история переводов. {walletId} - id кошелька.
Ответы: 1) статус 404, если кошелек не найден 2) статус 200 + json объект - массив объектов с полями: "from" - id исходящего кошелька, "to" - id входящего кошелька, "amount" - сумма перевода, "time" - время перевода в формате RFC 3339

Конфиг файл лежит в app/configs/config.yaml, директория с базой данных монитруется в ./databaseData, можно изменить в `docker-compose.yaml`

# Launch
* `docker-compose up --build`

# Demo
```
$ curl -X POST -i localhost:8000/api/v1/wallet

HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Fri, 01 Mar 2024 21:47:57 GMT
Content-Length: 59

{"id":"6abac28f-6c36-43f9-b7ff-a75ea559cf5d","balance":100}
```
```
$ curl -X POST localhost:8000/api/v1/wallet/6abac28f-6c36-43f9-b7ff-a75ea559cf5d/send 
-d '{"to": "45b68c5b-39f4-44b3-8ca0-b20bb835cd2c", "amount": 15}'

{"message":"successful transfer"}
```
```
$ curl -X GET localhost:8000/api/v1/wallet/6abac28f-6c36-43f9-b7ff-a75ea559cf5d/history

[{"time":"2024-03-01T21:55:08.496649Z",
"from":"6abac28f-6c36-43f9-b7ff-a75ea559cf5d",
"to":"45b68c5b-39f4-44b3-8ca0-b20bb835cd2c",
"amount":15},
{"time":"2024-03-01T21:56:33.889395Z",
"from":"6abac28f-6c36-43f9-b7ff-a75ea559cf5d",
"to":"45b68c5b-39f4-44b3-8ca0-b20bb835cd2c",
"amount":7}]
```
