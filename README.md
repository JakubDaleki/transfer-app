The purpose of this application is to build a system that could be used for efficient money transfers (something like bank platform). The main goal of this project is to learn golang and check how convinient the language is in developing microservices with popular dependencies like kafka or postgressql.


![System architecture](https://github.com/JakubDaleki/transfer-app/blob/main/arch-diagram.png?raw=true)

This is not a final architecture as it's still evolving. Another consideration is to use key-value in memory db to store balances. Each time this service goes down it would have to process kafka logs to retrieve current balance for each user. It can be costly but in this scenario I wouldn't expect KV store to be down very often and it would benefit in high performance thanks in-memory features. In addition, kafka allows to chain operations and perform account balance addition after subtraction (asynchroniously) and because of that standard, heavy db with distributed transaction features is not needed.

To run the application simply install docker and run command `docker compose up` in transfer-app directory and enjoy the following endpoints to play with the application!

## Endpoints

| Name          | HTTP Method  | Route            | Body Fields        | JWT required |
|---------------|--------------|------------------|--------------------|--------------|
| Register user | POST         | /register        | username, password |      ❌      |
| JWT Auth      | POST         | /authentication  | username, password |      ✅      |
| Get Balance   | GET          | /account/balance |                    |      ✅      |
| Make Transfer | POST         | /account/transfer| to, amount         |      ✅      |
