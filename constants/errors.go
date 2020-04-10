package constants

import "trackcoro/models"

const (
	NotAuthorizedCode = "NOT_AUTHORIZED"
	BadRequestCode    = "BAD_REQUEST"
	InternalErrorCode = "INTERNAL_ERROR"

	AdminNotExistsCode         = "ADMIN_NOT_EXISTS"
	SONotExistsCode            = "SO_NOT_EXISTS"
	SOAlreadyExistsCode        = "SO_ALREADY_EXISTS"
	SONotRegisteredByAdminCode = "SO_NOT_REGISTERED_BY_ADMIN"
	QuarantineNotExistsCode    = "QUARANTINE_NOT_EXISTS"

	UploadFileContentReadCode = "FILE_CONTENT_ERROR"
	UploadFileFailureCode     = "FILE_CONTENT_ERROR"
)

var (
	NotAuthorizedError = models.Error{Code: NotAuthorizedCode, Message: "You are not authorized to perform action"}
	BadRequestError    = models.Error{Code: BadRequestCode, Message: "Required field(s) are not provided"}
	InternalError      = models.Error{Code: InternalErrorCode, Message: "Could not perform action"}

	AdminNotExistsError         = models.Error{Code: AdminNotExistsCode, Message: "Admin does not exists"}
	SONotExistsError            = models.Error{Code: SONotExistsCode, Message: "SO does not exists"}
	SOAlreadyExistsError        = models.Error{Code: SOAlreadyExistsCode, Message: "SO already exists"}
	SONotRegisteredByAdminError = models.Error{Code: SONotRegisteredByAdminCode, Message: "SO is not registered by current admin"}
	QuarantineNotExistsError    = models.Error{Code: QuarantineNotExistsCode, Message: "Quarantine does not exists"}

	UploadFileContentReadError = models.Error{Code: UploadFileContentReadCode, Message: "Unable to read content of uploaded file"}
	UploadFileFailureError     = models.Error{Code: UploadFileFailureCode, Message: "Unable to upload file"}
)
