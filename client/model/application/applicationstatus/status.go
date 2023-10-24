package applicationstatus

type Type string

// see: https://github.com/cisco-lockhart/aegis/blob/master/modules/configs/src/main/java/com/cisco/lockhart/applications/models/ApplicationStatus.java
const (
	Init        Type = "INIT"        // used to send sns/sqs message for ops ticket
	Requested   Type = "REQUESTED"   // the cdFMC is being provisioned
	Active      Type = "ACTIVE"      // application is ready for use
	Unreachable Type = "UNREACHABLE" // when heartbeat fails
)
