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
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	passkeysv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/passkeys/v1"

	adaptersmocks "github.com/a-novel/golib/loggers/adapters/mocks"
	"github.com/a-novel/golib/testutils"

	"github.com/a-novel/uservice-passkeys/pkg/dao"
	"github.com/a-novel/uservice-passkeys/pkg/handlers"
	"github.com/a-novel/uservice-passkeys/pkg/services"
	servicesmocks "github.com/a-novel/uservice-passkeys/pkg/services/mocks"
)

func TestGetPasskey(t *testing.T) {
	reward, err := structpb.NewStruct(map[string]interface{}{"type": "reward"})
	require.NoError(t, err)

	testCases := []struct {
		name string

		metadata map[string]string
		request  *passkeysv1.GetServiceExecRequest

		callServiceWith *services.GetPasskeyRequest
		serviceResp     *services.GetPasskeyResponse
		serviceErr      error

		expect     *passkeysv1.GetServiceExecResponse
		expectCode codes.Code
	}{
		{
			name: "OK",

			metadata: map[string]string{
				"password": "passkey",
			},
			request: &passkeysv1.GetServiceExecRequest{
				Id:        "id",
				Namespace: "namespace",
				Validate:  true,
			},

			callServiceWith: &services.GetPasskeyRequest{
				ID:        "id",
				Namespace: "namespace",
				Passkey:   "passkey",
				Validate:  true,
			},
			serviceResp: &services.GetPasskeyResponse{
				ID:        "id",
				Namespace: "namespace",
				Reward:    map[string]interface{}{"type": "reward"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},

			expect: &passkeysv1.GetServiceExecResponse{
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

			metadata: map[string]string{},
			request: &passkeysv1.GetServiceExecRequest{
				Id:        "id",
				Namespace: "namespace",
			},

			callServiceWith: &services.GetPasskeyRequest{
				ID:        "id",
				Namespace: "namespace",
			},
			serviceResp: &services.GetPasskeyResponse{
				ID:        "id",
				Namespace: "namespace",
				Reward:    map[string]interface{}{"type": "reward"},
				ExpiresAt: lo.ToPtr(time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: lo.ToPtr(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},

			expect: &passkeysv1.GetServiceExecResponse{
				Id:        "id",
				Namespace: "namespace",
				Reward:    reward,
				ExpiresAt: timestamppb.New(time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)),
				CreatedAt: timestamppb.New(time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)),
				UpdatedAt: timestamppb.New(time.Date(2021, 2, 1, 0, 0, 0, 0, time.UTC)),
			},
		},
		{
			name: "InvalidRequest",

			metadata: map[string]string{
				"password": "passkey",
			},
			request: &passkeysv1.GetServiceExecRequest{
				Id:        "id",
				Namespace: "namespace",
				Validate:  true,
			},

			callServiceWith: &services.GetPasskeyRequest{
				ID:        "id",
				Namespace: "namespace",
				Passkey:   "passkey",
				Validate:  true,
			},
			serviceErr: services.ErrInvalidGetPasskeyRequest,

			expectCode: codes.InvalidArgument,
		},
		{
			name: "InvalidPasskey",

			metadata: map[string]string{
				"password": "passkey",
			},
			request: &passkeysv1.GetServiceExecRequest{
				Id:        "id",
				Namespace: "namespace",
				Validate:  true,
			},

			callServiceWith: &services.GetPasskeyRequest{
				ID:        "id",
				Namespace: "namespace",
				Passkey:   "passkey",
				Validate:  true,
			},
			serviceErr: dao.ErrInvalidPasskey,

			expectCode: codes.PermissionDenied,
		},
		{
			name: "InternalError",

			metadata: map[string]string{
				"password": "passkey",
			},
			request: &passkeysv1.GetServiceExecRequest{
				Id:        "id",
				Namespace: "namespace",
				Validate:  true,
			},

			callServiceWith: &services.GetPasskeyRequest{
				ID:        "id",
				Namespace: "namespace",
				Passkey:   "passkey",
				Validate:  true,
			},
			serviceErr: errors.New("uwups"),

			expectCode: codes.Internal,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			service := servicesmocks.NewMockGetPasskey(t)
			logger := adaptersmocks.NewMockGRPC(t)

			ctx := metadata.NewIncomingContext(context.Background(), metadata.New(testCase.metadata))

			service.
				On("Exec", ctx, testCase.callServiceWith).
				Return(testCase.serviceResp, testCase.serviceErr)

			logger.On("Report", handlers.GetPasskeyServiceName, mock.Anything)

			handler := handlers.NewGetPasskey(service, logger)
			resp, err := handler.Exec(ctx, testCase.request)

			testutils.RequireGRPCCodesEqual(t, err, testCase.expectCode)
			require.Equal(t, testCase.expect, resp)

			service.AssertExpectations(t)
			logger.AssertExpectations(t)
		})
	}
}
