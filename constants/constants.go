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
	FileNameTimeFormat = "2006-01-02T15:04:05"
)

var (
	AdminRole      = "admin"
	QuarantineRole = "quarantine"
	SORole         = "supervising_officer"
)
var (
	QuarantineNotExistsError         = "quarantine does not exists"
	SONotExistsError                 = "so does not exists"
	AdminNotExistsError              = "admin does not exists"
	TimeParseError                   = "could not parse to time format"
	EnvVariableNotFoundError         = "environment variable not found"
	SONotRegisteredByAdminError      = "so is not registered by current admin"
	QuarantineNotRegisteredBySOError = "quarantine is not registered by current so"
)
