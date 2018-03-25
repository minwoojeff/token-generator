# VARS #
ifeq ($(CLIENT_ID),)
	export CLIENT_ID=default
endif
ifeq ($(CLIENT_SECRET),)
	export CLIENT_SECRET=password
endif


# COMMANDS #

default:
	go build
	./token-generator

test:
	go test

coverage:
	go test -cover
