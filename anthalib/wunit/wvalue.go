// wunit/wvalue.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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
// 1 Royal College St, London NW1 0NH UK

package wunit

import (
	"sync"
)

/**********************************************************/

// a float with a RWMutex... may not be necessary
type wfloat struct{
	value float64
	mutex *sync.RWMutex
}

func (wv *wfloat)MakeMutex(){
	if wv.mutex==nil{
		m:=sync.RWMutex{}
		wv.mutex=&m
	}
}

// wrap a float in the wvalue structure
func NewWFloat(v float64) wfloat{
	wv:=wfloat{}
	wv.MakeMutex()
	wv.SetValue(v)
	return wv
}

// set float value
func (wv *wfloat)SetValue(v float64) float64{
	wv.mutex.Lock()
	defer wv.mutex.Unlock()
	r:=wv.value
	wv.value=v
	return r
}

// get float value
func (wv *wfloat)GetValue()float64{
	wv.mutex.RLock()
	defer wv.mutex.RUnlock()
	return wv.value
}

/********************************************/

// value structure wrapping a string 
type wstring struct{
	value string
	mutex *sync.RWMutex
}

func (wv *wstring)MakeMutex(){
	if wv.mutex==nil{
		m:=sync.RWMutex{}
		wv.mutex=&m
	}
}
func NewWString(v string) wstring{
	wv:=wstring{}
	wv.MakeMutex()
	wv.SetValue(v)
	return wv
}

// set the string value
func (wv *wstring)SetValue(s string)string{
	wv.mutex.Lock()
	defer wv.mutex.Unlock()
	r:=wv.value
	wv.value=s
	return r
}

// get the string value
func (wv *wstring)GetValue()string{
	wv.mutex.RLock()
	defer wv.mutex.RUnlock()
	return wv.value
}

/********************************************/

// an int plus a mutex
type wint struct{
	value int
	mutex *sync.RWMutex
}

func (wv *wint)MakeMutex(){
	if wv.mutex==nil{
		m:=sync.RWMutex{}
		wv.mutex=&m
	}
}

// wrap an int in the wvalue structure
func NewWInt(v int) wint{
	wv:=wint{}
	wv.MakeMutex()
	wv.SetValue(v)
	return wv
}

// set int value
func (wv *wint)SetValue(i int)int{
	wv.mutex.Lock()
	defer wv.mutex.Unlock()
	r:=wv.value
	wv.value=i
	return r
}

// get int value
func (wv *wint)GetValue()int{
	wv.mutex.RLock()
	defer wv.mutex.RUnlock()
	return wv.value
}
