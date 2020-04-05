package constants

var (
	DBConnectionString = "DB_CONNECTION"
	AdminMobileNumber = "ADMIN_MOBILE_NUMBER"
)

var (
	Empty = ""
	Authorization = "Authorization"
	MobileNumber = "MobileNumber"
	DetailsTimeFormat = "2006-01-02"
	FileNameTimeFormat = "2006-01-02T15:04:05"
)

var (
	AdminRole      = "admin"
	QuarantineRole = "quarantine"
)
var (
	NotExists = "quarantine does not exists"
	TimeParseError = "could not parse to time format"
	EnvVariableNotFound = "environment variable not found"
)