package testing

import (
	"fmt"
	"github.com/CiscoDevnet/terraform-provider-cdo/go-client/model/ftd/tier"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

type Model struct {
	TransactionSubmissionTime time.Time
	TransactionUid            uuid.UUID

	CdgUid  uuid.UUID
	CdgName string

	SdcUid uuid.UUID

	AsaUid         uuid.UUID
	AsaName        string
	AsaCreatedDate time.Time
	AsaHost        string
	AsaPort        string
	AsaUsername    string
	AsaPassword    string

	FtdUid              uuid.UUID
	FtdName             string
	FtdAccessPolicyName string
	FtdAccessPolicyUid  uuid.UUID
	FtdPerformanceTier  tier.Type

	DuoAdminPanelUid    uuid.UUID
	DuoAdminPanelName   string
	DuoAdminPanelHost   string
	DuoAdminPanelLabels []string

	IosUid      uuid.UUID
	IosName     string
	IosUsername string
	IosPassword string
	IosHost     string
	IosPort     string

	FmcHost       string
	FmcDomainUuid uuid.UUID

	TenantUid uuid.UUID
	BaseUrl   string
}

func NewRandomModel() Model {
	return Model{
		TransactionSubmissionTime: time.Now(),
		TransactionUid:            uuid.New(),

		CdgUid:  uuid.New(),
		CdgName: randString("cdg"),

		SdcUid: uuid.New(),

		AsaUid:         uuid.New(),
		AsaName:        randString("asa"),
		AsaCreatedDate: time.Now().Add(-3 * time.Second),
		AsaHost:        randHost(),
		AsaPort:        "443",
		AsaUsername:    randString("asa-username"),
		AsaPassword:    randString("asa-password"),

		FtdUid:              uuid.New(),
		FtdName:             randString("ftd"),
		FtdAccessPolicyName: randString("ftd-access-policy"),
		FtdAccessPolicyUid:  uuid.New(),
		FtdPerformanceTier:  tier.FTDv5,

		DuoAdminPanelUid:    uuid.New(),
		DuoAdminPanelName:   randString("duo-admin-panel"),
		DuoAdminPanelHost:   randHost(),
		DuoAdminPanelLabels: []string{"lab1", "lab2", "lab3"},

		IosUid:      uuid.New(),
		IosName:     randString("ios"),
		IosUsername: randString("ios-username"),
		IosPassword: randString("ios-password"),
		IosHost:     randHost(),
		IosPort:     "443",

		FmcHost:       randHost(),
		FmcDomainUuid: uuid.New(),

		TenantUid: uuid.New(),
		BaseUrl:   "https://unit-test.cisco.com",
	}
}

func randString(prefix string) string {
	return fmt.Sprintf("test-%s-%d", prefix, randInt())
}

func randInt() int {
	return rand.Intn(1000)
}

func randHost() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Intn(32), rand.Intn(32), rand.Intn(32), rand.Intn(32))
}
