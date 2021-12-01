package singleprocess

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/hashicorp/waypoint/internal/server"
	pb "github.com/hashicorp/waypoint/internal/server/gen"
	serverptypes "github.com/hashicorp/waypoint/internal/server/ptypes"
)

func TestServiceTrigger(t *testing.T) {
	ctx := context.Background()

	// Create our server
	impl, err := New(WithDB(testDB(t)))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	type Req = pb.UpsertTriggerRequest

	t.Run("create and update", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: serverptypes.TestValidTrigger(t, nil),
		})
		require.NoError(err)
		require.NotNil(resp)
		result := resp.Trigger
		require.NotEmpty(result.Id)

		// Let's write some data
		testName := "TestyTest"
		result.Name = testName
		resp, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
			Trigger: result,
		})
		require.NoError(err)
		require.NotNil(resp)
		result = resp.Trigger
		require.Equal(result.Name, testName)
	})

	t.Run("update non-existent", func(t *testing.T) {
		require := require.New(t)

		// Create, should get an ID back
		resp, err := client.UpsertTrigger(ctx, &Req{
			Trigger: serverptypes.TestValidTrigger(t, &pb.Trigger{
				Id: "nope",
			}),
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceTrigger_GetTrigger(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	triggerId := resp.Trigger.Id

	type Req = pb.UpsertTriggerRequest

	t.Run("get existing", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a status report
		resp, err := client.GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: resp.Trigger.Id,
			},
		})
		require.NoError(err)
		require.NotNil(resp.Trigger)
		require.NotEmpty(resp.Trigger.Id)
		require.Equal(triggerId, resp.Trigger.Id)
	})

	t.Run("get non-existing", func(t *testing.T) {
		require := require.New(t)

		// get, should fail
		resp, err := client.GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: "nope",
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})
}

func TestServiceTrigger_ListTriggersSimple(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	_, err = client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})

	t.Run("list", func(t *testing.T) {
		require := require.New(t)

		respList, err := client.ListTriggers(ctx, &pb.ListTriggerRequest{})
		require.NoError(err)
		require.Equal(len(respList.Triggers), 3)
	})
}

func TestServiceTrigger_ListTriggersWithFilters(t *testing.T) {
	//ctx := context.Background()

	//// Create our server
	//db := testDB(t)
	//impl, err := New(WithDB(db))
	//require.NoError(t, err)
	//client := server.TestServer(t, impl)

	// TODO:
	// create some triggers across workspaces, projects, apps, labels
	// test that listing with filters returns expected number of triggers
}

func TestServiceTrigger_DeleteTrigger(t *testing.T) {
	ctx := context.Background()

	// Create our server
	db := testDB(t)
	impl, err := New(WithDB(db))
	require.NoError(t, err)
	client := server.TestServer(t, impl)

	resp, err := client.UpsertTrigger(ctx, &pb.UpsertTriggerRequest{
		Trigger: serverptypes.TestValidTrigger(t, nil),
	})
	triggerId := resp.Trigger.Id

	type Req = pb.UpsertTriggerRequest

	t.Run("get existing then delete", func(t *testing.T) {
		require := require.New(t)

		// Get, should return a status report
		resp, err := client.GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: resp.Trigger.Id,
			},
		})
		require.NoError(err)
		require.NotNil(resp.Trigger)
		require.NotEmpty(resp.Trigger.Id)
		require.Equal(triggerId, resp.Trigger.Id)

		_, err = client.DeleteTrigger(ctx, &pb.DeleteTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: triggerId,
			},
		})
		require.NoError(err)

		// get, should fail
		resp, err = client.GetTrigger(ctx, &pb.GetTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: triggerId,
			},
		})
		require.Error(err)
		require.Nil(resp)
		st, ok := status.FromError(err)
		require.True(ok)
		require.Equal(codes.NotFound, st.Code())
	})

	t.Run("delete non-existing", func(t *testing.T) {
		require := require.New(t)

		resp, err := client.DeleteTrigger(ctx, &pb.DeleteTriggerRequest{
			Ref: &pb.Ref_Trigger{
				Id: "nope",
			},
		})
		require.NoError(err)
		require.NotNil(resp)
	})
}
