// Package users provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package users

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/labstack/echo/v4"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// GetUsers returns a list of users
	// (GET /api/users)
	GetUsers(ctx echo.Context, params GetUsersParams) error
	// CreateUser will create a new user
	// (POST /api/users)
	CreateUser(ctx echo.Context) error
	// DeleteUser allows to delete a specific userId
	// (DELETE /api/users/{userId})
	DeleteUser(ctx echo.Context, userId int32) error
	// GetUser will retrieve in backend all information about a specific userId
	// (GET /api/users/{userId})
	GetUser(ctx echo.Context, userId int32) error
	// UpdateUser allows to modifiy information about a specific userId
	// (PUT /api/users/{userId})
	UpdateUser(ctx echo.Context, userId int32) error
	// ChangeUserPassword allows a user to change it's own password
	// (PUT /api/users/{userId}/changepassword)
	ChangeUserPassword(ctx echo.Context, userId int32) error
	// GetLogin allows client to do a preflight prepare for a login
	// (GET /login)
	GetLogin(ctx echo.Context) error
	// LoginUser allows client to try to authenticate, and then receive a valid JWT
	// (POST /login)
	LoginUser(ctx echo.Context) error
	// GetResetPasswordEmail allows a client to do a password reset and receive a new password link to his email
	// (GET /resetpassword)
	GetResetPasswordEmail(ctx echo.Context) error
	// SendResetPassword will send an email with a reset password url valid for one hour
	// (POST /resetpassword)
	SendResetPassword(ctx echo.Context) error
	// GetResetPasswordToken allows a client to do a password reset and receive a new password link to his email
	// (GET /resetpassword/{resetPasswordToken})
	GetResetPasswordToken(ctx echo.Context, resetPasswordToken string) error
	// ResetPassword will change password if resetPasswordToken in url is still valid
	// (POST /resetpassword/{resetPasswordToken})
	ResetPassword(ctx echo.Context, resetPasswordToken string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GetUsers converts echo context to params.
func (w *ServerInterfaceWrapper) GetUsers(ctx echo.Context) error {
	var err error

	ctx.Set(JWTAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetUsersParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUsers(ctx, params)
	return err
}

// CreateUser converts echo context to params.
func (w *ServerInterfaceWrapper) CreateUser(ctx echo.Context) error {
	var err error

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateUser(ctx)
	return err
}

// DeleteUser converts echo context to params.
func (w *ServerInterfaceWrapper) DeleteUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.DeleteUser(ctx, userId)
	return err
}

// GetUser converts echo context to params.
func (w *ServerInterfaceWrapper) GetUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetUser(ctx, userId)
	return err
}

// UpdateUser converts echo context to params.
func (w *ServerInterfaceWrapper) UpdateUser(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UpdateUser(ctx, userId)
	return err
}

// ChangeUserPassword converts echo context to params.
func (w *ServerInterfaceWrapper) ChangeUserPassword(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ChangeUserPassword(ctx, userId)
	return err
}

// GetLogin converts echo context to params.
func (w *ServerInterfaceWrapper) GetLogin(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetLogin(ctx)
	return err
}

// LoginUser converts echo context to params.
func (w *ServerInterfaceWrapper) LoginUser(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.LoginUser(ctx)
	return err
}

// GetResetPasswordEmail converts echo context to params.
func (w *ServerInterfaceWrapper) GetResetPasswordEmail(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetResetPasswordEmail(ctx)
	return err
}

// SendResetPassword converts echo context to params.
func (w *ServerInterfaceWrapper) SendResetPassword(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.SendResetPassword(ctx)
	return err
}

// GetResetPasswordToken converts echo context to params.
func (w *ServerInterfaceWrapper) GetResetPasswordToken(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "resetPasswordToken" -------------
	var resetPasswordToken string

	err = runtime.BindStyledParameterWithLocation("simple", false, "resetPasswordToken", runtime.ParamLocationPath, ctx.Param("resetPasswordToken"), &resetPasswordToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resetPasswordToken: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetResetPasswordToken(ctx, resetPasswordToken)
	return err
}

// ResetPassword converts echo context to params.
func (w *ServerInterfaceWrapper) ResetPassword(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "resetPasswordToken" -------------
	var resetPasswordToken string

	err = runtime.BindStyledParameterWithLocation("simple", false, "resetPasswordToken", runtime.ParamLocationPath, ctx.Param("resetPasswordToken"), &resetPasswordToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter resetPasswordToken: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ResetPassword(ctx, resetPasswordToken)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.GET(baseURL+"/api/users", wrapper.GetUsers)
	router.POST(baseURL+"/api/users", wrapper.CreateUser)
	router.DELETE(baseURL+"/api/users/:userId", wrapper.DeleteUser)
	router.GET(baseURL+"/api/users/:userId", wrapper.GetUser)
	router.PUT(baseURL+"/api/users/:userId", wrapper.UpdateUser)
	router.PUT(baseURL+"/api/users/:userId/changepassword", wrapper.ChangeUserPassword)
	router.GET(baseURL+"/login", wrapper.GetLogin)
	router.POST(baseURL+"/login", wrapper.LoginUser)
	router.GET(baseURL+"/resetpassword", wrapper.GetResetPasswordEmail)
	router.POST(baseURL+"/resetpassword", wrapper.SendResetPassword)
	router.GET(baseURL+"/resetpassword/:resetPasswordToken", wrapper.GetResetPasswordToken)
	router.POST(baseURL+"/resetpassword/:resetPasswordToken", wrapper.ResetPassword)

}
