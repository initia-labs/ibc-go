package tendermintattestor

import (
	"bytes"
	"slices"

	sdk "github.com/cosmos/cosmos-sdk/types"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/ibc-go/v8/modules/light-clients/07-tendermint-attestor/keys"
)

func (cs ClientState) VerifySignatures(
	ctx sdk.Context,
	proofBytes []byte,
	attestations []*Attestation,
) error {
	if cs.Threshold == 0 {
		return nil
	} else if len(attestations) < int(cs.Threshold) {
		return errorsmod.Wrapf(ErrUnauthorizedAttestation, "not enough attestations: %d < %d", len(attestations), cs.Threshold)
	}
	for _, attestation := range attestations {
		var pubKey keys.PubKey
		if index := slices.IndexFunc(cs.AttestorPubkeys, func(attestorPubkey PubKey) bool {
			return bytes.Equal(attestorPubkey.Key, attestation.PubKey)
		}); index >= 0 {
			pubKey = convertPubKey(cs.AttestorPubkeys[index])
			if pubKey == nil {
				return errorsmod.Wrapf(ErrUnauthorizedAttestation, "invalid pubkey type: %s", cs.AttestorPubkeys[index])
			}
		} else {
			return ErrUnauthorizedAttestation
		}

		if err := pubKey.Verify(proofBytes, attestation.Signature); err != nil {
			return errorsmod.Wrap(ErrInvalidAttestation, err.Error())
		}
	}
	return nil
}

func convertPubKey(pubKey PubKey) keys.PubKey {
	switch pubKey.Type {
	case PubKeyType_SECP256K1:
		return keys.NewSecp256K1PubKey(pubKey.Key)
	case PubKeyType_ETH_SECP256K1:
		return keys.NewEthSecp256K1PubKey(pubKey.Key)
	}
	return nil
}
