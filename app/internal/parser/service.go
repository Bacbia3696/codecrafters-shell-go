package parser

// Service provides command parsing functionality
type Service struct{}

// NewService creates a new parser service
func NewService() *Service {
	return &Service{}
}

// ParseLine parses a command line into arguments and redirection targets
func (s *Service) ParseLine(line string) (args []string, outputFile string, errorFile string, err error) {
	return ParseLine(line)
}
