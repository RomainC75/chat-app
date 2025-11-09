Chat App in progress - started 25 october 2025

## Chat

This chat app is just an excuse to implement : 
- Domain Driven Design principles,
- TDD principles,
- CQRS,
- CICD pipelines,
- Port&Adapters,
- some Design patterns

in other words : to implement some good stuff ;-)

## Lunch on localhost 
```
cp .env.template .env
cp back/config/config.yml.template back/config/config.yml

docker compose build 
docker compose up
```

go to http://localhost:5173

Todo : 

- temp adapters for persistence,
- error handlers,
- test-containers for repo,
- CQRS
- Context timeout for requests

- save each message in db