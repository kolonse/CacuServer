package function

type Function interface {
	Parse(string)
	SetCallArgs(args ...interface{})
	Call() (interface{}, error)
}
