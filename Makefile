# VARS #
ifeq ($(CLIENT_ID),)
	export CLIENT_ID=default
endif
ifeq ($(CLIENT_SECRET),)
	export CLIENT_SECRET=password
endif
ifeq ($(CLIENT_URL),)
	export CLIENT_URL=http://localhost:3000
endif

# COMMANDS #

default:
	./scripts/bootstrap.sh
	go build
	./token-generator

test:
	go test

coverage:
	go test -cover

acceptance:
	pip3 install -q -r requirements.txt
	python3 acceptance_test.py ${CLIENT_URL}

clean:
	rm -rf tmp/
