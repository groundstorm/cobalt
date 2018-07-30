package storage

import (
	"github.com/groundstorm/cobalt/src/users"
)

// MemoryStorageUser overrides users.User so we can more easily implement the
// memory storage backend
type memoryStorageUser struct {
	users.User
	password string
}
