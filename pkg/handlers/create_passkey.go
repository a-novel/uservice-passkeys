package handlers

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"buf.build/gen/go/a-novel/proto/grpc/go/passkeys/v1/passkeysv1grpc"
	passkeysv1 "buf.build/gen/go/a-novel/proto/protocolbuffers/go/passkeys/v1"

	"github.com/a-novel/golib/grpc"
	"github.com/a-novel/golib/loggers/adapters"

	"github.com/a-novel/uservice-passkeys/pkg/services"
)

const CreatePasskeyServiceName = "create_passkey"

type CreatePasskey interface {
	passkeysv1grpc.CreateServiceServer
}

type createPasskeyImpl struct {
	service services.CreatePasskey
}

var handleCreatePasskeyError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidCreatePasskeyRequest, codes.InvalidArgument).
	Handle

func (handler *createPasskeyImpl) Exec(
	ctx context.Context, request *passkeysv1.CreateServiceExecRequest,
) (*passkeysv1.CreateServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.CreatePasskeyRequest{
		Namespace: request.GetNamespace(),
		Passkey:   ExtractPasskey(ctx),
		Reward:    grpc.StructOptionalProto(request.GetReward()),
		ExpiresIn: grpc.DurationOptionalProto(request.GetExpiresIn()),
	})
	if err != nil {
		return nil, handleCreatePasskeyError(err)
	}

	reward, err := grpc.StructOptional(res.Reward)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "convert reward: %v", err)
	}

	return &passkeysv1.CreateServiceExecResponse{
		Id:        res.ID,
		Namespace: res.Namespace,
		Reward:    reward,
		ExpiresAt: grpc.TimestampOptional(res.ExpiresAt),
		CreatedAt: timestamppb.New(res.CreatedAt),
	}, nil
}

func NewCreatePasskey(service services.CreatePasskey, logger adapters.GRPC) CreatePasskey {
	handler := &createPasskeyImpl{service: service}
	return grpc.ServiceWithMetrics(CreatePasskeyServiceName, handler, logger)
}
