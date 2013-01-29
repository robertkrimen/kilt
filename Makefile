.PHONY: test clean install

test:
	go test -i 
	go test -v 


clean:

install: test
	go install
