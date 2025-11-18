package store

import "sync"

func NewRepository() *Repository {
	return &Repository{
		catalogs: make(map[int][]string),
	}
}

type Repository struct {
	catalogs map[int][]string
	mu       sync.Mutex
}

func (r *Repository) Get(key int) []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.catalogs[key]
}

func (r *Repository) Set(key int, values []string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.catalogs[key] = values
}

func (r *Repository) Append(key int, value string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.catalogs[key] = append(r.catalogs[key], value)
}

func (r *Repository) Delete(key int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.catalogs, key)
}
