package main

// Loader is an interface for abstracting the transformations necessary to load
// samples onto cards for various devices
type Loader interface {
	// Adds the file at the given path to the card
	AddSample(path string) (err error, full bool)
}
