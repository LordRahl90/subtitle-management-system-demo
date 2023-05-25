package servers

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"translations/domains/tms"
	"translations/requests"
	"translations/responses"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateSubtitle_WithNoTranslationUploaded(t *testing.T) {
	req := requests.Subtitle{
		Name:           gofakeit.CarMaker(),
		Filename:       strings.ToLower(gofakeit.BuzzWord()),
		SourceLanguage: "en",
		TargetLanguage: "de",
		Content: []requests.Content{
			{
				TimeRange: "00:00:12.00 - 00:01:20.00",
				Content:   "I am Arwen - I've come to help you.",
			},
			{
				TimeRange: "00:03:55.00 - 00:04:20.00",
				Content:   "Come back to the light.",
			},
		},
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	token := createToken(t)

	res := requestHelper(t, http.MethodPost, "/sts", token, b)
	require.Equal(t, http.StatusCreated, res.Code)

	var result responses.Subtitle
	err = json.Unmarshal(res.Body.Bytes(), &result)
	require.NoError(t, err)
	require.NotEmpty(t, result)

	assert.Len(t, result.Content, 2)
}

func TestCreateSubtitle_WithTranslationUploaded(t *testing.T) {
	ctx := context.Background()
	// upload some translations
	tr := []*tms.Translation{
		{
			SourceLanguage: "en",
			TargetLanguage: "de",
			Source:         "I am Arwen - I've come to help you.",
			Target:         "Ich bin Arwen - Ich bin gekommen, um dir zu helfen.",
		},
		{
			SourceLanguage: "en",
			TargetLanguage: "de",
			Source:         "Come back to the light.",
			Target:         "Komm zurück zum Licht.",
		},
	}

	for _, v := range tr {
		require.NoError(t, server.translateService.Create(ctx, v))
	}

	req := requests.Subtitle{
		Name:           gofakeit.CarMaker(),
		Filename:       strings.ToLower(gofakeit.BuzzWord()),
		SourceLanguage: "en",
		TargetLanguage: "de",
		Content: []requests.Content{
			{
				TimeRange: "00:00:12.00 - 00:01:20.00",
				Content:   "I am Arwen - I've come to help you.",
			},
			{
				TimeRange: "00:03:55.00 - 00:04:20.00",
				Content:   "Come back to the light.",
			},
		},
	}

	b, err := json.Marshal(req)
	require.NoError(t, err)
	require.NotNil(t, b)

	token := createToken(t)

	res := requestHelper(t, http.MethodPost, "/sts", token, b)
	require.Equal(t, http.StatusCreated, res.Code)

	var result responses.Subtitle
	err = json.Unmarshal(res.Body.Bytes(), &result)
	require.NoError(t, err)
	require.NotEmpty(t, result)
	assert.NotEmpty(t, result.ID)
	assert.Equal(t, req.Name, result.Name)
	assert.Equal(t, req.Filename, result.Filename)

	assert.Len(t, result.Content, 2)
	contents := result.Content

	assert.NotEmpty(t, contents[0].ID)
	assert.Equal(t, req.Content[0].TimeRange, contents[0].TimeRange)
	assert.Equal(t, "Ich bin Arwen - Ich bin gekommen, um dir zu helfen.", contents[0].Content)

	assert.NotEmpty(t, contents[1].ID)
	assert.Equal(t, req.Content[1].TimeRange, contents[1].TimeRange)
	assert.Equal(t, "Komm zurück zum Licht.", contents[1].Content)
}

func TestUploadSubtitleFiles(t *testing.T) {
	token := createToken(t)
	fileNames := []string{
		"./testdata/subtitles.txt",
		"./testdata/subtitles_de.txt",
	}

	ctx := context.Background()
	// upload some translations
	tr := []*tms.Translation{
		{
			SourceLanguage: "en",
			TargetLanguage: "de",
			Source:         "I am Arwen - I've come to help you.",
			Target:         "Ich bin Arwen - Ich bin gekommen, um dir zu helfen.",
		},
		{
			SourceLanguage: "en",
			TargetLanguage: "de",
			Source:         "Come back to the light.",
			Target:         "Komm zurück zum Licht.",
		},
		{
			TargetLanguage: "en",
			SourceLanguage: "de",
			Target:         "I am Arwen - I've come to help you.",
			Source:         "Ich bin Arwen - Ich bin gekommen, um dir zu helfen.",
		},
		{
			TargetLanguage: "en",
			SourceLanguage: "de",
			Target:         "Come back to the light.",
			Source:         "Komm zurück zum Licht.",
		},
	}
	for _, v := range tr {
		require.NoError(t, server.translateService.Create(ctx, v))
	}

	var b bytes.Buffer
	w := httptest.NewRecorder()
	writer := multipart.NewWriter(&b)

	for _, v := range fileNames {
		form, err := writer.CreateFormFile("files", v)
		require.NoError(t, err)

		file, err := os.Open(v)
		require.NoError(t, err)
		_, err = io.Copy(form, file)
		require.NoError(t, err)
		file.Close()
	}
	name, err := writer.CreateFormField("name")
	require.NoError(t, err)
	_, err = name.Write([]byte("lord of the rings"))
	require.NoError(t, err)

	sourceLang, err := writer.CreateFormField("source_language")
	require.NoError(t, err)
	_, err = sourceLang.Write([]byte("en"))
	require.NoError(t, err)

	targetLang, err := writer.CreateFormField("target_language")
	require.NoError(t, err)
	_, err = targetLang.Write([]byte("de"))
	require.NoError(t, err)
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, "/sts/upload", &b)
	require.NoError(t, err)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)

	server.Router.ServeHTTP(w, req)

	println(w.Body.String())

	// require.Equal(t, http.StatusOK, w.Code)
	// exp := `"translations uploaded successfully"`
	// assert.Equal(t, exp, w.Body.String())

}
