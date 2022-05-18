run:
	go run main.go

mock:
	mockgen -source=database/db_interface.go -destination=database/mocks/db_mock.go -package=mocks