// Licensed under the Apache License, Version 2.0
// Details: https://raw.githubusercontent.com/maniksurtani/quotaservice/master/LICENSE

package config

import (
	"io"
)

// ConfigPersister is an interface that persists configs and notifies a channel of changes.
type ConfigPersister interface {
	// PersistAndNotify persists a marshalled configuration passed in.
	PersistAndNotify(marshalledConfig io.Reader) error
	// ConfigChangedWatcher returns a channel that is notified whenever configuration changes are
	// detected. Changes are coalesced so that a single notification may be emitted for multiple
	// changes.
	ConfigChangedWatcher() chan struct{}
	// ReadPersistedConfig provides a reader to a marshalled config previously persisted.
	ReadPersistedConfig() (marshalledConfig io.Reader, err error)
}
