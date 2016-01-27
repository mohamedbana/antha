// /antharun/frontend.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package main

import (
	"fmt"

	"github.com/antha-lang/antha/internal/github.com/twinj/uuid"
	"github.com/antha-lang/antha/microArch/equipment"
	"github.com/antha-lang/antha/microArch/equipment/manual"
	"github.com/antha-lang/antha/microArch/equipment/void"
	"github.com/antha-lang/antha/microArch/logger"
	"github.com/antha-lang/antha/target"
)

const (
	DEBUG  = "debug"
	REMOTE = "remote"
	CUI    = "cui"
)

type Options struct {
	Kind   string
	Target *target.Target
	URI    string
}

type Frontend struct {
	shutdowns []func() error
}

func getEq(opts Options) (equipment.Equipment, error) {
	id := uuid.NewV4().String()
	switch opts.Kind {
	case CUI:
		eq := manual.NewAnthaManualCUI(id)
		logger.RegisterMiddleware(eq.Cui)
		return eq, nil
	case DEBUG:
		return void.NewVoidEquipment(id), nil
	case REMOTE:
		return manual.NewAnthaManualGrpc(id, opts.URI), nil
	default:
		return nil, fmt.Errorf("unknown frontend %q", opts.Kind)
	}
}

func NewFrontend(opts Options) (*Frontend, error) {
	t := opts.Target
	if eq, err := getEq(opts); err != nil {
		return nil, err
	} else if err := eq.Init(); err != nil {
		return nil, err
	} else {
		t.AddLiquidHandler(eq)
		return &Frontend{
			shutdowns: []func() error{func() error {
				return eq.Shutdown()
			}},
		}, nil
	}
}

func (a *Frontend) Shutdown() (err error) {
	for _, fn := range a.shutdowns {
		if e := fn(); e != nil {
			err = e
		}
	}
	return
}
