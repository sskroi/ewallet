App config: app/configs/config.yaml

Port: 8000

Postgres:
    user: postgres
    password: qwerty
    port: 5433
    dbname: postgres

Default mounting point for postgres: ./databaseData

Run: docker-compose up --build


Information about the balance of wallets is stored
in the database with an accuracy of two decimal places. 
The "amount" parameter will be rounded to two decimal
places during the transfer.