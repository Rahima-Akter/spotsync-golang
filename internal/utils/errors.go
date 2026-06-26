package utils

import "errors"

var (

	// 404 Not Found
	ErrNotFound = errors.New("resource not found")

	// 400 Bad Request
	ErrDuplicateEmail = errors.New("email already exists")

	// 401 Unauthorized
	ErrInvalidCredentials = errors.New("invalid email or password")

	//401 Unauthorized
	ErrUnauthorized = errors.New("unauthorized access")

	// 403 Forbidden
	ErrForbidden = errors.New("forbidden: insufficient permissions")

	//  409 Conflict
	ErrZoneFull = errors.New("parking zone is at full capacity")

	// 409 Conflict
	ErrDuplicateReservation = errors.New("license plate already has an active reservation")
)
