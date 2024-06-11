package dto

import (
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/actors/builtin/miner"
)

//miner 信息
type MinerInfo struct{
	Miner string
	Country   string
	City      string
	Latitude  float64
	Longitude float64
	MinerInfo *miner.MinerInfo
	MinerPower *api.MinerPower
}
