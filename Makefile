.PHONY: test

test:
	make -C ./account_ms test

ci:
	make -C ./account_ms ci
	make -C ./domain ci