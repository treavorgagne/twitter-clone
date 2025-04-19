# twitter-clone
MVP dockerized twitter-clone.

## Stack
- client Next.js using React.js
- restful go api
- db mysql
- redis cache

## Deployment
Here are the commands in order to run each component of this project. 

### .env setup
Create a `.env` file with the following keys needed for the `compose.yaml` and go files in the root directory of the project.

```yaml
DBUSER=
DBPASS=
DBADDRESS=
DBPORT=
REDISPORT=
REDISADDRESS=
```
----
### Docker MySQL Database
Run the following docker compose file from the `server` directory. 

```bash
docker compose -f ./db/compose.yaml --env-file .env up
```
----
### Redis Cache
Run the following docker compose file from the `server` directory. 

```bash
docker compose -f ./redis/compose.yaml --env-file .env up
```
This is used to cache `GET` requests to the go server.
----
### Server
With `go` installed and cd into `server` directory. Then run the follwing to start the server:

```
go run main.go
```
----
### Client
----
