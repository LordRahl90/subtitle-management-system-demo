package servers

import (
	"encoding/json"
	"net/http"
	"testing"
	"translations/requests"

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

	res := requestHelper(t, http.MethodPost, "/tms", token, b)
	require.Equal(t, http.StatusCreated, res.Code)
}
