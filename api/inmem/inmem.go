package inmem

import (
	"sync"

	url "xsurl/api/shortenurl"
)

type urlRepository struct {
	mtx  sync.RWMutex
	urls map[url.ID]*url.URL
}

func (r *urlRepository) Store(u *url.URL) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	r.urls[u.ID] = u
	return nil
}

func (r *urlRepository) FindByID(id url.ID) (*url.URL, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	if val, ok := r.urls[id]; ok {
		return val, nil
	}
	return nil, url.ErrUnknownURL
}

func (r *urlRepository) FindAll() []*url.URL {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	u := make([]*url.URL, 0, len(r.urls))
	for _, val := range r.urls {
		u = append(u, val)
	}
	return u
}

func NewURLRepository() url.URLRepository {
	return &urlRepository{
		urls: make(map[url.ID]*url.URL),
	}
}
