export SPANNER_PROJECT=test-project

run-emulator:
	docker run -p 9010:9010 -p 9020:9020 gcr.io/cloud-spanner-emulator/emulator:1.1.1

run-test:
	SPANNER_EMULATOR_HOST=localhost:9010 go test ./...