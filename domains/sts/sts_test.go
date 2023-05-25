package sts

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/brianvoe/gofakeit"
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
		cleanup()
		os.Exit(code)
	}()
	db, initError = setupTestDB()
	if initError != nil {
		log.Fatal(initError)
	}
	code = m.Run()
}

func TestCreateSubtitle(t *testing.T) {
	s := &Subtitle{
		Name:           gofakeit.CarMaker(),
		Filename:       strings.ToLower(gofakeit.BuzzWord()),
		SourceLanguage: "en",
	}

	ctx := context.Background()

	repo, err := New(db)
	require.NoError(t, err)
	require.NotNil(t, repo)

	err = repo.Create(ctx, s)
	require.NoError(t, err)
	assert.NotEmpty(t, s.ID)

	// t.Cleanup(func() {
	// 	if err := db.Exec("DELETE FROM subtitles WHERE id = ?", s.ID).Error; err != nil {
	// 		log.Fatal(err)
	// 	}
	// })

	res, err := repo.FindSubtitle(ctx, s.Name)
	require.NoError(t, err)
	require.NotEmpty(t, res)

	assert.Equal(t, s.ID, res.ID)
	assert.Equal(t, s.CreatedAt, res.CreatedAt)
}

// func TestCreateContent(t *testing.T) {
// 	c := &Content{
// 		SubtitleID: uuid.NewString(),
// 		TimeStart:  "00:01:20.00",
// 		TimeStop:   "00:02:00.00",
// 		Content:    "Hello World",
// 	}

// 	ctx := context.Background()

// 	repo, err := New(db)
// 	require.NoError(t, err)
// 	require.NotNil(t, repo)

// 	err = repo.CreateContent(ctx, c)
// 	require.NoError(t, err)
// 	assert.NotEmpty(t, c.ID)

// 	if err := db.Exec("DELETE FROM contents WHERE id = ?", c.ID).Error; err != nil {
// 		log.Fatal(err)
// 	}
// }

// func TestFindContentByTimestamp(t *testing.T) {
// 	subtitleID := uuid.NewString()
// 	ctx := context.Background()
// 	ids := []string{}

// 	repo, err := New(db)
// 	require.NoError(t, err)
// 	require.NotNil(t, repo)

// 	c := []*Content{
// 		{
// 			SubtitleID: subtitleID,
// 			TimeStart:  "00:00:12.00",
// 			TimeStop:   "00:01:20.00",
// 			TimeRange:  "00:00:12.00 - 00:01:20.00",
// 			Content:    "Ich bin Arwen - Ich bin gekommen, um dir zu helfen.",
// 		},
// 		{
// 			SubtitleID: subtitleID,
// 			TimeStart:  "00:03:55.00",
// 			TimeStop:   "00:04:20.00",
// 			TimeRange:  "00:03:55.00 - 00:04:20.00",
// 			Content:    "Komm zurück zum Licht.",
// 		}, {
// 			SubtitleID: subtitleID,
// 			TimeStart:  "00:04:59.00",
// 			TimeStop:   "00:05:30.00",
// 			TimeRange:  "00:04:59.00 - 00:05:30.00",
// 			Content:    " Nein, my Schatz!!.",
// 		},
// 	}

// 	t.Cleanup(func() {
// 		if err := db.Exec("DELETE FROM contents WHERE id IN ?", ids).Error; err != nil {
// 			log.Fatal(err)
// 		}
// 	})

// 	for _, v := range c {
// 		require.NoError(t, repo.CreateContent(ctx, v))
// 		ids = append(ids, v.ID)
// 	}

// 	res, err := repo.FindContentByTimeRange(ctx, subtitleID, "00:03:55.00 - 00:04:20.00", "00:04:59.00 - 00:05:30.00")
// 	require.NoError(t, err)
// 	require.NotNil(t, res)

// 	assert.Len(t, res, 2)
// }

func setupTestDB() (*gorm.DB, error) {
	env := os.Getenv("ENVIRONMENT")
	dsn := "root:@tcp(127.0.0.1:3306)/translations?charset=utf8mb4&parseTime=True&loc=Local"
	if env == "cicd" {
		dsn = "test_user:password@tcp(127.0.0.1:33306)/translations?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func cleanup() {
	if err := db.Exec("DELETE FROM contents").Error; err != nil {
		log.Fatal(err)
	}

	if err := db.Exec("DELETE FROM subtitles").Error; err != nil {
		log.Fatal(err)
	}
}
