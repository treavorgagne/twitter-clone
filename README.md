# twitter-clone
MVP dockerized twitter-clone.

## Stack
- client -> TBD
- server -> go (rest api)
- db -> mysql
- cache -> redis

## Deployment
Complete the following instructions listed below: 

1) [Download Launch Dependencies](#download-launch-dependencies)
2) [Pre Launch Steps](#pre-launch-steps)
3) [Launch Steps](#launch-steps)

----
### Download Launch Dependencies
These must be downloaded
- [go](https://go.dev/doc/install)
- [docker](https://docs.docker.com/get-started/get-docker/)

### Pre Launch Steps

#### Clone project repo
Run the following from the directory you want you the project to stored:
```bash
git clone https://github.com/treavorgagne/twitter-clone.git
``` 

#### Create Docker Network
Run the following docker cmd:
```bash
docker network create tc-network
```

#### Create .env file
Create a `.env` file in the root directory of the project.

```yaml
DBUSER=<USER>
DBPASS=<PASSWD>
DBADDRESS=twitter-clone-db
DBPORT=3308
REDISPORT=6379
REDISADDRESS=redis
SERVERPORT=8080
```

### Launch Steps
Excute the following steps after downloading dependencies, cloning repo, and create docker network:

#### Step 1: Containerized MySQL Database
Execute db docker compose file from the project root directory:

```bash
# starts mysql db container using compose file
docker compose -f ./db/compose.yaml --env-file .env up -d 
```

----
#### Step 2: Containerized Go Server
Execute go server docker compose file from the project root directory:

```bash
# build go server binary and container image
cd server && go mod download && go build -o go-server && docker build . -t go-server:0.0.1 && cd ..
# starts server container from the image we just built
docker compose -f ./server/compose.yaml --env-file .env up -d 
```

----
#### Step 3: Containerized Redis Cache (optional)
Execute redis docker compose file from the project root directory:

```bash
# starts redis cache container using compose file
docker compose -f ./redis/compose.yaml --env-file .env up -d 
```
This is used to cache `GET` requests to the go server.

----


