# twitter-clone
MVP dockerized twitter-clone.

## Stack
- client -> TBD
- server -> go (rest api)
- db -> mysql
- cache -> redis

## Deployment
Complete the following instructions listed below: 

1) [Dependencies](#dependencies)
2) [Project Setup Up](#project-setup-up)
3) [Containerized Deployment](#containerized-deployment)

----
### Dependencies
These must be downloaded
- [go](https://go.dev/doc/install)
- [docker](https://docs.docker.com/get-started/get-docker/)

----
### Project Setup Up

#### Step 1 - Clone Project Repo
Run the following from the directory you want you the project to stored:
```bash
git clone https://github.com/treavorgagne/twitter-clone.git
``` 

#### Step 2 - Create Docker Network
Run the following docker cmd:
```bash
docker network create tc-network
```

#### Step 3 - Create `.env` File
Add `.env` file in the root directory of the project using the following template to go after:

```yaml
DBUSER=<USER>
DBPASS=<PASSWD>
DBADDRESS=twitter-clone-db
DBPORT=3308
REDISPORT=6379
REDISADDRESS=redis
SERVERPORT=8080
```

---- 
### Containerized Deployment
Execute the following steps containerize the server, db, and cache:

```bash
# build go server binary and container image
cd server && go mod download && go build -o go-server && docker build . -t go-server:0.0.1 && cd ..
# starts server container from the image we just built
docker compose -f compose.yaml --env-file .env up -d 
```