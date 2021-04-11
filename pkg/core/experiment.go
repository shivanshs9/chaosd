// Copyright 2020 Chaos Mesh Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package core

import (
	"context"
	"encoding/json"
	"time"

	perr "github.com/pkg/errors"
)

const (
	Created   = "created"
	Success   = "success"
	Error     = "error"
	Scheduled = "scheduled"
	Destroyed = "destroyed"
	Revoked   = "revoked"
)

const (
	ProcessAttack = "process"
	NetworkAttack = "network"
	StressAttack  = "stress"
	DiskAttack    = "disk"
	HostAttack    = "host"
)

// ExperimentStore defines operations for working with experiments
type ExperimentStore interface {
	List(ctx context.Context) ([]*Experiment, error)
	ListByConditions(ctx context.Context, conds *SearchCommand) ([]*Experiment, error)
	ListByStatus(ctx context.Context, status string) ([]*Experiment, error)
	FindByUid(ctx context.Context, uid string) (*Experiment, error)
	Set(ctx context.Context, exp *Experiment) error
	Update(ctx context.Context, uid, status, msg string, command string) error
}

// Experiment represents an experiment instance.
type Experiment struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Uid       string    `gorm:"index:uid" json:"uid"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
	Message   string    `json:"error"`
	// TODO: need to improve
	Kind           string `json:"kind"`
	Action         string `json:"action"`
	RecoverCommand string `json:"recover_command"`

	cachedRequestCommand AttackConfig
}

func (exp *Experiment) GetRequestCommand() (AttackConfig, error) {
	if exp.cachedRequestCommand != nil {
		return exp.cachedRequestCommand, nil
	}

	var attackConfig AttackConfig
	switch exp.Kind {
	case ProcessAttack:
		attackConfig = &ProcessCommand{}
	case NetworkAttack:
		attackConfig = &NetworkCommand{}
	case HostAttack:
		attackConfig = &HostCommand{}
	case StressAttack:
		attackConfig = &StressCommand{}
	case DiskAttack:
		attackConfig = &DiskCommand{}
	default:
		return nil, perr.Errorf("chaos experiment kind %s not found", exp.Kind)
	}

	if err := json.Unmarshal([]byte(exp.RecoverCommand), attackConfig); err != nil {
		return nil, err
	}
	exp.cachedRequestCommand = attackConfig
	return attackConfig, nil
}
