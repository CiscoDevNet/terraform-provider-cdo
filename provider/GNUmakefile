default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	USER_DATA_SOURCE_USERNAME=terraform-provider-cdo-super-admin@lockhart.io \
    USER_DATA_SOURCE_ROLE=ROLE_SUPER_ADMIN \
    USER_DATA_SOURCE_IS_API_ONLY=false \
	TF_ACC=1 \
	go test ./... -v $(TESTARGS) -timeout 120m -count 1

TESTARGS=-run TestAccUserResource