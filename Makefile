test:
	make -C ./account_ms test
	make -C ./email_ms test
	make -C ./domain test

ci:
	make -C ./account_ms ci
	make -C ./email_ms ci
	make -C ./domain ci