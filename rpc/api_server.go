package rpc

import (
	"context"
	"fmt"

	"github.com/bytom/rpc/pb"
	"github.com/bytom/wallet"
	"github.com/bytom/protocol"
	"github.com/bytom/util"
	"github.com/bytom/crypto/ed25519/chainkd"

	empty "github.com/golang/protobuf/ptypes/empty"
)

type ApiService struct {
	rpcServer *Rpc

	wallet        *wallet.Wallet
	chain         *protocol.Chain
}

func (s *ApiService)GetState(ctx context.Context, req *rpcpb.NonParamsRequest) (*rpcpb.GetStateResponse, error) {
	return &rpcpb.GetStateResponse{Status:"OK"}, nil
}

func (s *ApiService)CreateKey(ctx context.Context, req *rpcpb.CreateKeyRequest) (*rpcpb.CreateKeyResponse, error) {
	xpub, err := s.wallet.Hsm.XCreate(req.Alias, req.Password)
	if err != nil {
		return nil, err
	}

	xpubBytes := util.JsonEncode(xpub)
	return &rpcpb.CreateKeyResponse{Xpub:string(xpubBytes)}, nil
}

func (s *ApiService)ListKey(ctx context.Context, req *rpcpb.NonParamsRequest) (*rpcpb.ListKeyResponse, error) {
	var xpubs []string
	for _, xpub := range s.wallet.Hsm.ListKeys() {
		xpubs = append(xpubs, string(util.JsonEncode(xpub)))
	}
	return &rpcpb.ListKeyResponse{Xpub:xpubs}, nil
}

func (s *ApiService)DeleteKey(ctx context.Context, req *rpcpb.DeleteKeyRequest) (*empty.Empty, error) {
	xpub := new(chainkd.XPub)
	if err := xpub.UnmarshalText([]byte(req.Xpub)); err != nil {
		return nil, fmt.Errorf("delete-key: %v", err.Error())
	}
	if err := s.wallet.Hsm.XDelete(*xpub, req.Password); err != nil {
		return  nil, err
	}
	return &empty.Empty{}, nil
}