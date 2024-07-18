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

# How To Try
```
1. Export microservice-project.postman_collection.json to postman
2. Register your account with /auth/register endpoint
3. Get your jwt from /auth/login endpoint
4. set variable bearer_token with jwt value from login reponse
5. Now, you can use the API.
```
