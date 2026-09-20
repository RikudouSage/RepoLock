package errors

type SerializableError interface {
	SerializeToStruct() any
}

type ErrorWithStatusCode interface {
	StatusCode() int
}

type ErrorWithCode interface {
	GetCode() string
}
