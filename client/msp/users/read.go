package users

import (
	"context"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/url"
	mapset "github.com/deckarep/golang-set/v2"
)

func ReadCreatedUsersInTenant(ctx context.Context, client http.Client, readInput MspUsersInput) (*[]UserDetails, error) {
	client.Logger.Printf("Reading users in tenant %s\n", readInput.TenantUid)

	// create a map of the users that were created
	// find the list of deleted users by removing from the list every time a user is found in the response
	readUserDetailsMap := map[string]UserDetails{}
	for _, createdUser := range readInput.Users {
		readUserDetailsMap[createdUser.Username] = createdUser
	}

	limit := 200
	offset := 0
	count := 1
	var readUrl string
	var userPage UserPage
	foundUsernames := mapset.NewSet[string]()

	for count > offset {
		client.Logger.Printf("Getting users from %d to %d\n", offset, offset+limit)
		readUrl = url.GetUsersInMspManagedTenant(client.BaseUrl(), readInput.TenantUid, limit, offset)
		req := client.NewGet(ctx, readUrl)
		if err := req.Send(&userPage); err != nil {
			return nil, err
		}
		for _, user := range userPage.Items {
			// add user to map if not present
			if _, exists := readUserDetailsMap[user.Username]; exists {
				client.Logger.Printf("Updating user information for %v\n", user)
				readUserDetailsMap[user.Username] = user
				foundUsernames.Add(user.Username)
			}
		}

		offset += limit
		count = userPage.Count
		client.Logger.Printf("Got %d users in tenant %s\n", count, readInput.TenantUid)
	}

	var readUserDetails []UserDetails
	for _, value := range readUserDetailsMap {
		// do not add in any users that were not found when we read from the API
		if foundUsernames.Contains(value.Username) {
			readUserDetails = append(readUserDetails, value)
		}
	}

	return &readUserDetails, nil
}
