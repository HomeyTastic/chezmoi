package cmd

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/twpayne/chezmoi/next/internal/chezmoi"
)

type diffCmdConfig struct {
	include   *chezmoi.IncludeSet
	recursive bool
	NoPager   bool
	Pager     string
}

func (c *Config) newDiffCmd() *cobra.Command {
	diffCmd := &cobra.Command{
		Use:     "diff [target]...",
		Short:   "Print the diff between the target state and the destination state",
		Long:    mustLongHelp("diff"),
		Example: example("diff"),
		RunE:    c.runDiffCmd,
		Annotations: map[string]string{
			persistentStateMode: persistentStateModeReadOnly,
		},
	}

	flags := diffCmd.Flags()
	flags.VarP(c.Diff.include, "include", "i", "include entry types")
	flags.BoolVar(&c.Diff.NoPager, "no-pager", c.Diff.NoPager, "disable pager")
	flags.BoolVarP(&c.Diff.recursive, "recursive", "r", c.Diff.recursive, "recursive")

	return diffCmd
}

func (c *Config) runDiffCmd(cmd *cobra.Command, args []string) error {
	sb := strings.Builder{}
	dryRunSystem := chezmoi.NewDryRunSystem(c.destSystem)
	gitDiffSystem := chezmoi.NewGitDiffSystem(dryRunSystem, &sb, c.normalizedDestDir, c.color)
	if err := c.applyArgs(gitDiffSystem, c.normalizedDestDir, args, c.Diff.include, c.Diff.recursive, c.Umask.FileMode()); err != nil {
		return err
	}
	return c.writeOutputString(sb.String())
}
