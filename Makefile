start_images:
	docker run --rm -d -p ${MONGO_PORT}:27017 --name letterme_mongo mongo:5.0 || true
	docker run --rm -d -p 5672:5672 -p 15672:15672 --name rabbitmq rabbitmq:3.9-management || true

test:
	make start_images
	make -C ./account_ms test
	make -C ./email_ms test
	make -C ./domain test

ci:
	make -C ./account_ms coverage
	make -C ./email_ms coverage
	make -C ./domain coverage