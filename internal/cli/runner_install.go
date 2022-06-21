package cli

import (
	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/clierrors"
	"github.com/hashicorp/waypoint/internal/installutil"
	"github.com/hashicorp/waypoint/internal/pkg/flag"
	"github.com/hashicorp/waypoint/internal/runnerinstall"
	"github.com/hashicorp/waypoint/pkg/server"
	pb "github.com/hashicorp/waypoint/pkg/server/gen"
	"github.com/hashicorp/waypoint/pkg/serverconfig"
	"github.com/posener/complete"
	empty "google.golang.org/protobuf/types/known/emptypb"
	"sort"
	"strings"
)

type RunnerInstallCommand struct {
	*baseCommand

	platform              []string `hcl:"platform,optional"`
	adopt                 bool     `hcl:"adopt,optional"`
	serverUrl             string   `hcl:"server_url,required"`
	id                    string   `hcl:"id,optional"`
	runnerProfileOdrImage string   `hcl:"odr_image,optional"`
	serverTls             bool     `hcl:"server_tls,optional"`
	serverTlsSkipVerify   bool     `hcl:"server_tls_skip_verify,optional"`
	serverRequireAuth     bool     `hcl:"server_require_auth,optional"`
}

func (c *RunnerInstallCommand) AutocompleteArgs() complete.Predictor {
	return complete.PredictNothing
}

func (c *RunnerInstallCommand) AutocompleteFlags() complete.Flags {
	return c.Flags().Completions()
}

func (c *RunnerInstallCommand) Flags() *flag.Sets {
	return c.flagSet(0, func(set *flag.Sets) {
		f := set.NewSet("Command Options")

		var runnerPlatforms []string
		for platformName := range runnerinstall.Platforms {
			runnerPlatforms = append(runnerPlatforms, platformName)
		}

		f.EnumVar(&flag.EnumVar{
			Name:   "platform",
			Usage:  "Platform to install the Waypoint runner into. If unset, uses the platform of the local context.",
			Values: runnerPlatforms,
			Target: &c.platform,
		})

		f.StringVar(&flag.StringVar{
			Name:   "server-addr",
			Usage:  "Address of the Waypoint server.",
			EnvVar: "WAYPOINT_ADDR",
			Target: &c.serverUrl,
		})

		f.StringVar(&flag.StringVar{
			Name:    "odr-image",
			Usage:   "Docker image for the on-demand runners. If not specified, it defaults to 'hashicorp/waypoint-odr:latest'.",
			Default: "hashicorp/waypoint-odr:latest",
			Target:  &c.runnerProfileOdrImage,
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "server-tls",
			Target:  &c.serverTls,
			Usage:   "If true, the Waypoint runner will connect to the server over TLS.",
			Default: true,
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "server-tls-skip-verify",
			Target: &c.serverTlsSkipVerify,
			Usage:  "If true, will not validate TLS cert presented by the server.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:   "server-require-auth",
			Target: &c.serverRequireAuth,
			Usage:  "If true, will send authentication details.",
		})

		f.BoolVar(&flag.BoolVar{
			Name:    "adopt",
			Usage:   "Adopt the runner after it is installed.",
			Default: true,
			Target:  &c.adopt,
		})

		f.StringVar(&flag.StringVar{
			Name:   "id",
			Usage:  "If this is set, the runner will use the specified id.",
			Target: &c.id,
		})

		sort.Strings(runnerPlatforms)

		for _, name := range runnerPlatforms {
			platform := runnerinstall.Platforms[name]
			platformSet := set.NewSet(name + " Options")
			platform.InstallFlags(platformSet)
		}
	})
}

func (c *RunnerInstallCommand) Help() string {
	return formatHelp(`
Usage: waypoint runner install [options]

  Install a Waypoint runner to an existing platform. The platform should be
  specified as kubernetes, nomad, ecs, or docker.

  With the '-adopt' flag set to true (the default), this command will adopt the
  runner after it is installed. The install will attempt to install a runner for
  the server configured in the current Waypoint context.

` + c.Flags().Help())
}

func (c *RunnerInstallCommand) Synopsis() string {
	return "Install a Waypoint runner to Kubernetes, Nomad, ECS, or Docker"
}

