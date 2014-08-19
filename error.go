package govenmo

// Error stores error information from the Venmo API.
// This is used internally.
type Error struct {
	Message string
	Code    int
}
