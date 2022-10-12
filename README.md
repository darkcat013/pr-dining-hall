# Dining hall

Laboratory work nr2 on Network programming course.

## Build and run application in docker

```bash
docker compose up --build
```

## Run application locally

```bash
go run .
```

## URLs

```url
http://host.docker.internal:8081
http://host.docker.internal:8083
http://host.docker.internal:8085
http://host.docker.internal:8087
```

If it gives ECONNREFUSED (connection refused), the workaround is to find host.docker.internal in the hosts file and replace the ip attributed to it with 127.0.0.1

## Local URLs

```url
http://localhost:8081
http://localhost:8083
http://localhost:8085
http://localhost:8087
```
