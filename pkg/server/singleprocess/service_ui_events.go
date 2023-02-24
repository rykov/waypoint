package singleprocess
import (
	"context"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/waypoint/pkg/server/hcerr"
	serverptypes "github.com/hashicorp/waypoint/pkg/server/ptypes"

	pb "github.com/hashicorp/waypoint/pkg/server/gen"
)

func (s *Service) UI_ListEvents(
	ctx context.Context,
	req *pb.UI_ListEventsRequest,
) (*pb.UI_ListEventsResponse, error) {
	if err := serverptypes.ValidateUIListEventsRequest(req); err != nil {
		return nil, err
	}

	eventBundles, pagination, err := s.state(ctx).EventListBundles(ctx, req.Pagination)
	if err != nil {
		return nil, hcerr.Externalize(
			hclog.FromContext(ctx),
			err,
			"failed to list events",
		)
	}
	return &pb.UI_ListEventsResponse{
		Events:     eventBundles,
		Pagination: pagination,
	}, nil
}