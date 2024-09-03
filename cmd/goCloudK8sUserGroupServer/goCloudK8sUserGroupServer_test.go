package main

import (
	"encoding/json"
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/config"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-common-libs/pkg/gohttpclient"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/crypto"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/users"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

/*

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
*/

const (
	DEBUG                           = true
	assertCorrectStatusCodeExpected = "expected status code should be returned"
)

type testStruct struct {
	name           string
	contentType    string
	wantStatusCode int
	wantBody       string
	paramKeyValues map[string]string
	httpMethod     string
	url            string
	body           string
}

// TestMainExec is instantiating the "real" main code using the env variable (in your .env files if you use the Makefile rule)
func TestMainExec(t *testing.T) {
	listenPort := config.GetPortFromEnvOrPanic(defaultPort)
	listenIP := config.GetListenIpFromEnvOrPanic("0.0.0.0")
	listenAddr := fmt.Sprintf("%s://%s:%d", "http", listenIP, listenPort)
	fmt.Printf("INFO: 'Will start HTTP server listening on port %s'\n", listenAddr)

	newRequest := func(method, url string, body string) *http.Request {
		fmt.Printf("INFO: 💥💥'newRequest %s on %s ##BODY : %+v'\n", method, url, body)
		r, err := http.NewRequest(method, url, strings.NewReader(body))
		if err != nil {
			t.Fatalf("### ERROR http.NewRequest %s on [%s] error is :%v\n", method, url, err)
		}
		if method == http.MethodPost {
			// r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			r.Header.Set("Content-Type", "application/json")
		}
		return r
	}

	adminUser := config.GetAdminUserFromEnvOrPanic("admin")
	adminPassword := config.GetAdminPasswordFromEnvOrPanic()
	passwordHash := crypto.Sha256Hash(adminPassword)
	uLogin := users.UserLogin{
		PasswordHash: passwordHash,
		Username:     adminUser,
	}
	bodyCreateUser, err := json.Marshal(uLogin)
	if err != nil {
		fmt.Println(err)
		return
	}

	tests := []testStruct{
		{
			name:           "1: Get on default get handler should contain html tag",
			wantStatusCode: http.StatusOK,
			contentType:    MIMEHtmlCharsetUTF8,
			wantBody:       "<html",
			paramKeyValues: make(map[string]string, 0),
			httpMethod:     http.MethodGet,
			url:            "/",
			body:           "",
		},
		{
			name:           "2: Post on default get handler should return an http error method not allowed ",
			wantStatusCode: http.StatusMethodNotAllowed,
			contentType:    MIMEHtmlCharsetUTF8,
			wantBody:       "Method Not Allowed",
			paramKeyValues: make(map[string]string, 0),
			httpMethod:     http.MethodPost,
			url:            "/",
			body:           `{"junk":"test with junk text"}`,
		},
		{
			name:           "3: Get on nonexistent route should return an http error not found ",
			wantStatusCode: http.StatusNotFound,
			contentType:    MIMEHtmlCharsetUTF8,
			wantBody:       "page not found",
			paramKeyValues: make(map[string]string, 0),
			httpMethod:     http.MethodGet,
			url:            "/aroutethatwillneverexisthere",
			body:           "",
		},
		{
			name:           "4: POST to login with valid credential should return a JWT token ",
			wantStatusCode: http.StatusOK,
			contentType:    MIMEAppJSONCharsetUTF8,
			wantBody:       "token",
			paramKeyValues: make(map[string]string, 0),
			httpMethod:     http.MethodPost,
			url:            "/login",
			body:           string(bodyCreateUser),
		},
	}

	// starting main in his own go routine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		main()
	}()
	gohttpclient.WaitForHttpServer(listenAddr, 2*time.Second, 10)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			r := newRequest(tt.httpMethod, listenAddr+tt.url, tt.body)
			//r.Header.Set(HeaderContentType, tt.contentType)
			resp, err := http.DefaultClient.Do(r)
			if DEBUG {
				fmt.Printf("### %s : %s on %s\n", tt.name, r.Method, r.URL)
			}
			if err != nil {
				fmt.Printf("### GOT ERROR : %s\n%+v", err, resp)
				t.Fatal(err)
			}
			defer resp.Body.Close()
			assert.Equal(t, tt.wantStatusCode, resp.StatusCode, assertCorrectStatusCodeExpected)
			receivedJson, _ := io.ReadAll(resp.Body)

			if DEBUG {
				fmt.Printf("WANTED   :%T - %#v\n", tt.wantBody, tt.wantBody)
				fmt.Printf("RECEIVED :%T - %#v\n", receivedJson, string(receivedJson))
			}
			// check that receivedJson contains the specified tt.wantBody substring . https://pkg.go.dev/github.com/stretchr/testify/assert#Contains
			assert.Contains(t, string(receivedJson), tt.wantBody, "Response should contain what was expected.")
		})
	}
}
