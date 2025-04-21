LOCAL_BIN:=$(CURDIR)/bin


generate:
	make generate-migration

generate-migration:
	./bin/goose create ____ sql