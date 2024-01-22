package client

import "errors"

var (
	ActivationHandlerErr = errors.New("not set activation handler")
	LicenseVerifyErr     = errors.New("License file verification error")
	LicenseExpirationErr = errors.New("License file expiration")

	OnlineUsersErr = errors.New("Exceeded maximum number of online users")
)
