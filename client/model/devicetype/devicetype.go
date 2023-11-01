package devicetype

type Type string

const (
	Asa        Type = "ASA"
	Ios        Type = "IOS"
	CloudFmc   Type = "FMCE"
	CloudFtd   Type = "FTDC"
	GenericSSH Type = "GENERIC_SSH"
)
