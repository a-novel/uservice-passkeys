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

	"github.com/a-novel/uservice-passkeys/pkg/dao"
	"github.com/a-novel/uservice-passkeys/pkg/services"
)

const DeletePasskeyServiceName = "delete_passkey"

type DeletePasskey interface {
	passkeysv1grpc.DeleteServiceServer
}

type deletePasskeyImpl struct {
	service services.DeletePasskey
}

var handleDeletePasskeyError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidDeletePasskeyRequest, codes.InvalidArgument).
	Is(dao.ErrPasskeyNotFound, codes.NotFound).
	Is(dao.ErrInvalidPasskey, codes.PermissionDenied).
	Handle

func (handler *deletePasskeyImpl) Exec(
	ctx context.Context, request *passkeysv1.DeleteServiceExecRequest,
) (*passkeysv1.DeleteServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.DeletePasskeyRequest{
		ID:        request.GetId(),
		Namespace: request.GetNamespace(),
		Passkey:   ExtractPasskey(ctx),
		Validate:  request.GetValidate(),
	})
	if err != nil {
		return nil, handleDeletePasskeyError(err)
	}

	reward, err := grpc.StructOptional(res.Reward)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "convert reward: %v", err)
	}

	return &passkeysv1.DeleteServiceExecResponse{
		Id:        res.ID,
		Namespace: res.Namespace,
		Reward:    reward,
		ExpiresAt: grpc.TimestampOptional(res.ExpiresAt),
		CreatedAt: timestamppb.New(res.CreatedAt),
		UpdatedAt: grpc.TimestampOptional(res.UpdatedAt),
	}, nil
}

func NewDeletePasskey(service services.DeletePasskey, logger adapters.GRPC) DeletePasskey {
	handler := &deletePasskeyImpl{service: service}
	return grpc.ServiceWithMetrics(DeletePasskeyServiceName, handler, logger)
}
