package javascript

import (
	"time"

	"github.com/grafana/sobek"
)

type JavaScript struct {
	options Options
	vm      *sobek.Runtime
}

type Options struct {
	Script  string
	Timeout time.Duration
}

func New(options Options) *JavaScript {
	vm := sobek.New()

	return &JavaScript{
		vm:      vm,
		options: options,
	}
}

func (js *JavaScript) Execute() error {
	// init functions
	// call user function
	return nil
}
