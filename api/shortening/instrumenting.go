package shortening

import (
	"time"

	"github.com/go-kit/kit/metrics"
)

type instrumentingService struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func NewInstrumentingService(counter metrics.Counter, latency metrics.Histogram, s Service) Service {
	return &instrumentingService{
		requestCount:   counter,
		requestLatency: latency,
		next:           s,
	}
}

func (s *instrumentingService) ShortenURL(originURL string) (string, error) {
	defer func(begin time.Time) {
		s.requestCount.With("method", "ShortenURL").Add(1)
		s.requestLatency.With("method", "ShortenURL").Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.next.ShortenURL(originURL)
}
