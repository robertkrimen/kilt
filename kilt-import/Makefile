.PHONY: test release install clean

test:
	go test -i
	go test -v

release: test
	godocdown -signature . > README.markdown

install: test
	go install .

clean:
	rm -rf kilt
