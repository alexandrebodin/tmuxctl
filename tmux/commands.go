package tmux

var (
	// DefaultRunner is the default tmux command runner
	DefaultRunner = &Runner{}
)

// Exec runs a tmux command
func Exec(args ...string) (Result, error) {
	return DefaultRunner.Exec(args...)
}

// SendKeys sends keys to tmux (e.g to run a command)
func SendKeys(target, keys string) error {
	return DefaultRunner.SendKeys(target, keys)
}

// SendRawKeys sends keys to tmux (e.g to run a command)
func SendRawKeys(target, keys string) error {
	return DefaultRunner.SendRawKeys(target, keys)
}

// ListSessions returns the list of sessions currently running
func ListSessions() (map[string]SessionInfo, error) {
	return DefaultRunner.ListSessions()
}

// GetOptions get tmux options
func GetOptions() (*Options, error) {
	return DefaultRunner.GetOptions()
}
