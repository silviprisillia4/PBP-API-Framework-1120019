package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
)

func sendUserSuccessResponse(c echo.Context, users []User, message string) error {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	var response UsersResponse
	response.Status = 200
	response.Message = message
	return json.NewEncoder(c.Response()).Encode(response)
}

func sendLoginLogoutSuccessResponse(c echo.Context, login bool) error {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	var response UsersResponse
	response.Status = 200
	if login {
		response.Message = "Welcome!"

	} else {
		response.Message = "Good bye!"
	}
	return json.NewEncoder(c.Response()).Encode(response)
}

func sendBadRequestResponse(c echo.Context, errorMessage string) error {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusBadRequest)
	var response ErrorResponse
	response.Status = 400
	response.Message = errorMessage
	return json.NewEncoder(c.Response()).Encode(response)
}

func sendUnauthorizedResponse(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusUnauthorized)
	var response ErrorResponse
	response.Status = 401
	response.Message = "Unauthorized Access"
	return json.NewEncoder(c.Response()).Encode(response)
}

func sendNotFoundResponse(c echo.Context, errorMessage string) error {
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusNotFound)
	var response ErrorResponse
	response.Status = 404
	response.Message = errorMessage
	return json.NewEncoder(c.Response()).Encode(response)
}
