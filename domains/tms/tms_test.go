package tms

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db        *gorm.DB
	initError error
)

func TestMain(m *testing.M) {
	code := 1
	defer func() {
		os.Exit(code)
	}()
	db, initError = setupTestDB()
	if initError != nil {
		log.Fatal(initError)
	}
	code = m.Run()
}

func TestCreate(t *testing.T) {
	e := &Translation{
		SourceLanguage: "en",
		TargetLanguage: "de",
		Source:         "Hello World",
		Target:         "Hallo Welt",
	}

	repo, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, repo)

	require.NoError(t, repo.Create(context.TODO(), e))
	assert.NotEmpty(t, e.ID)

	t.Cleanup(func() {
		if err := db.Exec("DELETE FROM translations WHERE id = ?", e.ID).Error; err != nil {
			log.Fatal(err)
		}
	})
}

func TestFindTranslation(t *testing.T) {
	ctx := context.Background()
	ids := []string{}

	repo, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, repo)

	e := []*Translation{
		{
			Source:         "Hello World",
			Target:         "Hallo Welt",
			SourceLanguage: "en",
			TargetLanguage: "de",
		},
		{
			Source:         "Hello guys",
			Target:         "Hallo Leute",
			SourceLanguage: "en",
			TargetLanguage: "de",
		},
		{
			Source:         "I walk to the supermarket",
			Target:         "Ich gehe zum Supermarkt.",
			SourceLanguage: "en",
			TargetLanguage: "de",
		},
	}
	for _, v := range e {
		require.NoError(t, repo.Create(ctx, v))
		ids = append(ids, v.ID)
	}
	t.Cleanup(func() {
		if err := db.Exec("DELETE FROM translations WHERE id IN ?", ids).Error; err != nil {
			log.Fatal(err)
		}
	})

	res, err := repo.Find(ctx, "en", "de", "I walk to the supermarket")
	require.NoError(t, err)
	require.NotEmpty(t, res)
	assert.NotEmpty(t, res.Target)

	res, err = repo.FindByID(ctx, ids[2])
	require.NoError(t, err)
	require.NotEmpty(t, res)

	assert.Equal(t, res.Target, e[2].Target)
}

func TestFind_NonExistingTranslation(t *testing.T) {
	ctx := context.Background()
	repo, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, repo)

	res, err := repo.Find(ctx, "en", "da", "I walk to the supermarket")
	require.EqualError(t, err, gorm.ErrRecordNotFound.Error())
	require.Empty(t, res)
}

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/translations?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/translations?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