func (c *RunnerInstallCommand) Run(args []string) int {
	// Initialize. If we fail, we just exit since Init handles the UI.
	if err := c.Init(
		WithArgs(args),
		WithFlags(c.Flags()),
		WithNoLocalServer(), // no auth in local mode
		WithNoConfig(),
	); err != nil {
		return 1
	}

	ctx := c.Ctx
	log := c.Log.Named("install")
	defer c.Close()

	serverConfig, err := c.project.Client().GetServerConfig(ctx, &empty.Empty{})
	if err != nil {
		c.ui.Output(
			"Error getting server config.",
			clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}

	// If the user doesn't set a platform, set platform to the platform of the user's context
	platform := c.platform
	if len(platform) == 0 {
		platform = append(platform, serverConfig.Config.Platform)
	}

	p, ok := runnerinstall.Platforms[strings.ToLower(platform[0])]
	if !ok {
		c.ui.Output(
			"Error installing server into %q: unsupported platform",
			platform[0],
			terminal.WithErrorStyle(),
		)
		return 1
	}

	sg := c.ui.StepGroup()
	defer sg.Wait()

	s := sg.Add("Connecting to: %s", c.serverUrl)
	defer func() { s.Abort() }()

	var cookie string
	if c.adopt {
		cookie = serverConfig.Config.Cookie
	}

	if c.serverUrl == "" {
		c.ui.Output(
			"-server-addr must be supplied for adoption.",
			terminal.WithErrorStyle(),
		)
		return 1
	}

	client := c.project.Client()
	s.Update("Finished connecting to: %s", c.serverUrl)
	s.Status(terminal.StatusOK)
	s.Done()

	// We generate the ID if the user doesn't provide one
	// This ID is used later to adopt the runner
	id := c.id
	if id == "" {
		id, err = server.Id()
		if err != nil {
			c.ui.Output("Error generating runner ID: %s", clierrors.Humanize(err), terminal.WithErrorStyle())
			return 1
		}
	}

	s = sg.Add("Installing runner...")
	err = p.Install(ctx, &runnerinstall.InstallOpts{
		Log:        log,
		UI:         c.ui,
		Cookie:     cookie,
		ServerAddr: c.serverUrl,
		AdvertiseClient: &serverconfig.Client{
			Address:       c.serverUrl,
			Tls:           c.serverTls,
			TlsSkipVerify: c.serverTlsSkipVerify,
			RequireAuth:   c.serverRequireAuth,
		},
		Id: id,
	})
	if err != nil {
		c.ui.Output("Error installing runner: %s", clierrors.Humanize(err),
			terminal.WithErrorStyle(),
		)
		return 1
	}
	s.Update("Runner %s installed successfully to %s", id, platform[0])
	s.Status(terminal.StatusOK)
	s.Done()

	if c.adopt {
		err = installutil.AdoptRunner(ctx, c.ui, client, id, c.serverUrl)
		if err != nil {
			c.ui.Output("Error adopting runner: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}

		// Creating a new runner profile for the newly adopted runner
		s = sg.Add("Creating runner profile and targeting runner %s", strings.ToUpper(id))
		runnerProfile, err := client.UpsertOnDemandRunnerConfig(ctx, &pb.UpsertOnDemandRunnerConfigRequest{
			Config: &pb.OnDemandRunnerConfig{
				Name: platform[0] + "-" + strings.ToUpper(id),
				TargetRunner: &pb.Ref_Runner{
					Target: &pb.Ref_Runner_Id{
						Id: &pb.Ref_RunnerId{
							Id: strings.ToUpper(id),
						},
					},
				},
				OciUrl:     c.runnerProfileOdrImage,
				PluginType: platform[0],
			},
		})
		if err != nil {
			c.ui.Output("Error creating runner profile: %s", clierrors.Humanize(err),
				terminal.WithErrorStyle(),
			)
			return 1
		}
		s.Update("Runner profile %s created successfully.", runnerProfile.Config.Name)
		s.Status(terminal.StatusOK)
		s.Done()
		return 0
	}
	return 0
}
