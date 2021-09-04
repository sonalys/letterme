test:
	make -C ./account_ms test
	make -C ./email_ms test
	make -C ./domain test

ci:
	make -C ./account_ms coverage
	make -C ./email_ms coverage
	make -C ./domain coverage