package store

// Factory defines the storage interface.
type Factory interface {
	Talks() TalkStore
	Close() error
}
