package constants

import "trackcoro/models"

const (
	NotAuthorizedCode = "NOT_AUTHORIZED"
	BadRequestCode    = "BAD_REQUEST"
	InternalErrorCode = "INTERNAL_ERROR"

	AdminNotExistsCode              = "ADMIN_NOT_EXISTS"
	SONotExistsCode                 = "SO_NOT_EXISTS"
	SOAlreadyExistsCode             = "SO_ALREADY_EXISTS"
	SONotRegisteredByAdminCode      = "SO_NOT_REGISTERED_BY_ADMIN"
	QuarantineNotExistsCode         = "QUARANTINE_NOT_EXISTS"
	QuarantineAlreadyExistsCode     = "QUARANTINE_ALREADY_EXISTS"
	QuarantineNotRegisteredBySOCode = "QUARANTINE_NOT_REGISTERED_BY_SO"

	UploadFileContentReadCode = "FILE_CONTENT_ERROR"
	UploadFileFailureCode     = "FILE_UPLOAD_FAILED"

	DOBIncorrectFormatCode            = "DOB_INCORRECT_DATE_FORMAT"
	QuarantineDateIncorrectFormatCode = "QUARANTINE_INCORRECT_DATE_FORMAT"
	TravelDateIncorrectFormatCode     = "TRAVEL_INCORRECT_DATE_FORMAT"

	SendNotificationFailedCode = "SEND_NOTIFICATION_FAILED"
)

var (
	NotAuthorizedError = models.Error{Code: NotAuthorizedCode, Message: "You are not authorized to perform action"}
	BadRequestError    = models.Error{Code: BadRequestCode, Message: "Required field(s) are not provided"}
	InternalError      = models.Error{Code: InternalErrorCode, Message: "Could not perform action"}

	AdminNotExistsError              = models.Error{Code: AdminNotExistsCode, Message: "Admin does not exists"}
	SONotExistsError                 = models.Error{Code: SONotExistsCode, Message: "SO does not exists"}
	SOAlreadyExistsError             = models.Error{Code: SOAlreadyExistsCode, Message: "SO already exists"}
	SONotRegisteredByAdminError      = models.Error{Code: SONotRegisteredByAdminCode, Message: "SO is not registered by current admin"}
	QuarantineNotExistsError         = models.Error{Code: QuarantineNotExistsCode, Message: "Quarantine does not exists"}
	QuarantineAlreadyExistsError     = models.Error{Code: QuarantineAlreadyExistsCode, Message: "Quarantine already exists"}
	QuarantineNotRegisteredBySOError = models.Error{Code: QuarantineNotRegisteredBySOCode, Message: "Quarantine is not registered by current SO"}

	UploadFileContentReadError = models.Error{Code: UploadFileContentReadCode, Message: "Unable to read content of uploaded file"}
	UploadFileFailureError     = models.Error{Code: UploadFileFailureCode, Message: "Unable to upload file"}

	DOBIncorrectFormatError            = models.Error{Code: DOBIncorrectFormatCode, Message: "DOB is not in correct format"}
	QuarantineDateIncorrectFormatError = models.Error{Code: QuarantineDateIncorrectFormatCode, Message: "Quarantine started date is in not in correct format"}
	TravelDateIncorrectFormatError     = models.Error{Code: TravelDateIncorrectFormatCode, Message: "Travel date is in not in correct format"}

	SendNotificationFailedError = models.Error{Code: SendNotificationFailedCode, Message: "Could not send notification"}
)
