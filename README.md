# Dining hall
Laboratory work nr1 on Network programming course.

# Build and run application in docker
```
docker build -t dining-hall .
docker run -p 8081:8081 -it dining-hall
```
For Linux: 
```
docker build -t dining-hall .
docker run --add-host host.docker.internal:host-gateway -p 8081:8081 -it dining-hall
```

# Run application locally
```
go run .
```

# URL
```
http://host.docker.internal:8081
```