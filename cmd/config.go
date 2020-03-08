package cmd

import (
	"demo_project/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"os"
	"text/template"
)

// when updating this template, don't forget to update config.md!
const configTemplate = `[general]
# Log level
#
# debug=5, info=4, warning=3, error=2, fatal=1, panic=0
log_level={{ .General.LogLevel }}
`

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print the config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		t := template.Must(template.New("config").Parse(configTemplate))
		err := t.Execute(os.Stdout, config.Conf)
		if err != nil {
			return errors.Wrap(err, "execute config template error")
		}
		return nil
	},
}