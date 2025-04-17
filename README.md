# twitter-clone
MVP Twitter clone project to perform CRUD commands. Features, like a login will be omitted. For now this only has a single backend, but in the future I'm going to add interoperable components in different languages/frameworks that can be used interchangeably.

## Stack
- client Next.js using React.js
- server restful go api 
- db mysql

## Deployment

Here are the commands in order to run each component of this project.

### Docker MySQL Database
With `docker` and `mysql` installed run the following to create a mysql container from the lastest mysql image:

```bash
docker run -p 3308:3306 --name twitter-clone-db-container \
  -e MYSQL_ROOT_PASSWORD='<your_password>' \
  -e MYSQL_DATABASE=twitter \
  -d mysql:latest
```
or from `./server` directory
```bash
docker compose -f ./db/compose.yaml --env-file ./.env up
```

Then cd into the `./server/db` directory and add the SQL tables by running the following: 

```bash
mysql -u root -p < twitter.sql
```

----
### Server
With `go` installed and cd into `server` directory. Then run the follwing to start the server:

```
go run main.go
```
----
### redis cache
To start the redis cache run the following:

```bash
docker run -d --name my-redis -p 4444:6379 redis
```
or from `./server` directory
```bash
docker compose -f ./redis/compose.yaml
```
This is used to cache `GET` requests to the go server.

----
### Client

