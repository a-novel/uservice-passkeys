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

const UpdatePasskeyServiceName = "update_passkey"

type UpdatePasskey interface {
	passkeysv1grpc.UpdateServiceServer
}

type updatePasskeyImpl struct {
	service services.UpdatePasskey
}

var handleUpdatePasskeyError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidUpdatePasskeyRequest, codes.InvalidArgument).
	Handle

func (handler *updatePasskeyImpl) Exec(
	ctx context.Context, request *passkeysv1.UpdateServiceExecRequest,
) (*passkeysv1.UpdateServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.UpdatePasskeyRequest{
		ID:        request.GetId(),
		Namespace: request.GetNamespace(),
		Passkey:   ExtractPasskey(ctx),
		Reward:    grpc.StructOptionalProto(request.GetReward()),
		ExpiresIn: grpc.DurationOptionalProto(request.GetExpiresIn()),
	})
	if err != nil {
		return nil, handleUpdatePasskeyError(err)
	}

	reward, err := grpc.StructOptional(res.Reward)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "convert reward: %v", err)
	}

	return &passkeysv1.UpdateServiceExecResponse{
		Id:        res.ID,
		Namespace: res.Namespace,
		Reward:    reward,
		ExpiresAt: grpc.TimestampOptional(res.ExpiresAt),
		CreatedAt: timestamppb.New(res.CreatedAt),
		UpdatedAt: grpc.TimestampOptional(res.UpdatedAt),
	}, nil
}

func NewUpdatePasskey(service services.UpdatePasskey, logger adapters.GRPC) UpdatePasskey {
	handler := &updatePasskeyImpl{service: service}
	return grpc.ServiceWithMetrics(UpdatePasskeyServiceName, handler, logger)
}
