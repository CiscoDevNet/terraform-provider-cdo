package accesspolicies

import "github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/cloudfmc/internal"

type AccessPolicies = internal.Response

var Builder = internal.NewResponseBuilder
var New = internal.NewResponse

type Links = internal.Links
type Paging = internal.Paging
type Item = internal.Item

var NewItem = internal.NewItem
var NewLinks = internal.NewLinks
var NewPaging = internal.NewPaging
