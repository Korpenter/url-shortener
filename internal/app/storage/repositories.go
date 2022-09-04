package storage

type Repositories interface {
	Get(id string) (string, error)
	Add(id string) string
	NewID() int
}
