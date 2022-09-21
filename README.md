# Dining hall

Laboratory work nr1 on Network programming course.

## Build and run application in docker

```bash
docker build -t dining-hall .
docker run -p 8081:8081 -it dining-hall
```

For Linux:

```bash
docker build -t dining-hall .
docker run --add-host host.docker.internal:host-gateway -p 8081:8081 -it dining-hall
```

## Run application locally

```bash
go run .
```

## URL

```bash
http://host.docker.internal:8081
```

If it gives ECONNREFUSED (connection refused), the workaround is to find host.docker.internal in the hosts file and replace the ip attributed to it with 127.0.0.1
