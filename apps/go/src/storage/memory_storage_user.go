package storage

import (
	"github.com/groundstorm/cobalt/apps/go/src/models"
)

// MemoryStorageUser overrides User so we can more easily implement the
// memory storage backend
type memoryStorageUser struct {
	models.User
	password string
}
