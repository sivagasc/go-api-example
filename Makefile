build:
	@if [ ! -d bin ]; then echo "creating bin folder"; mkdir bin ; fi;
	@if [ -d bin/server ]; then rm bin/server; fi;
	@echo "building server"; 
	@cd cmd/server; \
	go build; 
	@mv cmd/server/server bin/server;
