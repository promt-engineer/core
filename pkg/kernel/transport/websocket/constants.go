package websocket

const (
	StatusInternalError = iota
	StatusSuccess
	StatusBadRequest
	StatusUnauthorized
	StatusActionNotFound
	StatusConflict
	StatusSessionExpired
	StatusValidationFailed
	StatusForbidden
	StatusPaymentRequirement
)
