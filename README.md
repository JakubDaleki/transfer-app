## Introduction
The purpose of this application is to build a system that could be used for efficient money transfers (something like bank platform). The main goal of this project is to learn golang and check how convinient the language is in developing microservices with popular dependencies like kafka or postgressql.

## Architecture overview

This is not a final architecture as it's still evolving. In the first approach I wanted to go with standard database instead of in-memory db. The reasons I rejected that solution were: kafka that already stores all the information about transactions on disk and performance of in-memory database. Now each time the `query-service` is restarted it needs to aggregate on kafka's log to retrieve up to date balance. For that purpose all transfers for a user are kept on the same partition to speed up summations. Thus, all transfer are added to kafka twice: for the first time to subtract an amount from a sender and later to make the addition for a receiver.

![System architecture](https://github.com/JakubDaleki/transfer-app/blob/main/arch-diagram.png?raw=true)

If this architecture hits performance limit, it can be easily scaled by adding more instances (partitions) of `query-service`, each partition containing a data for different subset of users.

## Endpoints

| Name          | HTTP Method  | Route            | Body Fields        | JWT required |
|---------------|--------------|------------------|--------------------|--------------|
| Register user | POST         | /register        | username, password |      ❌      |
| JWT Auth      | POST         | /authentication  | username, password |      ✅      |
| Get Balance   | GET          | /account/balance |                    |      ✅      |
| Make Transfer | POST         | /account/transfer| to, amount         |      ✅      |

## Starting the application

To run the application simply install docker and run command `docker compose up` in transfer-app directory and enjoy the following endpoints to play with the application!