package handlers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	passkeysv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/passkeys/v1"

	adaptersmocks "github.com/a-novel/golib/loggers/adapters/mocks"
	"github.com/a-novel/golib/testutils"

	"github.com/a-novel/uservice-passkeys/pkg/handlers"
	"github.com/a-novel/uservice-passkeys/pkg/services"
	servicesmocks "github.com/a-novel/uservice-passkeys/pkg/services/mocks"
)

func TestUpdatePasskey(t *testing.T) {
	reward, err := structpb.NewStruct(map[string]interface{}{"type": "reward"})
	require.NoError(t, err)

	testCases := []struct {
		name string

		metadata map[string]string
		request  *passkeysv1.UpdateServiceExecRequest

		callServiceWith *services.UpdatePasskeyRequest
		serviceResp     *services.UpdatePasskeyResponse
		serviceErr      error

		expect     *passkeysv1.UpdateServiceExecResponse
		expectCode codes.Code
	}{
		{
			name: "OK",

			metadata: map[string]string{
				"password": "passkey",
			},
			request: &passkeysv1.UpdateServiceExecRequest{
				Namespace: "namespace",
				Reward:    reward,
				ExpiresIn: durationpb.New(5 * time.Minute),
			},

			callServiceWith: &services.UpdatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
				Reward:    map[string]interface{}{"type": "reward"},
				ExpiresIn: lo.ToPtr(5 * time.Minute),
			},
			serviceResp: &services.UpdatePasskeyResponse{
				ID:        "id",
				Namespace: "namespace",
				Reward:    map[string]interface{}{"type": "reward"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},

			expect: &passkeysv1.UpdateServiceExecResponse{
				Id:        "id",
				Namespace: "namespace",
				Reward:    reward,
				ExpiresAt: timestamppb.New(time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: timestamppb.New(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "OK/Minimal",

			metadata: map[string]string{
				"password": "passkey",
			},
			request: &passkeysv1.UpdateServiceExecRequest{
				Namespace: "namespace",
			},

			callServiceWith: &services.UpdatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
			},
			serviceResp: &services.UpdatePasskeyResponse{
				ID:        "id",
				Namespace: "namespace",
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},

			expect: &passkeysv1.UpdateServiceExecResponse{
				Id:        "id",
				Namespace: "namespace",
				CreatedAt: timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "InvalidRequest",

			metadata: map[string]string{
				"password": "passkey",
			},

			request: &passkeysv1.UpdateServiceExecRequest{
				Namespace: "namespace",
			},

			callServiceWith: &services.UpdatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
			},

			serviceErr: services.ErrInvalidUpdatePasskeyRequest,

			expectCode: codes.InvalidArgument,
		},
		{
			name: "InternalError",

			metadata: map[string]string{
				"password": "passkey",
			},

			request: &passkeysv1.UpdateServiceExecRequest{
				Namespace: "namespace",
			},

			callServiceWith: &services.UpdatePasskeyRequest{
				Namespace: "namespace",
				Passkey:   "passkey",
			},

			serviceErr: errors.New("uwups"),

			expectCode: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := servicesmocks.NewMockUpdatePasskey(t)
			logger := adaptersmocks.NewMockGRPC(t)

			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(testCase.metadata))

			service.
				On("Exec", ctx, testCase.callServiceWith).
				Return(testCase.serviceResp, testCase.serviceErr)

			logger.On("Report", handlers.UpdatePasskeyServiceName, mock.Anything)

			handler := handlers.NewUpdatePasskey(service, logger)
			resp, err := handler.Exec(ctx, testCase.request)

			testutils.RequireGRPCCodesEqual(t, err, testCase.expectCode)
			require.Equal(t, testCase.expect, resp)

			service.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
