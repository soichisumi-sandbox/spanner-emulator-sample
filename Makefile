export SPANNER_PROJECT=test-project

start-emulator:
	docker run --rm -p 9010:9010 -p 9020:9020 --name spanner-emulator gcr.io/cloud-spanner-emulator/emulator:1.1.1

stop-emulator:
	docker stop spanner-emulator

run-test:
	SPANNER_EMULATOR_HOST=localhost:9010 go test ./...