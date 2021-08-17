mongo_port = 2564
mongo_address = "localhost:${mongo_port}"

start_mongo:
	docker rm -f letterme_mongo || true
	docker run --rm -d -p ${mongo_port}:27017 --name letterme_mongo mongo:5.0

test:
	LM_MONGO_SETTINGS='{ "hosts": [${mongo_address}], "db_name": "account_manager", "session_name": "account_manager_ms" }' \
	go test ./...