package accesspolicies

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"

type AccessPolicies = internal.Response

var NewAccessPoliciesBuilder = internal.NewResponseBuilder

type Links = internal.Links
type Paging = internal.Paging
type Item = internal.Item

var NewLinks = internal.NewLinks
var NewPaging = internal.NewPaging
