# VARS #
ifeq ($(CLIENT_ID),)
	export CLIENT_ID=default
endif
ifeq ($(CLIENT_SECRET),)
	export CLIENT_SECRET=password
endif
ifeq ($(CLIENT_URL),)
	export CLIENT_URL=http://localhost:80
endif

# COMMANDS #

default:
	./scripts/bootstrap.sh
	go build
	./token-generator

test:
	go test ./...

coverage_single:
	go test -coverprofile=c.out
	go tool cover -func=c.out

acceptance:
	pip3 install -q -r requirements.txt
	python3 acceptance_test.py ${CLIENT_URL}

clean:
	rm -rf tmp/ c.out
