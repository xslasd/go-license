package client

import "errors"

var (
	ActivationHandlerErr = errors.New("not set activation handler")
	LicenseVerifyErr     = errors.New("License file verification error")
	LicenseExpirationErr = errors.New("License file expiration")

	ActivationChecksValErr = errors.New("Activation check value error")

	OnlineUsersErr = errors.New("Exceeded maximum number of online users")
)
