Chat App in progress - started 25 october 2025


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
- hexagonal architecture,
- DDD 
- CQRS
- Context timeout for requests
  

- save each message in db