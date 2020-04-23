build:
	@if [ ! -d bin ]; then echo "creating bin folder"; mkdir bin ; fi;
	@if [ -d bin/server ]; then rm bin/server; fi;
	@echo "building server"; 
	@cd cmd/server; \
	go build -o ../../bin/server; 
	@echo "building client"; 
	@cd cmd/client; \
	go build -o ../../bin/client; 

build_run:
	make build && bin/server