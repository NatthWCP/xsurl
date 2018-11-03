package shortenurl

import (
	"errors"
	"strings"

	"github.com/pborman/uuid"
)

// ID is an indentifier of each URL comming from user
type ID string

// URL is an central object use on shortening service
type URL struct {
	ID         ID
	OriginURL  string
	ShortenURL string
}

// NewURL create URL object which contains ID and URL from user
func NewURL(id ID, ourl string, surl string) *URL {
	return &URL{
		ID:         id,
		OriginURL:  ourl,
		ShortenURL: surl,
	}
}

// URLRepository is used for URL Datastore accessing
type URLRepository interface {
	Store(URL *URL) error
	FindByID(id ID) (*URL, error)
	// FindByOriginURL(ourl string) (*URL, error)
}

// NextID generates a new ID by getting the first section of UUID
func NextID() ID {
	return ID(strings.Split(strings.ToUpper(uuid.New()), "-")[0])
}

// ErrUnknownURL occurs when the target URL couldn't be found
var ErrUnknownURL = errors.New("unknown URL")
