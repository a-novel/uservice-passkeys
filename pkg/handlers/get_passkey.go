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

const GetPasskeyServiceName = "get_passkey"

type GetPasskey interface {
	passkeysv1grpc.GetServiceServer
}

type getPasskeyImpl struct {
	service services.GetPasskey
}

var handleGetPasskeyError = grpc.HandleError(codes.Internal).
	Is(services.ErrInvalidGetPasskeyRequest, codes.InvalidArgument).
	Is(dao.ErrPasskeyNotFound, codes.NotFound).
	Is(dao.ErrInvalidPasskey, codes.PermissionDenied).
	Handle

func (handler *getPasskeyImpl) Exec(
	ctx context.Context, request *passkeysv1.GetServiceExecRequest,
) (*passkeysv1.GetServiceExecResponse, error) {
	res, err := handler.service.Exec(ctx, &services.GetPasskeyRequest{
		ID:        request.GetId(),
		Namespace: request.GetNamespace(),
		Passkey:   ExtractPasskey(ctx),
		Validate:  request.GetValidate(),
	})
	if err != nil {
		return nil, handleGetPasskeyError(err)
	}

	reward, err := grpc.StructOptional(res.Reward)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "convert reward: %v", err)
	}

	return &passkeysv1.GetServiceExecResponse{
		Id:        res.ID,
		Namespace: res.Namespace,
		Reward:    reward,
		ExpiresAt: grpc.TimestampOptional(res.ExpiresAt),
		CreatedAt: timestamppb.New(res.CreatedAt),
		UpdatedAt: grpc.TimestampOptional(res.UpdatedAt),
	}, nil
}

func NewGetPasskey(service services.GetPasskey, logger adapters.GRPC) GetPasskey {
	handler := &getPasskeyImpl{service: service}
	return grpc.ServiceWithMetrics(GetPasskeyServiceName, handler, logger)
}
