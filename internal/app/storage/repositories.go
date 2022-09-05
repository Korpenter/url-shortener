package storage

type Repository interface {
	Get(id string) (string, error)
	Add(id string) string
	NewID() int
}
