package main

import (
	"fmt"
	"github.com/lao-tseu-is-alive/go-cloud-k8s-user-group/pkg/gohttpclient"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
)

type testScenario struct {
	name           string
	wantStatusCode int
	wantBody       string
	r              *http.Request
}

func TestMainExec(t *testing.T) {

	newRequest := func(method, url string, body string) *http.Request {
		r, err := http.NewRequest(method, ts.URL+url, strings.NewReader(body))
		if err != nil {
			t.Fatalf("### ERROR http.NewRequest %s on [%s] \n", method, url)
		}
		return r
	}
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

	newUserJson := "{\"username\":\"cgil\", \"name\":\"Carlos GIL\", \"email\":\"c@gil.town\", \"password_hash\":\"4acf0b39d9c4766709a3689f553ac01ab550545ffa4544dfc0b2cea82fba02a3\"}"
	createUserJson := `{
  "bad_password_count": 0,
  "create_time": "2022-08-31T19:12:37.19623Z",
  "creator": 1,
  "email": "c@gil.town",
  "id": 10,
  "is_active": true,
  "is_admin": false,
  "is_locked": false,
  "name": "Carlos GIL",
  "password_hash": "4acf0b39d9c4766709a3689f553ac01ab550545ffa4544dfc0b2cea82fba02a3",
  "username": "cgil"
}`

	tt := []testScenario{
		{
			name:           "GetUsers, should return all the existing users as json",
			wantStatusCode: http.StatusOK,
			wantBody:       string(jsonInitialData),
			r:              newRequest(http.MethodGet, "/api/users", ""),
		},
		{
			name:           "CreateUser with a valid new user, should return a valid User",
			wantStatusCode: http.StatusCreated,
			wantBody:       createUserJson,
			r:              newRequest(http.MethodPost, "/api/users", newUserJson),
		},
	}

	listenAddr := fmt.Sprintf("http://localhost:%d/", defaultPort)
	err := os.Setenv("PORT", fmt.Sprintf("%d", defaultPort))
	if err != nil {
		t.Errorf("Unable to set env variable PORT")
		return
	}
	// starting main in his own go routine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		main()
	}()
	gohttpclient.WaitForHttpServer(listenAddr, 1*time.Second, 10)

	resp, err := http.Get(listenAddr)
	if err != nil {
		t.Fatalf("Cannot make http get: %v\n", err)
	}
	defer resp.Body.Close()
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Should return an http status ok")

	receivedHtml, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v\n", err)
	}

	// check that receivedJson contains the specified tt.wantBody substring . https://pkg.go.dev/github.com/stretchr/testify/assert#Contains
	assert.Contains(t, string(receivedHtml), "<html", "Response should contain the html tag.")
	//assert.Contains(t, string(receivedHtml), "\"request_id\":", "Response should contain the request_id field.")

}
