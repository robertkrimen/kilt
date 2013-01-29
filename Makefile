.PHONY: test clean install

test:
	go test -i 
	go test -v 
	$(MAKE) -C kilt-import $@


clean:
	$(MAKE) -C kilt-import $@

install: test
	go install
	$(MAKE) -C kilt-import $@
