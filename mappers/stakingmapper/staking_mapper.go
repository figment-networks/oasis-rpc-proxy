package stakingmapper

import (
	"github.com/figment-networks/oasis-rpc-proxy/grpc/state/statepb"
	"github.com/figment-networks/oasis-rpc-proxy/mappers/accountmapper"
	"github.com/oasislabs/oasis-core/go/staking/api"
)

func ToPb(rawStaking api.Genesis) *statepb.Staking {
	// Thresholds
	thresholds := map[int64][]byte{}
	for kind, quantity := range rawStaking.Parameters.Thresholds {
		thresholds[int64(kind)] = quantity.ToBigInt().Bytes()
	}

	// Reward Schedule
	var rewardSchedule []*statepb.RewardStep
	for _, step := range rawStaking.Parameters.RewardSchedule {
		rewardSchedule = append(rewardSchedule, &statepb.RewardStep{
			Scale: step.Scale.ToBigInt().Bytes(),
			Until: uint64(step.Until),
		})
	}

	// Slashing
	slashing := map[string]*statepb.Slash{}
	for reason, slash := range rawStaking.Parameters.Slashing {
		slashing[reason.String()] = &statepb.Slash{
			Amount:         slash.Amount.ToBigInt().Bytes(),
			FreezeInterval: uint64(slash.FreezeInterval),
		}
	}

	// Gas costs
	gasCosts := map[string]uint64{}
	for op, gas := range rawStaking.Parameters.GasCosts {
		gasCosts[string(op)] = uint64(gas)
	}

	// Undisable transfers from
	undisableTransfersFrom := map[string]bool{}
	for key, b := range rawStaking.Parameters.UndisableTransfersFrom {
		undisableTransfersFrom[key.String()] = b
	}

	// Ledger
	ledger := map[string]*statepb.Account{}
	for key, account := range rawStaking.Ledger {
		ledger[key.String()] = accountmapper.ToPb(*account)
	}

	// Delegations
	delegations := map[string]*statepb.DelegationEntry{}
	for validatorId, items := range rawStaking.Delegations {
		entries := map[string]*statepb.Delegation{}
		for escrowPublicKey, delegation := range items {
			entries[escrowPublicKey.String()] = &statepb.Delegation{
				Shares: delegation.Shares.ToBigInt().Bytes(),
			}
		}

		delegations[validatorId.String()] = &statepb.DelegationEntry{
			Entries: entries,
		}
	}

	// Debonding delegations
	debondingDelegations := map[string]*statepb.DebondingDelegationEntry{}
	for validatorId, items := range rawStaking.DebondingDelegations {
		innerEntries := map[string]*statepb.DebondingDelegationInnerEntry{}
		for escrowPublicKey, innerItems := range items {
			var dds []*statepb.DebondingDelegation
			for _, item := range innerItems {
				dds = append(dds, &statepb.DebondingDelegation{
					Shares:        item.Shares.ToBigInt().Bytes(),
					DebondEndTime: uint64(item.DebondEndTime),
				})
			}

			innerEntries[escrowPublicKey.String()] = &statepb.DebondingDelegationInnerEntry{
				DebondingDelegations: dds,
			}
		}

		debondingDelegations[validatorId.String()] = &statepb.DebondingDelegationEntry{
			Entries: innerEntries,
		}
	}

	return &statepb.Staking{
		TotalSupply: rawStaking.TotalSupply.ToBigInt().Bytes(),
		CommonPool:  rawStaking.CommonPool.ToBigInt().Bytes(),
		Parameters: &statepb.StakingParameters{
			Thresholds:                        thresholds,
			DebondingInterval:                 uint64(rawStaking.Parameters.DebondingInterval),
			RewardSchedule:                    rewardSchedule,
			SigningRewardThresholdNumerator:   rawStaking.Parameters.SigningRewardThresholdNumerator,
			SigningRewardThresholdDenominator: rawStaking.Parameters.SigningRewardThresholdDenominator,
			CommissionScheduleRules: &statepb.CommissionScheduleRules{
				RateBoundLead:      uint64(rawStaking.Parameters.CommissionScheduleRules.RateBoundLead),
				RateChangeInterval: uint64(rawStaking.Parameters.CommissionScheduleRules.RateChangeInterval),
				MaxBoundSteps:      int64(rawStaking.Parameters.CommissionScheduleRules.MaxBoundSteps),
				MaxRateSteps:       int64(rawStaking.Parameters.CommissionScheduleRules.MaxRateSteps),
			},
			Slashing:                  slashing,
			GasCosts:                  gasCosts,
			MinDelegationAmount:       rawStaking.Parameters.MinDelegationAmount.ToBigInt().Bytes(),
			DisableTransfers:          rawStaking.Parameters.DisableTransfers,
			DisableDelegation:         rawStaking.Parameters.DisableDelegation,
			UndisableTransfersFrom:    undisableTransfersFrom,
			FeeSplitVote:              rawStaking.Parameters.FeeSplitVote.ToBigInt().Bytes(),
			FeeSplitPropose:           rawStaking.Parameters.FeeSplitPropose.ToBigInt().Bytes(),
			RewardFactorEpochSigned:   rawStaking.Parameters.RewardFactorEpochSigned.ToBigInt().Bytes(),
			RewardFactorBlockProposed: rawStaking.Parameters.RewardFactorBlockProposed.ToBigInt().Bytes(),
		},
		Ledger:               ledger,
		Delegations:          delegations,
		DebondingDelegations: debondingDelegations,
	}
}