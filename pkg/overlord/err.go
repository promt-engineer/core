package overlord

import (
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	WrongFreeSpinIDMessage  = "user have not such free bet"
	UserBlockedMessage      = "user block"
	LowBalanceMessage       = "too low balance"
	IntegratorIsDeadMessage = "integrator is dead"
)

type ValidationError struct {
	Message string
}

func (err ValidationError) Error() string {
	return err.Message
}

const (
	CodeSessionExpired = 100
)

var (
	ErrBalanceTooLow             = errors.New("not enough money")
	ErrWrongSessionToken         = errors.New("wrong session token")
	ErrWrongFreeSpinID           = errors.New("wrong free spin id")
	ErrMarshaling                = errors.New("bad content for marshaling")
	ErrSessionTokenExpired       = errors.New("session token expired ")
	ErrUserHasDifferentCurrency  = errors.New("user has different currency")
	ErrUserIsBlocked             = errors.New("user is blocked")
	ErrIntegratorCriticalFailure = errors.New("integrator critical failure")
)

func mapError(err error) error {
	if err != nil {
		code, _ := status.FromError(err)
		switch code.Code() {
		case codes.PermissionDenied:
			return ErrBalanceTooLow
		case codes.Unauthenticated:
			return ErrWrongSessionToken
		case codes.InvalidArgument:
			if code.Message() == WrongFreeSpinIDMessage {
				return ErrWrongFreeSpinID
			}

			return ValidationError{Message: code.Message()}
		case CodeSessionExpired:
			return ErrSessionTokenExpired
		case codes.Aborted:
			if code.Message() == UserBlockedMessage {
				return ErrUserIsBlocked
			}

			if code.Message() == IntegratorIsDeadMessage {
				return ErrIntegratorCriticalFailure
			}

			return err
		case codes.AlreadyExists, codes.Canceled, codes.DataLoss,
			codes.DeadlineExceeded, codes.FailedPrecondition, codes.NotFound,
			codes.OK, codes.OutOfRange, codes.ResourceExhausted, codes.Internal,
			codes.Unavailable, codes.Unimplemented, codes.Unknown:
			fallthrough
		default:
			return err
		}
	}

	return nil
}
