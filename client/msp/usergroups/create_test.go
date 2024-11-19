package usergroups_test

import (
	"context"
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/http"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactionstatus"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/internal/publicapi/transaction/transactiontype"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/msp/usergroups"
	"github.com/google/uuid"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	netHttp "net/http"
	"sort"
	"strconv"
	"testing"
	"time"
)

// Function to generate user groups
func generateUserGroups(num int) []usergroups.MspManagedUserGroup {
	var createdUserGroups []usergroups.MspManagedUserGroup
	for i := 1; i <= num; i++ {
		uid := "uid" + strconv.Itoa(i) // Generate unique UID
		var role string
		if i%2 == 0 {
			role = "ROLE_SUPER_ADMIN"
		} else {
			role = "ROLE_ADMIN"
		}
		var notes string
		if i%2 == 0 {
			notes = "notes" + strconv.Itoa(i)
		}

		createdUserGroups = append(createdUserGroups, usergroups.MspManagedUserGroup{
			Uid:             uid,
			GroupIdentifier: "groupIdentifier" + strconv.Itoa(i),
			Name:            "name" + strconv.Itoa(i),
			Role:            role,
			Notes:           &notes,
		})
	}
	return createdUserGroups
}

func TestCreate(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("successfully create user groups in MSP-managed tenant", func(t *testing.T) {
		httpmock.Reset()
		var managedTenantUid = uuid.New().String()
		var notes = "This is a group of developers"
		var createInp = []usergroups.MspManagedUserGroupInput{
			{
				GroupIdentifier: "developers",
				IssuerUrl:       "https://okta.com/123456",
				Name:            "Developers",
				Role:            "ROLE_ADMIN",
				Notes:           &notes,
			},
			{
				GroupIdentifier: "managers",
				IssuerUrl:       "https://okta.com/123456",
				Name:            "Managers",
				Role:            "ROLE_READ_ONLY",
			},
		}
		var userGroupsInCdoTenant = generateUserGroups(250)
		var userGroupsWithIds []usergroups.MspManagedUserGroup
		for _, userGroup := range createInp {
			userGroupWithId := usergroups.MspManagedUserGroup{
				Uid:             uuid.New().String(),
				GroupIdentifier: userGroup.GroupIdentifier,
				IssuerUrl:       userGroup.IssuerUrl,
				Name:            userGroup.Name,
				Role:            userGroup.Role,
				Notes:           userGroup.Notes,
			}
			userGroupsInCdoTenant = append(userGroupsInCdoTenant, userGroupWithId)
			userGroupsWithIds = append(userGroupsWithIds, userGroupWithId)
		}
		firstUserGroupPage := usergroups.MspManagedUserGroupPage{Items: userGroupsInCdoTenant[:200], Count: len(userGroupsInCdoTenant), Limit: 200, Offset: 0}
		secondUserGroupPage := usergroups.MspManagedUserGroupPage{Items: userGroupsInCdoTenant[200:], Count: len(userGroupsInCdoTenant), Limit: 200, Offset: 200}
		var transactionUid = uuid.New().String()
		var inProgressTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USER_GROUPS_TO_TENANT,
			Status:          transactionstatus.IN_PROGRESS,
		}
		var doneTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USER_GROUPS_TO_TENANT,
			Status:          transactionstatus.DONE,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users/groups", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, inProgressTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			inProgressTransaction.PollingUrl,
			httpmock.NewJsonResponderOrPanic(200, doneTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users/groups?limit=200&offset=0", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, firstUserGroupPage),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users/groups?limit=200&offset=200", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, secondUserGroupPage),
		)

		actual, err := usergroups.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid token", 0, 0, time.Minute), managedTenantUid, &createInp)

		assert.NotNil(t, actual, "Created user groups should have not been nil")
		assert.Nil(t, err, "Created user groups operation should have not been an error")
		sort.Slice(userGroupsWithIds, func(i, j int) bool {
			return userGroupsWithIds[i].Uid < userGroupsWithIds[j].Uid
		})
		sort.Slice(*actual, func(i, j int) bool {
			return (*actual)[i].Uid < (*actual)[j].Uid
		})
		assert.Equal(t, userGroupsWithIds, *actual, "Created users operation should have been the same as the created tenant")
	})

	t.Run("user group creation transaction fails", func(t *testing.T) {
		httpmock.Reset()
		var managedTenantUid = uuid.New().String()
		var notes = "This is a group of developers"
		var createInp = []usergroups.MspManagedUserGroupInput{
			{
				GroupIdentifier: "developers",
				IssuerUrl:       "https://okta.com/123456",
				Name:            "Developers",
				Role:            "ROLE_ADMIN",
				Notes:           &notes,
			},
			{
				GroupIdentifier: "managers",
				IssuerUrl:       "https://okta.com/123456",
				Name:            "Managers",
				Role:            "ROLE_READ_ONLY",
			},
		}
		var transactionUid = uuid.New().String()
		var inProgressTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USERS_TO_TENANT,
			Status:          transactionstatus.IN_PROGRESS,
		}
		var errorTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USERS_TO_TENANT,
			Status:          transactionstatus.ERROR,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users/groups", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, inProgressTransaction),
		)
		httpmock.RegisterResponder(
			netHttp.MethodGet,
			inProgressTransaction.PollingUrl,
			httpmock.NewJsonResponderOrPanic(200, errorTransaction),
		)

		actual, err := usergroups.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid token", 0, 0, time.Minute), managedTenantUid, &createInp)

		assert.Nil(t, actual, "Created user groups should be nil")
		assert.NotNil(t, err, "Created user groups in tenant operation should have an error")
		assert.Equal(t, usergroups.CreateError{
			Err:               publicapi.NewTransactionErrorFromTransaction(errorTransaction),
			CreatedResourceId: &managedTenantUid,
		}, *err, "created transaction error does not match")
	})

	t.Run("user group creation API call fails with an error transaction", func(t *testing.T) {
		httpmock.Reset()
		var managedTenantUid = uuid.New().String()
		var notes = "This is a group of developers"
		var createInp = []usergroups.MspManagedUserGroupInput{
			{
				GroupIdentifier: "developers",
				IssuerUrl:       "https://okta.com/123456",
				Name:            "Developers",
				Role:            "ROLE_ADMIN",
				Notes:           &notes,
			},
			{
				GroupIdentifier: "managers",
				IssuerUrl:       "https://okta.com/123456",
				Name:            "Managers",
				Role:            "ROLE_READ_ONLY",
			},
		}
		var transactionUid = uuid.New().String()
		var errorTransaction = transaction.Type{
			TransactionUid:  transactionUid,
			TenantUid:       uuid.New().String(),
			EntityUid:       managedTenantUid,
			EntityUrl:       "https://unittest.cdo.cisco.com/api/rest/v1/msp/tenants/" + managedTenantUid,
			PollingUrl:      "https://unittest.cdo.cisco.com/api/rest/v1/transactions/" + transactionUid,
			SubmissionTime:  "2024-09-10T20:10:00Z",
			LastUpdatedTime: "2024-10-10T20:10:00Z",
			Type:            transactiontype.MSP_ADD_USERS_TO_TENANT,
			Status:          transactionstatus.ERROR,
		}

		httpmock.RegisterResponder(
			netHttp.MethodPost,
			fmt.Sprintf("/api/rest/v1/msp/tenants/%s/users/groups", managedTenantUid),
			httpmock.NewJsonResponderOrPanic(200, errorTransaction),
		)

		actual, err := usergroups.Create(context.Background(), *http.MustNewWithConfig(baseUrl, "valid token", 0, 0, time.Minute), managedTenantUid, &createInp)

		assert.Nil(t, actual, "Created user groups should be nil")
		assert.NotNil(t, err, "Created user groups in tenant operation should have an error")
		var emptyCreatedResourceId = ""
		assert.Equal(t, usergroups.CreateError{
			Err:               publicapi.NewTransactionErrorFromTransaction(errorTransaction),
			CreatedResourceId: &emptyCreatedResourceId,
		}, *err, "created transaction error does not match")
	})
}
