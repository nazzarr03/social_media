package utils

import (
	"math/rand"
	"net/url"
	"time"

	"github.com/nazzarr03/social-media/models"
)

func CreateShortURL(longURL string) models.ShortURL {
	if longURL == "" {
		return models.ShortURL{}
	}

	_, err := url.ParseRequestURI(longURL)
	if err != nil {
		return models.ShortURL{}
	}

	shortURL := models.ShortURL{
		LongURL:  longURL,
		ShortKey: GenerateShortKey(),
	}

	return shortURL
}

func GenerateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const keyLength = 6
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, keyLength)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}
	return string(shortKey)
}
