package unknown

type unknownCommand struct{}

func New() *unknownCommand {
	return &unknownCommand{}
}

func (s *unknownCommand) Process(msg string) (string, error) {
	return "не знаю такую команду", nil
}
