package keeper

import (
	"context"
	"errors"

	storetypes "cosmossdk.io/store/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v8/modules/core/exported"
	ibctm "github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint"
)

func (k *Keeper) SetPostUpdateHandler(postUpdateHandler func(context.Context, string, int64, *cmtproto.ValidatorSet) error) {
	k.postUpdateHandler = postUpdateHandler
}

func (k Keeper) handlePostUpdate(ctx sdk.Context, clientState exported.ClientState, clientStore storetypes.KVStore, clientMsg exported.ClientMessage) error {
	if k.postUpdateHandler == nil || clientState.ClientType() != exported.Tendermint {
		return errors.New("not set post handler")
	}

	header := clientMsg.(*ibctm.Header)
	tmState := clientState.(*ibctm.ClientState)

	return k.postUpdateHandler(ctx, tmState.ChainId, header.Header.Height, header.ValidatorSet)
}
