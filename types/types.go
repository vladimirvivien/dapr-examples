package types

type Order struct {
	ID        string
	Items     []string
	Received  bool
	Completed bool
}