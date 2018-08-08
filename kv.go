package kvapi

// KV stores and retrieves keys and values, scoped by UID.
type KV interface {
	// Get retrieves a value by key for the user.
	Get(UID, key string) string

	// Set stores a value by key for the user.
	Set(UID, key, value string)
}
