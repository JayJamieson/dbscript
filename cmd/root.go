package cmd

import (
	"fmt"
	"os"

	"github.com/grafana/sobek"
	"github.com/spf13/cobra"
)

var script = `
function handle() {
	dbscript.ok({
		"foo": "bar",
		"bar": 123
	});
}
`

var rootCmd = &cobra.Command{
	Use:   "dbscript",
	Short: "Easily script your CDC events",
	Long: `Script your CDC events with Javascript. Build complex or simple pre-processing
pipelines in Javascript.

Javscript files can be hot reloading to avoid downtime.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		vm := sobek.New()

		vm.SetFieldNameMapper(sobek.TagFieldNameMapper("json", true))

		vm.GlobalObject().Set("dbscript", struct {
			Ok func(event any) `json:"ok"`
		}{
			Ok: func(event any) {
				fmt.Println(event)
			},
		})

		program, err := sobek.Compile("", script, false)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		_, err = vm.RunProgram(program)

		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(1)
		}

		var handle func() (interface{}, error)
		handleFunction := vm.Get("handle")

		if handleFunction == nil {
			fmt.Println("handle is not defined")
			os.Exit(1)
		}

		err = vm.ExportTo(handleFunction, &handle)
		if err != nil {
			return
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
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dbscript.yaml)")

	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
