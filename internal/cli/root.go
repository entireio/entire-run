package cli

import (
	"github.com/spf13/cobra"
)

type Options struct {
	Version string
	Env     EntireEnv
}

// Execute runs the plugin root command with the real process environment.
func Execute(version string) error {
	return NewRootCommand(Options{
		Version: version,
		Env:     EnvFromOS(),
	}).Execute()
}

func NewRootCommand(opts Options) *cobra.Command {
	if opts.Version == "" {
		opts.Version = "dev"
	}

	var yolo bool
	cmd := &cobra.Command{
		Use:           "entire-run [agent] [args...]",
		Short:         "Launch an Entire-enabled agent in the current directory",
		Args:          cobra.ArbitraryArgs,
		SilenceUsage:  true,
		SilenceErrors: true,
		Long: `entire-run lists the agents enabled for Entire in the current
repository, lets you pick one, and launches that agent in the current
directory.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLauncher(cmd.Context(), cmd.OutOrStdout(), args, launcherOptions{Yolo: yolo})
		},
	}

	cmd.Flags().BoolVar(&yolo, "yolo", false, "Launch the selected agent in its no-approval mode when supported")
	cmd.Flags().SetInterspersed(false)
	cmd.AddCommand(newDoctorCommand(opts.Env))
	cmd.AddCommand(newVersionCommand(opts.Version))
	return cmd
}
