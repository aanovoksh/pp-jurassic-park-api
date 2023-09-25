# The Jurassic Park API - Go Backend

## Tech Stack

* Go v1.21.1
  * gin v1.9.1
  * gorm v1.25.4
  * testify v1.8.4
* PostgreSQL
  * latest docker image

## Getting Started

1. Ensure Docker is installed: [Docker Installation Guide](https://docs.docker.com/get-docker/)
2. Run the services (Go server + PostgreSQL) with Docker Compose
```sh
docker-compose up --build
```
3. Access the API at `http://localhost:8080`
4. Run tests inside docker container by 
```sh
docker-compose exec api go test ./...
```

## API Overview
Following API endpoints are built to manage the Jurassic Park

| Route | HTTP Method | Description |  
| ------ | ------ | ------ | 
| `/cages` | GET | Query all cage details, including enclosed dinosaurs. Filterable by power status. |
| `/cages/:id` | GET | Query single cage details, including enclosed dinosaurs. | 
| `/cages` | POST | Create a new cage. | 
| `/cages/:id` | PATCH | Update power status in the existing cage. | 
| `/cages/:id` | DELETE | Delete the cage. | 
| `/dinosaurs` | GET | Query all dinosaur details. Filterable by species. |
| `/dinosaurs/:id` | GET | Query single dinosaur details. |
| `/dinosaurs` | POST | Add new dinosaur to existing cage in the Park. | 
| `/dinosaurs/:id` | PATCH | Move dinosaur from one cage to another. | 
| `/dinosaurs/:id` | DELETE | Remove dinosaur from the Park. | 

**Note**: Postman Collection is added to the github repo as well as a top level file. Feel free to pull it, import and play with existing APIs.

## Codebase Structure
```
/pp-jurassic-park-api
├── /cmd
│   └── /ypi
│       └── main.go         # The application’s entry point
├── /internal                     
│   ├── /api
│   │   ├── handlers        # HTTP handlers for cages and dinosaurs 
│   │   ├── models          # API Models
│   │   └── stransform      # Helpers for model transformations
│   ├── /db
│   │   ├── models          # DB Models
│   │   └── db.go           # DB connection and migration logic
│   └── /tests              # End-to-end handler tests
├── go.mod
├── go.mod
├── dockerfile
└── docker-compose.yml
```

## Follow-ups and Improvements
* Handler Code Refactoring
  * Given limited time and scope for this project, and considering that the overall logic of the API is simple enough, I have made a decision to keep it all in just a handler. But as the API scales and the logic expands, it would be good to split core business logic and db logic from handler into service/repository.
* DB Schema Management
  * As of right now, DB schema is powered by two models and DB is auto-migrated on a server start. As the API scales, we should look into supporting DB versioning, migrations and rollbacks in a more resilient and scalable way.
* Test Data Lifecycle
  * As of right now, the test data is being set once before all the tests are run. This is the simplest options to start with, but as the API scales, it would be good to revisit it and have a proper test data setup and teardown for each individual test case.
* Logging and Metrics
  * Adding proper logging and metrics will help the team to monitor health of the system and triage any issues if they were to occur.
* DB connection strings
  * As of right now, connection strings are just hardcoded into the codebase. Moving it to a secure place is a must before releasing this code.
* Swagger Docs
  * Adding automated API documentation will help the team to keep track of existing APIs and have proper documentation of it.
* Linting
  * Adding automated linting can ensure that consistent styling is used across the entire application and that best practices are being followed.
* RBAC 
  * Right now, there is no access control to the API, meaning that anyone can call any endpoints. Adding role-based access control will help with ensuring proper access to the tools.
* API updates to support concurrency (if needed)
  * Right now, the assumption is taken that Cages and Dinosaurs will be managed by a fairly small group of scientists and doctors (not a lot of dino experts out there anyways), thus the APIs are not fully safe from race conditions. If the park will scale rapidly, we can update the APIs to use BD locking as a safeguard for concurrent updates:
    * DB Locking: records can be locked for the duration of the transaction. It will prevent other threads from updating the same records. Openning DB transaction and selecting records "FOR UPDATE" is the common way of locking the records.
* Dinosaur Species Management
  * As of right now, all supported species are hardcoded in the apimodels package. It is an acceptable solution for a given exercise, but if we were to scale the park and if we know that the list of known species will grow over time, it would be better to make it more configurable. For example, all species and whether they are carnivores or herbivores can be moved to DB and CRUD api around known species can be added.
