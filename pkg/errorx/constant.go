package errorx

const (
	DefaultBadRequestID            = "bad_request"              // 400
	DefaultUnauthorizedID          = "unauthorized"             // 401
	DefaultForbiddenID             = "forbidden"                // 403
	DefaultNotFoundID              = "not_found"                // 404
	DefaultMethodNotAllowedID      = "method_not_allowed"       // 405
	DefaultRequestTimeoutID        = "request_timeout"          // 408
	DefaultConflictID              = "conflict"                 // 409
	DefaultRequestEntityTooLargeID = "request_entity_too_large" // 413
	DefaultTooManyRequestsID       = "too_many_requests"        // 429
	DefaultInternalServerErrorID   = "internal_server_error"    // 500
)
