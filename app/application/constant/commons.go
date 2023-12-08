package constant

const (
	ApplicationLayer = "Application"
)

type KindErrors string

var (
	ServiceFailed KindErrors = "SERVICE_FAILED"
)

type ReasonErrors string

var (
	SaveUserErrorDescription ReasonErrors = "Error saving user"
)
 