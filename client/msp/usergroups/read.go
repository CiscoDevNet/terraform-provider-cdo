package usergroups

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	mapset "github.com/deckarep/golang-set/v2"
)

func ReadCreatedUserGroupsInTenant(ctx context.Context, client http.Client, tenantUid string, readInput *[]MspManagedUserGroupInput) (*[]MspManagedUserGroup, error) {
	client.Logger.Printf("Reading user groups in tenant %s\n", tenantUid)
	// create a map of the user groups that were created
	// find the list of deleted users by removing from the list every time a user is found in the response
	readUserGroupDetailsMap := map[string]MspManagedUserGroup{}
	for _, createdUserGroup := range *readInput {
		readUserGroupDetailsMap[createdUserGroup.GroupIdentifier] = MspManagedUserGroup{
			GroupIdentifier: createdUserGroup.GroupIdentifier,
			IssuerUrl:       createdUserGroup.IssuerUrl,
			Role:            createdUserGroup.Role,
			Name:            createdUserGroup.Name,
			Notes:           createdUserGroup.Notes,
		}
	}

	limit := 200
	offset := 0
	count := 1
	var readUrl string
	var userGroupPage MspManagedUserGroupPage
	foundUsernames := mapset.NewSet[string]()

	for count > offset {
		client.Logger.Printf("Getting users from %d to %d\n", offset, offset+limit)
		readUrl = url.GetUserGroupsInMspManagedTenant(client.BaseUrl(), tenantUid, limit, offset)
		req := client.NewGet(ctx, readUrl)
		if err := req.Send(&userGroupPage); err != nil {
			return nil, err
		}
		for _, userGroup := range userGroupPage.Items {
			// add userGroup to map if not present
			if _, exists := readUserGroupDetailsMap[userGroup.GroupIdentifier]; exists {
				client.Logger.Printf("Updating user group information for %v\n", userGroup)
				readUserGroupDetailsMap[userGroup.GroupIdentifier] = userGroup
				foundUsernames.Add(userGroup.GroupIdentifier)
			}
		}

		offset += limit
		count = userGroupPage.Count
		client.Logger.Printf("Got %d user groups in tenant %s\n", count, tenantUid)
	}

	var readUserDetails []MspManagedUserGroup
	for _, value := range readUserGroupDetailsMap {
		// do not add in any users that were not found when we read from the API
		if foundUsernames.Contains(value.GroupIdentifier) {
			readUserDetails = append(readUserDetails, value)
		}
	}

	return &readUserDetails, nil
}
