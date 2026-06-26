package utils

import "github.com/labstack/echo/v4"

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`   // omitempty means this field won't appear if empty
	Errors  interface{} `json:"errors,omitempty"` // For validation errors or error details
}

func SuccessResponse(c echo.Context, statusCode int, message string, data interface{}) error {
	return c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, statusCode int, message string, errors interface{}) error {
	return c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Errors:  errors,
	})
}
