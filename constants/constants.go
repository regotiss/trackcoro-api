package constants

var (
	DBConnectionString = "DB_CONNECTION"
	AdminMobileNumber  = "ADMIN_MOBILE_NUMBER"
)

var (
	Empty              = ""
	Authorization      = "Authorization"
	MobileNumber       = "MobileNumber"
	DetailsTimeFormat  = "2006-01-02"
)

var (
	AdminRole      = "admin"
	QuarantineRole = "quarantine"
	SORole         = "supervising_officer"
)
var (
	QuarantineNotExistsError         = "quarantine does not exists"
	TimeParseError                   = "could not parse to time format"
	EnvVariableNotFoundError         = "environment variable not found"
	QuarantineNotRegisteredBySOError = "quarantine is not registered by current so"
)
