default: testacc

# https://unix.stackexchange.com/questions/235223/makefile-include-env-file
include	.github-action.env  # common variables
include .github-action.local.env  # secret variables
export

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 \
	go test ./... -v $(TESTARGS) -timeout 120m -count 1

TESTARGS=-run TestAccFtdDataSource