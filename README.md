# weather-on-mars-interface
Simple GraphQL service that serves as an interface to NASA's InSight: Mars Weather Service API​.

### To get everything running, follow these steps:

1. By running `docker-compose up -d` we pull & build all required images and bring up all the
containers for the application. 

2. By running `go run .` we run API Server. API is exposed on port: 8080.

### TODO:

This API is still work in progress. We need to add a few additional requirements to fulfill functional requirements. First to handle are:

- database management as defined in functional requirements
- pagination
- GraphQL Subscription
- Unit tests
- Code refactoring