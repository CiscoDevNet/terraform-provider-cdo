package auth

import (
	"encoding/json"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/featureflag"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/user/auth/role"
	"strings"
)

type Info struct {
	UserAuthentication Authentication `json:"userAuthentication"`
}

type Authority struct {
	Authority role.Type `json:"authority"`
}

type Authentication struct {
	Authorities   []Authority `json:"authorities"`
	Details       Details     `json:"details"`
	Authenticated bool        `json:"authenticated"`
	Principle     string      `json:"principal"`
	Name          string      `json:"name"`
}

type Details struct {
	TenantUid              string `json:"TenantUid"`
	TenantName             string `json:"TenantName"`
	SseTenantUid           string `json:"sseTenantUid"`
	TenantOrganizationName string `json:"TenantOrganizationName"`
	TenantDbFeatures       string `json:"TenantDbFeatures"`
	TenantUserRoles        string `json:"TenantUserRoles"`
	TenantDatabaseName     string `json:"TenantDatabaseName"`
	TenantPayType          string `json:"TenantPayType"`
}

func (info *Info) HasFeatureFlagEnabled(featureFlag featureflag.Type) bool {
	featureMap := info.getFeatureMap()
	enabled, ok := featureMap[featureFlag.String()]
	return ok && enabled
}

func (info *Info) getFeatureMap() map[string]bool {
	featureMap := map[string]bool{}
	err := json.Unmarshal([]byte(info.UserAuthentication.Details.TenantDbFeatures), &featureMap)
	if err != nil {
		// feature flag received is not a json
		panic(fmt.Sprintf("feature flag received from authentication service is not in valid format: %s", info.UserAuthentication.Details.TenantDbFeatures))
	}
	// convert feature flag keys to lowercase
	// TODO: make use of lh-feature
	normalizedFeatureMap := map[string]bool{}
	for k, v := range featureMap {
		normalizedFeatureMap[strings.ToLower(k)] = v
	}

	return normalizedFeatureMap
}
