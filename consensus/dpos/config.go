// Copyright 2018 The Fractal Team Authors
// This file is part of the fractal project.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package dpos

import (
	"fmt"
	"math/big"
	"sync/atomic"
	"time"

	"github.com/fractalplatform/fractal/utils/fdb"
	"github.com/fractalplatform/fractal/utils/rlp"
)

// DefaultConfig configures
var DefaultConfig = &Config{
	MaxURLLen:            512,
	UnitStake:            big.NewInt(1000),
	ProducerMinQuantity:  big.NewInt(10),
	VoterMinQuantity:     big.NewInt(1),
	ActivatedMinQuantity: big.NewInt(100),
	BlockInterval:        3000,
	BlockFrequency:       6,
	ProducerScheduleSize: 3,
	DelayEcho:            2,
	AccountName:          "ftsystemdpos",
	SystemName:           "ftsystemio",
	SystemURL:            "www.fractalproject.com",
	ExtraBlockReward:     big.NewInt(1),
	BlockReward:          big.NewInt(5),
	Decimals:             18,
}

// Config dpos configures
type Config struct {
	// consensus fileds
	MaxURLLen            uint64   // url length
	UnitStake            *big.Int // state unit
	ProducerMinQuantity  *big.Int // min quantity
	VoterMinQuantity     *big.Int // min quantity
	ActivatedMinQuantity *big.Int // min active quantity
	BlockInterval        uint64
	BlockFrequency       uint64
	ProducerScheduleSize uint64
	DelayEcho            uint64
	AccountName          string
	SystemName           string
	SystemURL            string
	ExtraBlockReward     *big.Int
	BlockReward          *big.Int
	Decimals             uint64

	// cache files
	decimal    atomic.Value
	blockInter atomic.Value
	epochInter atomic.Value
	safeSize   atomic.Value
}

// EncodeRLP  encoding the consensus fileds.
func (cfg *Config) EncodeRLP() ([]byte, error) {
	return rlp.EncodeToBytes(cfg)
}

// DecodeRLP decoding the consensus fields.
func (cfg *Config) DecodeRLP(data []byte) error {
	return rlp.DecodeBytes(data, &cfg)
}

func (cfg *Config) decimals() *big.Int {
	if decimal := cfg.decimal.Load(); decimal != nil {
		return decimal.(*big.Int)
	}
	decimal := big.NewInt(1)
	for i := uint64(0); i < cfg.Decimals; i++ {
		decimal = new(big.Int).Mul(decimal, big.NewInt(10))
	}
	cfg.decimal.Store(decimal)
	return decimal
}

func (cfg *Config) unitStake() *big.Int {
	return new(big.Int).Mul(cfg.UnitStake, cfg.decimals())
}

func (cfg *Config) extraBlockReward() *big.Int {
	return new(big.Int).Mul(cfg.ExtraBlockReward, cfg.decimals())
}
func (cfg *Config) blockReward() *big.Int {
	return new(big.Int).Mul(cfg.BlockReward, cfg.decimals())
}
func (cfg *Config) blockInterval() uint64 {
	if blockInter := cfg.blockInter.Load(); blockInter != nil {
		return blockInter.(uint64)
	}
	blockInter := cfg.BlockInterval * uint64(time.Millisecond)
	cfg.blockInter.Store(blockInter)
	return blockInter
}

func (cfg *Config) epochInterval() uint64 {
	if epochInter := cfg.epochInter.Load(); epochInter != nil {
		return epochInter.(uint64)
	}
	epochInter := cfg.blockInterval() * cfg.BlockFrequency * cfg.ProducerScheduleSize
	cfg.epochInter.Store(epochInter)
	return epochInter
}

func (cfg *Config) consensusSize() uint64 {
	if safeSize := cfg.safeSize.Load(); safeSize != nil {
		return safeSize.(uint64)
	}

	safeSize := cfg.ProducerScheduleSize*2/3 + 1
	cfg.safeSize.Store(safeSize)
	return safeSize
}

func (cfg *Config) slot(timestamp uint64) uint64 {
	return ((timestamp + cfg.blockInterval()/10) / cfg.blockInterval() * cfg.blockInterval())
}

func (cfg *Config) nextslot(timestamp uint64) uint64 {
	return cfg.slot(timestamp) + cfg.blockInterval()
}

func (cfg *Config) getoffset(timestamp uint64) uint64 {
	offset := uint64(timestamp) % cfg.epochInterval()
	offset /= cfg.blockInterval() * cfg.BlockFrequency
	return offset
}

func (cfg *Config) epoch(timestamp uint64) uint64 {
	return timestamp / cfg.epochInterval()
}

// Write writes the dpos config settings to the database.
func (cfg *Config) Write(db fdb.Database, key []byte) error {
	data, err := cfg.EncodeRLP()
	if err != nil {
		return fmt.Errorf("Failed to rlp encode dpos config --- %v", err)
	}
	if err := db.Put(key, data); err != nil {
		return fmt.Errorf("Failed to store dpos config ---- %v", err)
	}
	return nil
}

// Read retrieves the consensus settings from the database.
func (cfg *Config) Read(db fdb.Database, key []byte) error {
	data, err := db.Get(key)
	if err != nil {
		return fmt.Errorf("Failed to load dpos config --- %v", err)
	}
	if err := cfg.DecodeRLP(data); err != nil {
		return fmt.Errorf("Failed to rlp decode dpos config --- %v", err)
	}
	return nil
}
