package chezmoi

// An EmptyPersistentState is a PersistentState that returns nil for all reads
// and panics on any writes.
type EmptyPersistentState struct{}

// Close does nothing.
func (EmptyPersistentState) Close() error { return nil }

// CopyTo does nothing.
func (EmptyPersistentState) CopyTo(PersistentState) error { return nil }

// Delete panics.
func (EmptyPersistentState) Delete([]byte, []byte) error { panic(nil) }

// ForEach does nothing.
func (EmptyPersistentState) ForEach([]byte, func([]byte, []byte) error) error { return nil }

// Get returns nil.
func (EmptyPersistentState) Get([]byte, []byte) ([]byte, error) { return nil, nil }

// Set panics.
func (EmptyPersistentState) Set([]byte, []byte, []byte) error { panic(nil) }
