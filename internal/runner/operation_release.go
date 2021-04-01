package runner

import (
	"context"

	"github.com/hashicorp/go-hclog"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
	"github.com/hashicorp/waypoint/internal/core"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
)

func (r *Runner) executeReleaseOp(
	ctx context.Context,
	log hclog.Logger,
	job *pb.Job,
	project *core.Project,
) (*pb.Job_Result, error) {
	app, err := project.App(job.Application.Application)
	if err != nil {
		return nil, err
	}

	op, ok := job.Operation.(*pb.Job_Release)
	if !ok {
		// this shouldn't happen since the call to this function is gated
		// on the above type match.
		panic("operation not expected type")
	}

	// Our target deployment
	target := op.Release.Deployment

	// Get our last release. If its the same generation, then release is
	// a no-op and return this value. We only do this if we have a generation.
	// We SHOULD but if we have an old client, its possible we don't.
	var release *pb.Release
	if target.Generation != nil {
		resp, err := r.client.ListReleases(ctx, &pb.ListReleasesRequest{
			Application:   app.Ref(),
			Workspace:     project.WorkspaceRef(),
			PhysicalState: pb.Operation_CREATED,
			LoadDetails:   pb.Release_DEPLOYMENT,
			Order: &pb.OperationOrder{
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
		})
		if err != nil {
			return nil, err
		}
		for _, r := range resp.Releases {
			if r.Preload != nil && r.Preload.Deployment != nil {
				d := r.Preload.Deployment
				if d.Generation != nil && d.Generation.Id == target.Generation.Id {
					release = r
				}
			}

			// We always break cause we really only want the first release.
			// This is the most recent release.
			break
		}
	}

	// If we're pruning, then let's query the deployments we want to prune
	// ahead of time so that fails fast.
	var pruneDeploys []*pb.Deployment
	if op.Release.Prune {
		// Determine the number of deployments to keep around.
		retain := 2
		if op.Release.PruneRetainOverride {
			retain = int(op.Release.PruneRetain) + 1 // add 1 to make this the total number
		}

		log.Debug("pruning requested, gathering deployments to prune",
			"retain", retain)
		resp, err := r.client.ListDeployments(ctx, &pb.ListDeploymentsRequest{
			Application:   app.Ref(),
			Workspace:     project.WorkspaceRef(),
			PhysicalState: pb.Operation_CREATED,
			Order: &pb.OperationOrder{
				Order: pb.OperationOrder_COMPLETE_TIME,
				Desc:  true,
			},
		})
		if err != nil {
			return nil, err
		}

		// If we have less than the prune amount, then we do nothing. Otherwise
		// we prune away the ones we're definitely keeping.
		if len(resp.Deployments) <= retain {
			log.Debug("less than the limit deployments exists, no pruning")
			resp.Deployments = nil
		} else {
			resp.Deployments = resp.Deployments[retain:]
		}

		// Assign to short character var since we'll manipulate it a lot
		ds := make([]*pb.Deployment, 0, len(resp.Deployments))
		for _, d := range resp.Deployments {
			// If this is the deployment we're releasing, then do NOT delete it.
			if d.Id == op.Release.Deployment.Id {
				continue
			}

			// Ignore deployments with the same generation, because this
			// means that they share underlying resources.
			if target.Generation != nil && d.Generation != nil {
				if target.Generation.Id == d.Generation.Id {
					continue
				}
			}

			// TODO this should instead check against the app's platform component
			// and ignore any deployments that are NOT the app's current platform
			// component (ya dig?)
			if d.Component.Name != op.Release.Deployment.Component.Name {
				continue
			}

			// Mark for deletion
			ds = append(ds, d)
		}

		log.Info("will prune deploys", "len", len(ds))
		pruneDeploys = ds
	}

	// Do the release
	if release == nil {
		release, _, err = app.Release(ctx, op.Release.Deployment)
		if err != nil {
			return nil, err
		}
	} else {
		log.Info("not releasing since last released deploy has a matching generation",
			"gen", target.Generation.Id)
	}

	result := &pb.Job_Result{
		Release: &pb.Job_ReleaseResult{
			Release: release,
		},
	}

	// Do the pruning
	if len(pruneDeploys) > 0 {
		log.Info("pruning deploys", "len", len(pruneDeploys))
		app.UI.Output("Pruning old deployments...", terminal.WithHeaderStyle())
		for _, d := range pruneDeploys {
			app.UI.Output("Deployment: %s (v%d)", d.Id, d.Sequence, terminal.WithInfoStyle())
			if err := app.DestroyDeploy(ctx, d); err != nil {
				return result, err
			}
		}
	}

	return result, nil
}
