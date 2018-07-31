package storage

import (
	"github.com/groundstorm/cobalt/src/models"
)

// MemoryStorageUser overrides User so we can more easily implement the
// memory storage backend
type memoryStorageUser struct {
	models.User
	password string
}
