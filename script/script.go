package script

type Script interface {
	SetCallArgs(...interface{})
	Parse(str string)
	Call() (interface{}, error)
}
