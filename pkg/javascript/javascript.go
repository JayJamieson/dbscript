package javascript

import (
	"fmt"
	"os"
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

	js.vm.SetFieldNameMapper(sobek.TagFieldNameMapper("json", true))

	js.vm.GlobalObject().Set("dbscript", &Runtime{
		Context: runtimeCtx{},
	})

	program, err := sobek.Compile("", js.options.Script, false)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	_, err = js.vm.RunProgram(program)

	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}

	var handle func() (interface{}, error)
	handleFunction := js.vm.Get("handle")

	if handleFunction == nil {
		fmt.Println("handle is not defined")
		os.Exit(1)
	}

	err = js.vm.ExportTo(handleFunction, &handle)
	if err != nil {
		return err
	}

	output, err := handle()
	if err != nil {
		switch e := err.(type) {
		case *sobek.InterruptedError:
			fmt.Printf("%v\n", e)
			os.Exit(1)

		case *sobek.Exception:
			fmt.Printf("%v\n", e)
			os.Exit(1)
		}
	}
	fmt.Printf("output: %v\n", output)
	return nil
}
