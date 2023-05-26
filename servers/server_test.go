package servers

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"translations/domains/core"
	"translations/domains/users"

	"github.com/brianvoe/gofakeit"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	server          *Server
	db              *gorm.DB
	initErr         error
	signingSecret   = "hello-world"
	outputDirectory = "./testdata/outputs/"
)

func TestMain(m *testing.M) {
	var code = 1
	defer func() {
		if db == nil {
			log.Fatal("db not initialized")
		}
		cleanup()
		os.Exit(code)
	}()
	db, initErr = setupTestDB()
	if initErr != nil {
		log.Fatal(initErr)
	}

	s, err := New(db, signingSecret, outputDirectory)
	if err != nil {
		log.Fatal(err)
	}
	server = s
	server.signingSecret = signingSecret
	code = m.Run()
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/translations?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/translations?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func cleanup() {
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		log.Fatal(err)
	}
	if err := db.Exec("DELETE FROM translations").Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Exec("DELETE FROM contents").Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Exec("DELETE FROM subtitles").Error; err != nil {
		log.Fatal(err)
	}
}

func requestHelper(t *testing.T, method, path, token string, payload []byte) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	var (
		req *http.Request
		err error
	)

	if len(payload) == 0 {
		req, err = http.NewRequest(method, path, nil)
	} else {
		fmt.Printf("\n\nReq: %s\n\n", payload)
		req, err = http.NewRequest(method, path, bytes.NewBuffer(payload))
	}

	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		fmt.Printf("\n\nToken: %s\n\n", token)
		req.Header.Set("Authorization", "Bearer "+token)
	}
	server.Router.ServeHTTP(w, req)
	require.NotNil(t, w)
	fmt.Printf("\n\nResponse: %s\n\n", w.Body.String())
	return w
}

func createToken(t *testing.T) string {
	t.Helper()
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate(server.signingSecret)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}
