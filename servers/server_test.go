package servers

import (
	"bytes"
	"database/sql"
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
	server        *Server
	db            *gorm.DB
	initErr       error
	signingSecret = "hello-world"
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
		panic(initErr)
		// log.Fatal(initErr)
	}
	if db == nil {
		panic("db failed to initialize")
	}

	s, err := New(db)
	if err != nil {
		panic(err)
		//log.Fatal(err)
	}
	server = s
	server.signingSecret = signingSecret
	code = m.Run()
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/translations?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33061)/translations?charset=utf8mb4&parseTime=True&loc=Local"
	}

	// db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// if err != nil {
	// 	panic(err)
	// }
	sqlDB, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	return gormDB, err
}

func cleanup() {
	if db == nil {
		panic("DB is not initialized!!!")
	}
	if err := db.Exec("DELETE FROM users").Error; err != nil {
		panic(err)
		// log.Fatal(err)
	}
	if err := db.Exec("DELETE FROM translations").Error; err != nil {
		panic(err)
		// log.Fatal(err)
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
		req, err = http.NewRequest(method, path, bytes.NewBuffer(payload))
	}

	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	server.Router.ServeHTTP(w, req)
	require.NotNil(t, w)
	fmt.Println(w.Body.String())
	return w
}

func createToken(t *testing.T) string {
	t.Helper()
	td := core.TokenData{
		UserID:   uuid.NewString(),
		Email:    gofakeit.Email(),
		UserType: string(users.UserTypeAdmin),
	}
	token, err := td.Generate()
	require.NoError(t, err)
	require.NotEmpty(t, token)

	return token
}
