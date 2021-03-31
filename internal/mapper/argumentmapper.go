package mapper

// Argument represents a function to map arguments to resource UUIDs
type Argument func(arg string) (uuid string, err error)
