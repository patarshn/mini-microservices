# Checklist
| Feat / Tech                | Implemented | Desc                         |
|----------------------------|-------------|------------------------------|
| JWT                        |✅             | Authentication               |
| PostgreSQL                 |✅             | Database: SQL OLTP           |
| MongoDB                    |✅             | Database: NoSQL              |
| ClickHouse                 |✅             | Database: SQL OLAP           |
| MinIO                      |✅             | Object Storage               |
| Hot-Cold Storage           |✅             |Hot Storage: Hot archival makes this information available instantly. <br>Cold Storage: Cold archival is the long-term archival of data, and it is not intended to be accessed often.|
| Docker                     |✅             |Containerization |
| Microservices              |✅             |Architecture|
| Database Migration         |✅             |                              |
| Swagger UI                 |✅             | API Documentation            |
| Postman Collection         |✅             |                              |
| Unit Test                  |             |                              |
| Integration Test           |             |                              |
| SAGA Pattern               |             | Microservice pattern         |
| Apache Kafka               |             | Distributed Messaging System |
| Nginx                      |✅             | Web Server                   |
| Cron Job                   |✅             |                              |
| NodeJS (Fastify Framework) |✅             |Auth Service|
| Golang                     |✅             |Product Service, Transaction Service|
| MVC Architectur            |✅             |Auth Service | 
| Clean Architecture         |✅             |Product Service, Transaction Service|


# How To Start

### 1. Build Docker Compose
```sh 
docker-compose build
```

### 2. Run the database
```sh 
docker-compose up postgres clickhouse minio mongo -d
```
you can make sure it running or not with "docker ps"


### 3. Run migration
Sometimes, migration will fail because the database is not ready to use. For now, make sure the database is on hehe 🙃.
##### Run migration for product-service
``` 
docker-compose run product-service-migrate bash -c "../migrate"
```
##### Run migration for transaction service
``` 
docker-compose run transaction-service-migrate bash -c "../migrate"
```
##### Run migration for minio
``` 
docker-compose run minio-migrate
```

### 4. Run Service
```
docker-compose up auth-service product-service transaction-service transaction-service-cron nginx
```
Normally the service will run on the following port.
|Service | PORT | HOST| NGINX |
|--------|------|-----|-------|
|Auth Service| 8081 | localhost:8081 | localhost:8080/auth |
|Product Service| 8082 | localhost:8082 | localhost:8080/product |
|Transaction Service| 8083 | localhost:8083 | localhost:8080/transaction |

# How To Try
```
1. Export microservice-project.postman_collection.json to postman
2. Register your account with /auth/register endpoint
3. Get your jwt from /auth/login endpoint
4. set variable bearer_token with jwt value from login reponse
5. Now, you can use the API.
```
