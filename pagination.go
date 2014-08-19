package govenmo

// Pagination stores the 'next' link that indicates a continuation of the Venmo response.
// This may be incorrect when retrieved from the sandbox environments.
// This is used internally.
type Pagination struct {
	Next string
}
