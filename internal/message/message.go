package message

type ResultMessage interface {
	String() string
	CSV() string
}
