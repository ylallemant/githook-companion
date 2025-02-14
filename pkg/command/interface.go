package command

type Command interface {
	Execute() (string, error)
	AddArg(argument string)
}
