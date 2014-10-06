package schema

type Schema interface {
	KeySize() int
	MACSize() int
}
