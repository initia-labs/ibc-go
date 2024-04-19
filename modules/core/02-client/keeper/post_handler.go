package keeper

import (
	"context"

	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

func (k *Keeper) SetPostUpdateHandler(postUpdateHandler func(context.Context, string, int64, *cmtproto.ValidatorSet) error) {
	k.postUpdateHandler = postUpdateHandler
}

func (k Keeper) handlePostUpdate(ctx sdk.Context, clientID string, clientState exported.ClientState, clientMsg exported.ClientMessage) error {
	// ignore if the handler is not set or the client is not a tendermint client
	if k.postUpdateHandler == nil || clientState.ClientType() != exported.Tendermint {
		return nil
	}

	header := clientMsg.(*ibctm.Header)
	return k.postUpdateHandler(ctx, clientID, header.Header.Height, header.ValidatorSet)
}
