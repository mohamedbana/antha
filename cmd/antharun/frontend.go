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
	"io"
	"os"

	lhclient "github.com/antha-lang/antha/driver/lh/pb/client"
	"github.com/antha-lang/antha/microArch/logger"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
	"github.com/antha-lang/antha/target/mixer"
)

const (
	DEBUG  = "debug"
	REMOTE = "remote"
)

type FrontendOpt struct {
	Kind     string
	Target   *target.Target
	MixerOpt mixer.Opt
	URI      string
}

type Frontend struct {
	shutdowns []func() error
}

func (a *Frontend) Shutdown() (err error) {
	for _, fn := range a.shutdowns {
		if e := fn(); e != nil {
			err = e
		}
	}
	return
}

type middleware struct {
	out io.Writer
}

func (a *middleware) Log(lvl logger.LogLevel, ts int64, source, msg string, extra ...interface{}) {
	if lvl == logger.TRACK {
		return
	}
	fmt.Fprint(a.out, msg)
	fmt.Fprint(a.out, extra...)
	fmt.Fprint(a.out, "\n")
}

func (a *middleware) Measure(ts int64, source, msg string, extra ...interface{}) {}

func (a *middleware) Sensor(ts int64, source, msg string, extra ...interface{}) {}

func NewFrontend(opt FrontendOpt) (*Frontend, error) {
	t := opt.Target

	switch opt.Kind {
	case DEBUG:
		t.AddDevice(human.New(human.Opt{CanMix: true}))
	case REMOTE:
		if m, err := mixer.New(opt.MixerOpt, lhclient.NewDriver(opt.URI)); err != nil {
			return nil, err
		} else {
			t.AddDevice(m)
			t.AddDevice(human.New(human.Opt{CanMix: false}))
		}
	default:
		return nil, fmt.Errorf("unknown frontend %q", opt.Kind)
	}

	mw := &middleware{out: os.Stdout}
	logger.RegisterMiddleware(mw)

	ret := &Frontend{}
	ret.shutdowns = append(ret.shutdowns, func() error {
		logger.UnregisterMiddleware(mw)
		return nil
	})

	return ret, nil
}
