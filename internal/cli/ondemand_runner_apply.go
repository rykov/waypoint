package cli

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	hcljson "github.com/hashicorp/hcl/v2/json"
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	"github.com/posener/complete"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OndemandRunnerApplyCommand struct {
	*baseCommand

	flagId           string
	flagOCIUrl       string
	flagEnvVars      []string
	flagPluginType   string
	flagPluginConfig string
	flagDefault      bool
}

func (c *OndemandRunnerApplyCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	flagSet := c.Flags()
	if err := c.Init(
		WithArgs(args),
		WithFlags(flagSet),
		WithNoConfig(),
	); err != nil {
		return 1
	}
	args = flagSet.Args()
	ctx := c.Ctx

	sg := c.ui.StepGroup()
	defer sg.Wait()

	var s terminal.Step

	defer func() {
		if s != nil {
			s.Abort()
		}
	}()

	var (
		od      *pb.OndemandRunner
		updated bool
	)

	if c.flagId != "" {
		s = sg.Add("Checking for an existing ondemand runner: %s", c.flagId)
		// Check for an existing project of the same name.
		resp, err := c.project.Client().GetOndemandRunner(ctx, &pb.GetOndemandRunnerRequest{
			OndemandRunner: &pb.Ref_OndemandRunner{
				Id: c.flagId,
			},
		})
		if status.Code(err) == codes.NotFound {
			// If the error is a not found error, act as though there is no error
			// and the project is nil so that we can handle that later.
			resp = nil
			err = nil
		}
		if err != nil {
			c.ui.Output(
				"Error checking for project: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}

		od = resp.OndemandRunner
		s.Update("Updating ondemand runner %q...", od.Id)
		updated = true
	} else {
		s = sg.Add("Creating new ondemand runner")
		od = &pb.OndemandRunner{}
	}

	// If we were specified a file then we're going to load that up.
	if c.flagPluginConfig != "" {
		path, err := filepath.Abs(c.flagPluginConfig)
		if err != nil {
			c.ui.Output(
				"Error loading HCL file specified with the -plugin-config flag: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			c.ui.Output(
				"Error reading HCL plugin config: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)

			return 1
		}

		od.PluginConfig = data

		switch filepath.Ext(path) {
		case ".hcl":
			od.ConfigFormat = pb.Project_HCL
			_, diag := hclsyntax.ParseConfig(data, "<waypoint-hcl>", hcl.Pos{})
			if diag.HasErrors() {
				c.ui.Output(
					"Syntax errors in file specified with -plugin-config: %s",
					clierrors.Humanize(diag),
					terminal.WithErrorStyle(),
				)

				return 1
			}

		case ".json":
			od.ConfigFormat = pb.Project_JSON
			_, diag := hcljson.Parse(data, "<waypoint-hcl>")
			if diag.HasErrors() {
				c.ui.Output(
					"Syntax errors in file specified with -plugin-config: %s",
					clierrors.Humanize(diag),
					terminal.WithErrorStyle(),
				)

				return 1
			}

		default:
			c.ui.Output(
				"File specified via -plugin-config must end in '.hcl' or '.json'",
				terminal.WithErrorStyle(),
			)

			return 1
		}
	}

	od.OciUrl = c.flagOCIUrl
	od.EnvironmentVariables = map[string]string{}
	od.Default = c.flagDefault

	for _, kv := range c.flagEnvVars {
		idx := strings.IndexByte(kv, '=')
		if idx != -1 {
			od.EnvironmentVariables[kv[:idx]] = kv[idx+1:]
		}
	}

	od.PluginType = c.flagPluginType

	// Upsert
	_, err := c.project.Client().UpsertOndemandRunner(ctx, &pb.UpsertOndemandRunnerRequest{
		OndemandRunner: od,
	})
	if err != nil {
		c.ui.Output(
			"Error upserting ondemand runner: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	if updated {
		s.Update("Ondemand Runner updated")
	} else {
		s.Update("Ondemand Runner created")
	}
	s.Done()

	return 0
}

func (c *OndemandRunnerApplyCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(sets *flag.Sets) {
		f := sets.NewSet("Command Options")

		f.StringVar(&flag.StringVar{
			Name:    "id",
			Target:  &c.flagId,
			Default: "",
			Usage:   "The id an existing ondemand runner to update.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "oci-url",
			Target:  &c.flagOCIUrl,
			Default: "hashicorp/waypoint:stable",
			Usage:   "The url for the OCI image to launch for the ondemand runner.",
		})

		f.StringSliceVar(&flag.StringSliceVar{
			Name:   "env-vars",
			Target: &c.flagEnvVars,
			Usage: "Environment variable to expose to the ondemand runner. Typically used to " +
				"introduce configuration for the plugins that the runner will execute.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "plugin-type",
			Target:  &c.flagPluginType,
			Default: "",
			Usage:   "The type of the plugin to launch for the ondemand runner, such as aws-ecs, kubernetes, etc.",
		})

		f.StringVar(&flag.StringVar{
			Name:    "project-config",
			Target:  &c.flagPluginConfig,
			Default: "",
			Usage: "Path to a an hcl file that contains the configuration for the plugin." +
				"This is only necessary when the plugin's defaults need to be adjusted for " +
				"the environment the plugin will launch the ondemand runner.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "default",
			Target:  &c.flagDefault,
			Default: false,
			Usage: "Indicates that this ondemand runner should be used by any project that doesn't " +
				"otherwise specify it's own ondemand runner.",
		})
	})
}

func (c *OndemandRunnerApplyCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *OndemandRunnerApplyCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *OndemandRunnerApplyCommand) Synopsis() string {
	return "Create or update an ondemand runner."
}

func (c *OndemandRunnerApplyCommand) Help() string {
	return formatHelp(`
Usage: waypoint ondemand-runner apply [OPTIONS]

  Create or update an ondemand runner.

  This will create a new ondemand runner with the given options. If
	a ondemand runner with the same id already exists, this will update the
	existing runner using the fields that are set.

` + c.Flags().Help())
}
