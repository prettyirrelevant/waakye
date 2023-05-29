BASE_DIR = $(shell pwd)
KILISHI_DIR = $(BASE_DIR)/kilishi
MASA_DIR = $(BASE_DIR)/masa
ASARO_DIR = $(BASE_DIR)/asaro

run-kilishi:
	@cd $(KILISHI_DIR) && air -c .air.toml

run-asaro:
	@cd $(ASARO_DIR) && flask run

run-masa:
	@cd $(MASA_DIR) && npm run dev
