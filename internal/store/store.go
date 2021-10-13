package store

// Factory defines the storage interface.
type Factory interface {
	Talks() TalkStore
	Messages() MessageStore
	Close() error
}
