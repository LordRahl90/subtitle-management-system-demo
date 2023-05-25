package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"translations/domains/tms"
	"translations/requests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTranslation(t *testing.T) {
	token := createToken(t)
	req := requests.Translation{
		SourceLanguage: "en",
		TargetLanguage: "de",
		Source:         "Hello World",
		Target:         "Hallo Welt",
	}
	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotEmpty(t, b)

	fmt.Printf("\n\nReq: %s\n\nToken: %s\n\n", b, token)

	res := requestHelper(t, http.MethodPost, "/tms", token, b)
	require.Equal(t, http.StatusCreated, res.Code)
}

func TestTranslateSentence(t *testing.T) {
	ctx := context.Background()
	token := createToken(t)

	err := server.translateService.Create(ctx, &tms.Translation{
		SourceLanguage: "en",
		TargetLanguage: "de",
		Source:         "Hello World",
		Target:         "Hallo Welt",
	})
	require.NoError(t, err)

	req := requests.Translation{
		SourceLanguage: "en",
		TargetLanguage: "de",
		Source:         "Hello World",
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	res := requestHelper(t, http.MethodPost, "/tms/translate", token, b)
	require.Equal(t, http.StatusOK, res.Code)
	exp := `{"target":"Hallo Welt"}`
	assert.Equal(t, exp, res.Body.String())
}

func TestUploadTranslation(t *testing.T) {
	token := createToken(t)
	fileName := "./testdata/translation.json"
	var b bytes.Buffer
	w := httptest.NewRecorder()

	writer := multipart.NewWriter(&b)

	file, err := os.Open(fileName)
	require.NoError(t, err)
	defer file.Close()

	form, err := writer.CreateFormFile("file", fileName)
	require.NoError(t, err)

	_, err = io.Copy(form, file)
	require.NoError(t, err)
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "/tms/upload", &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	server.Router.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
	exp := `"translations uploaded successfully"`
	assert.Equal(t, exp, w.Body.String())
}
