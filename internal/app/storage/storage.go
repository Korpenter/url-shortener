package storage

// Repository interface for storage instances
type Repository interface {
	Get(id string) (string, error)
	Add(id string) string
	NewID() int
}
