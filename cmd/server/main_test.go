package main

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/structpb"

	"buf.build/gen/go/a-novel/proto/grpc/go/passkeys/v1/passkeysv1grpc"
	passkeysv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/passkeys/v1"

	anovelgrpc "github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/testutils"
)

func init() {
	go main()
}

var servicesToTest = []string{
	"create",
	"delete",
	"get",
	"update",
}

func TestIntegrationHealth(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	// Create the RPC client.
	pool := anovelgrpc.NewConnPool()
	conn, err := pool.Open("0.0.0.0", 8080, anovelgrpc.ProtocolHTTP)
	require.NoError(t, err)

	healthClient := healthpb.NewHealthClient(conn)

	testutils.WaitConn(t, conn)

	for _, service := range servicesToTest {
		res, err := healthClient.Check(context.Background(), &healthpb.HealthCheckRequest{Service: service})
		require.NoError(t, err)
		require.Equal(t, healthpb.HealthCheckResponse_SERVING, res.Status)
	}
}

func TestIntegrationCRUD(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode.")
	}

	reward, err := structpb.NewStruct(map[string]interface{}{"type": "reward"})
	require.NoError(t, err)
	newReward, err := structpb.NewStruct(map[string]interface{}{"type": "new-reward"})
	require.NoError(t, err)

	// Create the RPC client.
	pool := anovelgrpc.NewConnPool()
	conn, err := pool.Open("0.0.0.0", 8080, anovelgrpc.ProtocolHTTP)
	require.NoError(t, err)

	testutils.WaitConn(t, conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	createPasskeyClient := passkeysv1grpc.NewCreateServiceClient(conn)
	deletePasskeyClient := passkeysv1grpc.NewDeleteServiceClient(conn)
	getPasskeyClient := passkeysv1grpc.NewGetServiceClient(conn)
	updatePasskeyClient := passkeysv1grpc.NewUpdateServiceClient(conn)

	// Create passkey
	createCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("password", "my-secret-password"))
	createData, err := createPasskeyClient.Exec(createCTX, &passkeysv1.CreateServiceExecRequest{
		Namespace: "namespace",
		Reward:    reward,
	})
	require.NoError(t, err)
	require.Equal(t, "namespace", createData.Namespace)
	require.Equal(t, reward.AsMap(), createData.Reward.AsMap())

	// Get passkey (no password)
	getData, err := getPasskeyClient.Exec(ctx, &passkeysv1.GetServiceExecRequest{
		Id:        createData.Id,
		Namespace: "namespace",
		Validate:  false,
	})
	require.NoError(t, err)
	require.Equal(t, "namespace", getData.Namespace)
	require.Equal(t, reward.AsMap(), getData.Reward.AsMap())

	// Get passkey (with password)
	getCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("password", "my-secret-password"))
	getData, err = getPasskeyClient.Exec(getCTX, &passkeysv1.GetServiceExecRequest{
		Id:        createData.Id,
		Namespace: "namespace",
		Validate:  true,
	})
	require.NoError(t, err)
	require.Equal(t, "namespace", getData.Namespace)
	require.Equal(t, reward.AsMap(), getData.Reward.AsMap())

	// Get passkey (with wrong password)
	getCTX = metadata.NewOutgoingContext(ctx, metadata.Pairs("password", "wrong-password"))
	_, err = getPasskeyClient.Exec(getCTX, &passkeysv1.GetServiceExecRequest{
		Id:        createData.Id,
		Namespace: "namespace",
		Validate:  true,
	})
	require.Error(t, err)
	testutils.RequireGRPCCodesEqual(t, err, codes.PermissionDenied)

	// Update passkey
	updateCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("password", "my-new-secret-password"))
	updateData, err := updatePasskeyClient.Exec(updateCTX, &passkeysv1.UpdateServiceExecRequest{
		Id:        createData.Id,
		Namespace: "namespace",
		Reward:    newReward,
	})
	require.NoError(t, err)
	require.Equal(t, "namespace", updateData.Namespace)
	require.Equal(t, newReward.AsMap(), updateData.Reward.AsMap())
	require.Equal(t, createData.Id, updateData.Id)

	// Try new passkey
	getCTX = metadata.NewOutgoingContext(ctx, metadata.Pairs("password", "my-new-secret-password"))
	getData, err = getPasskeyClient.Exec(getCTX, &passkeysv1.GetServiceExecRequest{
		Id:        createData.Id,
		Namespace: "namespace",
		Validate:  true,
	})
	require.NoError(t, err)
	require.Equal(t, "namespace", getData.Namespace)
	require.Equal(t, newReward.AsMap(), getData.Reward.AsMap())

	// Delete passkey
	deleteCTX := metadata.NewOutgoingContext(ctx, metadata.Pairs("password", "my-new-secret-password"))
	_, err = deletePasskeyClient.Exec(deleteCTX, &passkeysv1.DeleteServiceExecRequest{
		Id:        createData.Id,
		Namespace: "namespace",
		Validate:  true,
	})
	require.NoError(t, err)

	// Try to get deleted passkey
	_, err = getPasskeyClient.Exec(ctx, &passkeysv1.GetServiceExecRequest{
		Id:        createData.Id,
		Namespace: "namespace",
		Validate:  false,
	})
	require.Error(t, err)
	testutils.RequireGRPCCodesEqual(t, err, codes.NotFound)
}
