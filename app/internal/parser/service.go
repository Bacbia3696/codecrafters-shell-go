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

// ParseLineWithMode parses a command line into arguments and redirection targets with append mode information
func (s *Service) ParseLineWithMode(line string) (args []string, outputFile string, errorFile string, outputAppend bool, errorAppend bool, err error) {
	return ParseLineWithMode(line)
}
