# How to run it

1. `docker-compose -f deployment/docker-compose.yaml up -d` to start service and db
2. `docker-compose -f deployment/docker-compose.yaml logs -f` 
   wait until `database system is ready to accept connections` this message appear in db_1 twice.
3. `sh script/run_migration.sh` to run db migration
4. After all the steps above are done, you can `curl localhost:8080` with these APIs
   1. `POST /api/v1/score` to update score in json (`{ "score": 100 }`) with `ClientId:<clientId>` Header
   2. `GET /api/vi/leaderboard` to retrieve top10 score leaderboard
