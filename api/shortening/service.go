package shortening

import (
	"crypto/md5"
	"encoding/base64"
	"errors"

	url "xsurl/api/shortenurl"
)

var ErrInvalidArgument = errors.New("invalid argument")

type Service interface {
	ShortenURL(originURL string) (string, error)
}

type service struct {
	urls url.URLRepository
}

func shortenURL(originURL string) string {
	// business rule for shortening the URL will be implement here
	hashedbytes := md5.Sum([]byte(originURL))
	return base64.StdEncoding.EncodeToString(hashedbytes[:])
}

func (s *service) ShortenURL(originURL string) (string, error) {
	if originURL == "" {
		return "", ErrInvalidArgument
	}

	id := url.NextID()
	shortenURL := shortenURL(originURL)

	u := url.NewURL(id, originURL, shortenURL)

	if err := s.urls.Store(u); err != nil {
		return "", err
	}

	return u.ShortenURL, nil
}

func NewService(urls url.URLRepository) Service {
	return &service{
		urls: urls,
	}
}
