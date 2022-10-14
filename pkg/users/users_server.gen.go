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
	// GroupList returns a list of groups
	// (GET /groups)
	GroupList(ctx echo.Context, params GroupListParams) error
	// GroupCreate will create a new group
	// (POST /groups)
	GroupCreate(ctx echo.Context) error
	// GroupDelete allows to delete a specific groupId
	// (DELETE /groups/{groupId})
	GroupDelete(ctx echo.Context, groupId int32) error
	// GroupGet will retrieve in backend all information about a specific groupId
	// (GET /groups/{groupId})
	GroupGet(ctx echo.Context, groupId int32) error
	// GroupUpdate allows to modify information about a specific groupId
	// (PUT /groups/{groupId})
	GroupUpdate(ctx echo.Context, groupId int32) error
	// UserList returns a list of users
	// (GET /users)
	UserList(ctx echo.Context, params UserListParams) error
	// UserCreate will create a new user
	// (POST /users)
	UserCreate(ctx echo.Context) error
	// UserDelete allows to delete a specific userId
	// (DELETE /users/{userId})
	UserDelete(ctx echo.Context, userId int32) error
	// UserGet will retrieve in backend all information about a specific userId
	// (GET /users/{userId})
	UserGet(ctx echo.Context, userId int32) error
	// UserUpdate allows to modify information about a specific userId
	// (PUT /users/{userId})
	UserUpdate(ctx echo.Context, userId int32) error
	// UserChangePassword allows a user to change it's own password
	// (PUT /users/{userId}/changepassword)
	UserChangePassword(ctx echo.Context, userId int32) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// GroupList converts echo context to params.
func (w *ServerInterfaceWrapper) GroupList(ctx echo.Context) error {
	var err error

	ctx.Set(JWTAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GroupListParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GroupList(ctx, params)
	return err
}

// GroupCreate converts echo context to params.
func (w *ServerInterfaceWrapper) GroupCreate(ctx echo.Context) error {
	var err error

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GroupCreate(ctx)
	return err
}

// GroupDelete converts echo context to params.
func (w *ServerInterfaceWrapper) GroupDelete(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "groupId" -------------
	var groupId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "groupId", runtime.ParamLocationPath, ctx.Param("groupId"), &groupId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter groupId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GroupDelete(ctx, groupId)
	return err
}

// GroupGet converts echo context to params.
func (w *ServerInterfaceWrapper) GroupGet(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "groupId" -------------
	var groupId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "groupId", runtime.ParamLocationPath, ctx.Param("groupId"), &groupId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter groupId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GroupGet(ctx, groupId)
	return err
}

// GroupUpdate converts echo context to params.
func (w *ServerInterfaceWrapper) GroupUpdate(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "groupId" -------------
	var groupId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "groupId", runtime.ParamLocationPath, ctx.Param("groupId"), &groupId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter groupId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GroupUpdate(ctx, groupId)
	return err
}

// UserList converts echo context to params.
func (w *ServerInterfaceWrapper) UserList(ctx echo.Context) error {
	var err error

	ctx.Set(JWTAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params UserListParams
	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserList(ctx, params)
	return err
}

// UserCreate converts echo context to params.
func (w *ServerInterfaceWrapper) UserCreate(ctx echo.Context) error {
	var err error

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserCreate(ctx)
	return err
}

// UserDelete converts echo context to params.
func (w *ServerInterfaceWrapper) UserDelete(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserDelete(ctx, userId)
	return err
}

// UserGet converts echo context to params.
func (w *ServerInterfaceWrapper) UserGet(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserGet(ctx, userId)
	return err
}

// UserUpdate converts echo context to params.
func (w *ServerInterfaceWrapper) UserUpdate(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserUpdate(ctx, userId)
	return err
}

// UserChangePassword converts echo context to params.
func (w *ServerInterfaceWrapper) UserChangePassword(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "userId" -------------
	var userId int32

	err = runtime.BindStyledParameterWithLocation("simple", false, "userId", runtime.ParamLocationPath, ctx.Param("userId"), &userId)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter userId: %s", err))
	}

	ctx.Set(JWTAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.UserChangePassword(ctx, userId)
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

	router.GET(baseURL+"/groups", wrapper.GroupList)
	router.POST(baseURL+"/groups", wrapper.GroupCreate)
	router.DELETE(baseURL+"/groups/:groupId", wrapper.GroupDelete)
	router.GET(baseURL+"/groups/:groupId", wrapper.GroupGet)
	router.PUT(baseURL+"/groups/:groupId", wrapper.GroupUpdate)
	router.GET(baseURL+"/users", wrapper.UserList)
	router.POST(baseURL+"/users", wrapper.UserCreate)
	router.DELETE(baseURL+"/users/:userId", wrapper.UserDelete)
	router.GET(baseURL+"/users/:userId", wrapper.UserGet)
	router.PUT(baseURL+"/users/:userId", wrapper.UserUpdate)
	router.PUT(baseURL+"/users/:userId/changepassword", wrapper.UserChangePassword)

}
