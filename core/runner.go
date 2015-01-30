package core

// Runner is the interface which describes objects
// that can execute tasks on a system.
type Runner interface {
	Run(Task) error
	Close()
}
