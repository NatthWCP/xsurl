package shortening

import (
	"time"

	kitlog "github.com/go-kit/kit/log"
)

type loggingService struct {
	logger kitlog.Logger
	next   Service
}

func NewLoggingService(logger kitlog.Logger, s Service) Service {
	return &loggingService{logger, s}
}

func (s *loggingService) ShortenURL(originURL string) (shortenURL string, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "ShortenURL",
			"originURL", originURL,
			"shortenURL", shortenURL,
			"err", err,
		)
	}(time.Now())

	return s.next.ShortenURL(originURL)
}
