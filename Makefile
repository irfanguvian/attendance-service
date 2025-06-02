serve:
	CompileDaemon -command="./attendance-service" -color=true

run-migration:
	go run ./migrate/migrate.go