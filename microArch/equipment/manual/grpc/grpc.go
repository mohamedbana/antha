// microArch/equipment/manual/grpc/grpc.go: Part of the Antha language
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

package grpc

import (
	"encoding/json"
	"log"

	"github.com/antha-lang/antha/antha/anthalib/material"
	i_1 "github.com/antha-lang/antha/antha/anthalib/wtype"
	i_2 "github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/bvendor/google.golang.org/grpc"
	i_4 "github.com/antha-lang/antha/microArch/driver"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	i_3 "github.com/antha-lang/antha/microArch/driver/liquidhandling"
	pb "github.com/antha-lang/antha/microArch/equipment/manual/grpc/ExtendedLiquidhandlingDriver"
)

type Driver struct {
	C pb.ExtendedLiquidhandlingDriverClient
	// ignore the below: it's just there to ensure we use all imports
	d i_3.ExtendedLiquidhandlingDriver
}

func NewDriver(address string) *Driver {
	var d Driver
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Cannot initialize driver")
	}
	d.C = pb.NewExtendedLiquidhandlingDriverClient(conn)
	return &d
}
func Encodeinterface(arg interface{}) *pb.AnyMessage {
	s, _ := json.Marshal(arg)
	ret := pb.AnyMessage{string(s)}
	return &ret
}
func Decodeinterface(msg *pb.AnyMessage) interface{} {
	var v interface{}
	json.Unmarshal([]byte(msg.Arg_1), &v)
	return v
}

//func main() {
//	d := NewDriver()
//	d.Go()
//}
func EncodePtrToLHProperties(arg *i_3.LHProperties) *pb.PtrToLHPropertiesMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHPropertiesMessage{
		EncodeLHProperties(*arg),
	}
	return &ret
}
func DecodePtrToLHProperties(arg *pb.PtrToLHPropertiesMessage) *i_3.LHProperties {
	ret := DecodeLHProperties(arg.Arg_1)
	return &ret
}
func EncodeArrayOfstring(arg []string) *pb.ArrayOfstring {
	a := make([]string, len(arg))
	for i, v := range arg {
		a[i] = (string)(v)
	}
	ret := pb.ArrayOfstring{
		a,
	}
	return &ret
}
func DecodeArrayOfstring(arg *pb.ArrayOfstring) []string {
	ret := make(([]string), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = (string)(v)
	}
	return ret
}
func EncodeArrayOffloat64(arg []float64) *pb.ArrayOfdouble {
	a := make([]float64, len(arg))
	for i, v := range arg {
		a[i] = (float64)(v)
	}
	ret := pb.ArrayOfdouble{
		a,
	}
	return &ret
}
func DecodeArrayOffloat64(arg *pb.ArrayOfdouble) []float64 {
	ret := make(([]float64), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = (float64)(v)
	}
	return ret
}
func EncodeArrayOfbool(arg []bool) *pb.ArrayOfbool {
	a := make([]bool, len(arg))
	for i, v := range arg {
		a[i] = (bool)(v)
	}
	ret := pb.ArrayOfbool{
		a,
	}
	return &ret
}
func DecodeArrayOfbool(arg *pb.ArrayOfbool) []bool {
	ret := make(([]bool), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = (bool)(v)
	}
	return ret
}
func EncodeMapstringinterfaceMessage(arg map[string]interface{}) *pb.MapstringAnyMessageMessage {
	a := make([]*pb.MapstringAnyMessageMessageFieldEntry, 0, len(arg))
	for k, v := range arg {
		fe := EncodeMapstringinterfaceMessageFieldEntry(k, v)
		a = append(a, &fe)
	}
	ret := pb.MapstringAnyMessageMessage{
		a,
	}
	return &ret
}
func EncodeMapstringinterfaceMessageFieldEntry(k string, v interface{}) pb.MapstringAnyMessageMessageFieldEntry {
	ret := pb.MapstringAnyMessageMessageFieldEntry{
		(string)(k),
		Encodeinterface(v),
	}
	return ret
}
func DecodeMapstringinterfaceMessage(arg *pb.MapstringAnyMessageMessage) map[string]interface{} {
	a := make(map[(string)](interface{}), len(arg.MapField))
	for _, fe := range arg.MapField {
		k, v := DecodeMapstringinterfaceMessageFieldEntry(fe)
		a[k] = v
	}
	return a
}
func DecodeMapstringinterfaceMessageFieldEntry(arg *pb.MapstringAnyMessageMessageFieldEntry) (string, interface{}) {
	k := (string)(arg.Key)
	v := Decodeinterface(arg.Value)
	return k, v
}
func EncodeCommandStatus(arg i_4.CommandStatus) *pb.CommandStatusMessage {
	ret := pb.CommandStatusMessage{(bool)(arg.OK), int64(arg.Errorcode), (string)(arg.Msg)}
	return &ret
}
func DecodeCommandStatus(arg *pb.CommandStatusMessage) i_4.CommandStatus {
	ret := i_4.CommandStatus{(bool)(arg.Arg_1), (int)(arg.Arg_2), (string)(arg.Arg_3)}
	return ret
}
func EncodeArrayOfint(arg []int) *pb.ArrayOfint64 {
	a := make([]int64, len(arg))
	for i, v := range arg {
		a[i] = int64(v)
	}
	ret := pb.ArrayOfint64{
		a,
	}
	return &ret
}
func DecodeArrayOfint(arg *pb.ArrayOfint64) []int {
	ret := make(([]int), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = (int)(v)
	}
	return ret
}
func EncodeLHProperties(arg i_3.LHProperties) *pb.LHPropertiesMessage {
	ret := pb.LHPropertiesMessage{
		(string)(arg.ID),
		int64(arg.Nposns),
		EncodeMapstringPtrToLHPositionMessage(arg.Positions),
		EncodeMapstringinterfaceMessage(arg.PlateLookup),
		EncodeMapstringstringMessage(arg.PosLookup),
		EncodeMapstringstringMessage(arg.PlateIDLookup),
		EncodeMapstringPtrToLHPlateMessage(arg.Plates),
		EncodeMapstringPtrToLHTipboxMessage(arg.Tipboxes),
		EncodeMapstringPtrToLHTipwasteMessage(arg.Tipwastes),
		EncodeMapstringPtrToLHPlateMessage(arg.Wastes),
		EncodeMapstringPtrToLHPlateMessage(arg.Washes),
		EncodeMapstringstringMessage(arg.Devices),
		(string)(arg.Model),
		(string)(arg.Mnfr),
		(string)(arg.LHType),
		(string)(arg.TipType),
		EncodeArrayOfPtrToLHHead(arg.Heads),
		EncodeArrayOfPtrToLHHead(arg.HeadsLoaded),
		EncodeArrayOfPtrToLHAdaptor(arg.Adaptors),
		EncodeArrayOfPtrToLHTip(arg.Tips),
		EncodeArrayOfstring(arg.Tip_preferences),
		EncodeArrayOfstring(arg.Input_preferences),
		EncodeArrayOfstring(arg.Output_preferences),
		EncodeArrayOfstring(arg.Tipwaste_preferences),
		EncodeArrayOfstring(arg.Waste_preferences),
		EncodeArrayOfstring(arg.Wash_preferences),
		EncodePtrToLHChannelParameter(arg.CurrConf),
		EncodeArrayOfPtrToLHChannelParameter(arg.Cnfvol),
		EncodeMapstringCoordinatesMessage(arg.Layout),
		int64(arg.MaterialType),
	}
	return &ret
}
func DecodeLHProperties(arg *pb.LHPropertiesMessage) i_3.LHProperties {
	ret := i_3.LHProperties{
		(string)(arg.Arg_1),
		(int)(arg.Arg_2),
		(map[string]*i_1.LHPosition)(DecodeMapstringPtrToLHPositionMessage(arg.Arg_3)),
		(map[string]interface{})(DecodeMapstringinterfaceMessage(arg.Arg_4)),
		(map[string]string)(DecodeMapstringstringMessage(arg.Arg_5)),
		(map[string]string)(DecodeMapstringstringMessage(arg.Arg_6)),
		(map[string]*i_1.LHPlate)(DecodeMapstringPtrToLHPlateMessage(arg.Arg_7)),
		(map[string]*i_1.LHTipbox)(DecodeMapstringPtrToLHTipboxMessage(arg.Arg_8)),
		(map[string]*i_1.LHTipwaste)(DecodeMapstringPtrToLHTipwasteMessage(arg.Arg_9)),
		(map[string]*i_1.LHPlate)(DecodeMapstringPtrToLHPlateMessage(arg.Arg_10)),
		(map[string]*i_1.LHPlate)(DecodeMapstringPtrToLHPlateMessage(arg.Arg_11)),
		(map[string]string)(DecodeMapstringstringMessage(arg.Arg_12)),
		(string)(arg.Arg_13),
		(string)(arg.Arg_14),
		(string)(arg.Arg_15),
		(string)(arg.Arg_16),
		([]*i_1.LHHead)(DecodeArrayOfPtrToLHHead(arg.Arg_17)),
		([]*i_1.LHHead)(DecodeArrayOfPtrToLHHead(arg.Arg_18)),
		([]*i_1.LHAdaptor)(DecodeArrayOfPtrToLHAdaptor(arg.Arg_19)),
		([]*i_1.LHTip)(DecodeArrayOfPtrToLHTip(arg.Arg_20)),
		([]string)(DecodeArrayOfstring(arg.Arg_21)),
		([]string)(DecodeArrayOfstring(arg.Arg_22)),
		([]string)(DecodeArrayOfstring(arg.Arg_23)),
		([]string)(DecodeArrayOfstring(arg.Arg_24)),
		([]string)(DecodeArrayOfstring(arg.Arg_25)),
		([]string)(DecodeArrayOfstring(arg.Arg_26)),
		nil, //@jmanart Manually changed as the case is not being handled in the pbtogo code properly
		(*i_1.LHChannelParameter)(DecodePtrToLHChannelParameter(arg.Arg_27)),
		([]*i_1.LHChannelParameter)(DecodeArrayOfPtrToLHChannelParameter(arg.Arg_28)),
		(map[string]i_1.Coordinates)(DecodeMapstringCoordinatesMessage(arg.Arg_29)),
		(material.MaterialType)(arg.Arg_30), //@jmanart manual fix, material lib was not found
	}
	return ret
}
func EncodeMapstringPtrToLHPositionMessage(arg map[string]*i_1.LHPosition) *pb.MapstringPtrToLHPositionMessageMessage {
	a := make([]*pb.MapstringPtrToLHPositionMessageMessageFieldEntry, 0, len(arg))
	for k, v := range arg {
		fe := EncodeMapstringPtrToLHPositionMessageFieldEntry(k, v)
		a = append(a, &fe)
	}
	ret := pb.MapstringPtrToLHPositionMessageMessage{
		a,
	}
	return &ret
}
func EncodeMapstringPtrToLHPositionMessageFieldEntry(k string, v *i_1.LHPosition) pb.MapstringPtrToLHPositionMessageMessageFieldEntry {
	ret := pb.MapstringPtrToLHPositionMessageMessageFieldEntry{
		(string)(k),
		EncodePtrToLHPosition(v),
	}
	return ret
}
func DecodeMapstringPtrToLHPositionMessage(arg *pb.MapstringPtrToLHPositionMessageMessage) map[string]*i_1.LHPosition {
	a := make(map[(string)](*i_1.LHPosition), len(arg.MapField))
	for _, fe := range arg.MapField {
		k, v := DecodeMapstringPtrToLHPositionMessageFieldEntry(fe)
		a[k] = v
	}
	return a
}
func DecodeMapstringPtrToLHPositionMessageFieldEntry(arg *pb.MapstringPtrToLHPositionMessageMessageFieldEntry) (string, *i_1.LHPosition) {
	k := (string)(arg.Key)
	v := DecodePtrToLHPosition(arg.Value)
	return k, v
}
func EncodeMapstringCoordinatesMessage(arg map[string]i_1.Coordinates) *pb.MapstringCoordinatesMessageMessage {
	a := make([]*pb.MapstringCoordinatesMessageMessageFieldEntry, 0, len(arg))
	for k, v := range arg {
		fe := EncodeMapstringCoordinatesMessageFieldEntry(k, v)
		a = append(a, &fe)
	}
	ret := pb.MapstringCoordinatesMessageMessage{
		a,
	}
	return &ret
}
func EncodeMapstringCoordinatesMessageFieldEntry(k string, v i_1.Coordinates) pb.MapstringCoordinatesMessageMessageFieldEntry {
	ret := pb.MapstringCoordinatesMessageMessageFieldEntry{
		(string)(k),
		EncodeCoordinates(v),
	}
	return ret
}
func DecodeMapstringCoordinatesMessage(arg *pb.MapstringCoordinatesMessageMessage) map[string]i_1.Coordinates {
	a := make(map[(string)](i_1.Coordinates), len(arg.MapField))
	for _, fe := range arg.MapField {
		k, v := DecodeMapstringCoordinatesMessageFieldEntry(fe)
		a[k] = v
	}
	return a
}
func DecodeMapstringCoordinatesMessageFieldEntry(arg *pb.MapstringCoordinatesMessageMessageFieldEntry) (string, i_1.Coordinates) {
	k := (string)(arg.Key)
	v := DecodeCoordinates(arg.Value)
	return k, v
}
func EncodePtrToLHPlate(arg *i_1.LHPlate) *pb.PtrToLHPlateMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHPlateMessage{
		EncodeLHPlate(*arg),
	}
	return &ret
}
func DecodePtrToLHPlate(arg *pb.PtrToLHPlateMessage) *i_1.LHPlate {
	ret := DecodeLHPlate(arg.Arg_1)
	return &ret
}
func EncodeLHTipbox(arg i_1.LHTipbox) *pb.LHTipboxMessage {
	ret := pb.LHTipboxMessage{(string)(arg.ID), (string)(arg.Boxname), (string)(arg.Type), (string)(arg.Mnfr), int64(arg.Nrows), int64(arg.Ncols), (float64)(arg.Height), EncodePtrToLHTip(arg.Tiptype), EncodePtrToLHWell(arg.AsWell), int64(arg.NTips), EncodeArrayOfArrayOfPtrToLHTip(arg.Tips), (float64)(arg.TipXOffset), (float64)(arg.TipYOffset), (float64)(arg.TipXStart), (float64)(arg.TipYStart), (float64)(arg.TipZStart)}
	return &ret
}
func DecodeLHTipbox(arg *pb.LHTipboxMessage) i_1.LHTipbox {
	ret := i_1.LHTipbox{(string)(arg.Arg_2), (string)(arg.Arg_3), (string)(arg.Arg_4), (string)(arg.Arg_5), (int)(arg.Arg_6), (int)(arg.Arg_7), (float64)(arg.Arg_8), (*i_1.LHTip)(DecodePtrToLHTip(arg.Arg_9)), (*i_1.LHWell)(DecodePtrToLHWell(arg.Arg_10)), (int)(arg.Arg_11), ([][]*i_1.LHTip)(DecodeArrayOfArrayOfPtrToLHTip(arg.Arg_12)), (float64)(arg.Arg_13), (float64)(arg.Arg_14), (float64)(arg.Arg_15), (float64)(arg.Arg_16), (float64)(arg.Arg_17)}
	return ret
}
func EncodeLHTipwaste(arg i_1.LHTipwaste) *pb.LHTipwasteMessage {
	ret := pb.LHTipwasteMessage{(string)(arg.ID), (string)(arg.Type), (string)(arg.Mnfr), int64(arg.Capacity), int64(arg.Contents), (float64)(arg.Height), (float64)(arg.WellXStart), (float64)(arg.WellYStart), (float64)(arg.WellZStart), EncodePtrToLHWell(arg.AsWell)}
	return &ret
}
func DecodeLHTipwaste(arg *pb.LHTipwasteMessage) i_1.LHTipwaste {
	ret := i_1.LHTipwaste{(string)(arg.Arg_1), (string)(arg.Arg_2), (string)(arg.Arg_3), (int)(arg.Arg_4), (int)(arg.Arg_5), (float64)(arg.Arg_6), (float64)(arg.Arg_7), (float64)(arg.Arg_8), (float64)(arg.Arg_9), (*i_1.LHWell)(DecodePtrToLHWell(arg.Arg_10))}
	return ret
}
func EncodePtrToLHHead(arg *i_1.LHHead) *pb.PtrToLHHeadMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHHeadMessage{
		EncodeLHHead(*arg),
	}
	return &ret
}
func DecodePtrToLHHead(arg *pb.PtrToLHHeadMessage) *i_1.LHHead {
	ret := DecodeLHHead(arg.Arg_1)
	return &ret
}
func EncodeLHChannelParameter(arg i_1.LHChannelParameter) *pb.LHChannelParameterMessage {
	ret := pb.LHChannelParameterMessage{(string)(arg.ID), (string)(arg.Name), EncodePtrToVolume(arg.Minvol), EncodePtrToVolume(arg.Maxvol), EncodePtrToFlowRate(arg.Minspd), EncodePtrToFlowRate(arg.Maxspd), int64(arg.Multi), (bool)(arg.Independent), int64(arg.Orientation), int64(arg.Head)}
	return &ret
}
func DecodeLHChannelParameter(arg *pb.LHChannelParameterMessage) i_1.LHChannelParameter {
	// this could be nil
	if arg == nil {
		// return an empty thing
		var v i_1.LHChannelParameter
		return v
	} else {
		ret := i_1.LHChannelParameter{(string)(arg.Arg_1), (string)(arg.Arg_2), (*i_2.Volume)(DecodePtrToVolume(arg.Arg_3)), (*i_2.Volume)(DecodePtrToVolume(arg.Arg_4)), (*i_2.FlowRate)(DecodePtrToFlowRate(arg.Arg_5)), (*i_2.FlowRate)(DecodePtrToFlowRate(arg.Arg_6)), (int)(arg.Arg_7), (bool)(arg.Arg_8), (int)(arg.Arg_9), (int)(arg.Arg_10)}
		return ret
	}
}
func EncodeMapstringstringMessage(arg map[string]string) *pb.MapstringstringMessage {
	a := make([]*pb.MapstringstringMessageFieldEntry, 0, len(arg))
	for k, v := range arg {
		fe := EncodeMapstringstringMessageFieldEntry(k, v)
		a = append(a, &fe)
	}
	ret := pb.MapstringstringMessage{
		a,
	}
	return &ret
}
func EncodeMapstringstringMessageFieldEntry(k string, v string) pb.MapstringstringMessageFieldEntry {
	ret := pb.MapstringstringMessageFieldEntry{
		(string)(k),
		(string)(v),
	}
	return ret
}
func DecodeMapstringstringMessage(arg *pb.MapstringstringMessage) map[string]string {
	a := make(map[(string)](string), len(arg.MapField))
	for _, fe := range arg.MapField {
		k, v := DecodeMapstringstringMessageFieldEntry(fe)
		a[k] = v
	}
	return a
}
func DecodeMapstringstringMessageFieldEntry(arg *pb.MapstringstringMessageFieldEntry) (string, string) {
	k := (string)(arg.Key)
	v := (string)(arg.Value)
	return k, v
}
func EncodePtrToLHTipbox(arg *i_1.LHTipbox) *pb.PtrToLHTipboxMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHTipboxMessage{
		EncodeLHTipbox(*arg),
	}
	return &ret
}
func DecodePtrToLHTipbox(arg *pb.PtrToLHTipboxMessage) *i_1.LHTipbox {
	ret := DecodeLHTipbox(arg.Arg_1)
	return &ret
}
func EncodeLHTip(arg i_1.LHTip) *pb.LHTipMessage {
	ret := pb.LHTipMessage{(string)(arg.ID), (string)(arg.Type), (string)(arg.Mnfr), (bool)(arg.Dirty), EncodePtrToVolume(arg.MaxVol), EncodePtrToVolume(arg.MinVol)}
	return &ret
}
func DecodeLHTip(arg *pb.LHTipMessage) i_1.LHTip {
	if arg == nil {
		return i_1.LHTip{}
	} else {
		ret := i_1.LHTip{(string)(arg.Arg_1), (string)(arg.Arg_2), (string)(arg.Arg_3), (bool)(arg.Arg_4), (*i_2.Volume)(DecodePtrToVolume(arg.Arg_5)), (*i_2.Volume)(DecodePtrToVolume(arg.Arg_6))}
		return ret
	}
}
func EncodeArrayOfPtrToLHChannelParameter(arg []*i_1.LHChannelParameter) *pb.ArrayOfPtrToLHChannelParameterMessage {
	a := make([]*pb.PtrToLHChannelParameterMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodePtrToLHChannelParameter(v)
	}
	ret := pb.ArrayOfPtrToLHChannelParameterMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfPtrToLHChannelParameter(arg *pb.ArrayOfPtrToLHChannelParameterMessage) []*i_1.LHChannelParameter {
	ret := make(([]*i_1.LHChannelParameter), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodePtrToLHChannelParameter(v)
	}
	return ret
}
func EncodeLHPlate(arg i_1.LHPlate) *pb.LHPlateMessage {
	ret := pb.LHPlateMessage{
		(string)(arg.ID),
		(string)(arg.Inst),
		(string)(arg.Loc),
		(string)(arg.PlateName),
		(string)(arg.Type),
		(string)(arg.Mnfr),
		int64(arg.WlsX),
		int64(arg.WlsY),
		int64(arg.Nwells),
		EncodeMapstringPtrToLHWellMessage(arg.HWells),
		(float64)(arg.Height),
		(string)(arg.Hunit),
		EncodeArrayOfArrayOfPtrToLHWell(arg.Rows),
		EncodeArrayOfArrayOfPtrToLHWell(arg.Cols),
		EncodePtrToLHWell(arg.Welltype),
		EncodeMapstringPtrToLHWellMessage(arg.Wellcoords),
		(float64)(arg.WellXOffset),
		(float64)(arg.WellYOffset),
		(float64)(arg.WellXStart),
		(float64)(arg.WellYStart),
		(float64)(arg.WellZStart),
	}
	return &ret
}
func DecodeLHPlate(arg *pb.LHPlateMessage) i_1.LHPlate {
	ret := i_1.LHPlate{(string)(arg.Arg_2), (string)(arg.Arg_3), (string)(arg.Arg_4), (string)(arg.Arg_5), (string)(arg.Arg_6), (string)(arg.Arg_7), (int)(arg.Arg_8), (int)(arg.Arg_9), (int)(arg.Arg_10), (map[string]*i_1.LHWell)(DecodeMapstringPtrToLHWellMessage(arg.Arg_11)), (float64)(arg.Arg_12), (string)(arg.Arg_13), ([][]*i_1.LHWell)(DecodeArrayOfArrayOfPtrToLHWell(arg.Arg_14)), ([][]*i_1.LHWell)(DecodeArrayOfArrayOfPtrToLHWell(arg.Arg_15)), (*i_1.LHWell)(DecodePtrToLHWell(arg.Arg_16)), (map[string]*i_1.LHWell)(DecodeMapstringPtrToLHWellMessage(arg.Arg_17)), (float64)(arg.Arg_18), (float64)(arg.Arg_19), (float64)(arg.Arg_20), (float64)(arg.Arg_21), (float64)(arg.Arg_22)}
	return ret
}
func EncodePtrToLHTipwaste(arg *i_1.LHTipwaste) *pb.PtrToLHTipwasteMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHTipwasteMessage{
		EncodeLHTipwaste(*arg),
	}
	return &ret
}
func DecodePtrToLHTipwaste(arg *pb.PtrToLHTipwasteMessage) *i_1.LHTipwaste {
	ret := DecodeLHTipwaste(arg.Arg_1)
	return &ret
}
func EncodeMapstringPtrToLHTipwasteMessage(arg map[string]*i_1.LHTipwaste) *pb.MapstringPtrToLHTipwasteMessageMessage {
	a := make([]*pb.MapstringPtrToLHTipwasteMessageMessageFieldEntry, 0, len(arg))
	for k, v := range arg {
		fe := EncodeMapstringPtrToLHTipwasteMessageFieldEntry(k, v)
		a = append(a, &fe)
	}
	ret := pb.MapstringPtrToLHTipwasteMessageMessage{
		a,
	}
	return &ret
}
func EncodeMapstringPtrToLHTipwasteMessageFieldEntry(k string, v *i_1.LHTipwaste) pb.MapstringPtrToLHTipwasteMessageMessageFieldEntry {
	ret := pb.MapstringPtrToLHTipwasteMessageMessageFieldEntry{
		(string)(k),
		EncodePtrToLHTipwaste(v),
	}
	return ret
}
func DecodeMapstringPtrToLHTipwasteMessage(arg *pb.MapstringPtrToLHTipwasteMessageMessage) map[string]*i_1.LHTipwaste {
	a := make(map[(string)](*i_1.LHTipwaste), len(arg.MapField))
	for _, fe := range arg.MapField {
		k, v := DecodeMapstringPtrToLHTipwasteMessageFieldEntry(fe)
		a[k] = v
	}
	return a
}
func DecodeMapstringPtrToLHTipwasteMessageFieldEntry(arg *pb.MapstringPtrToLHTipwasteMessageMessageFieldEntry) (string, *i_1.LHTipwaste) {
	k := (string)(arg.Key)
	v := DecodePtrToLHTipwaste(arg.Value)
	return k, v
}
func EncodeLHAdaptor(arg i_1.LHAdaptor) *pb.LHAdaptorMessage {
	ret := pb.LHAdaptorMessage{(string)(arg.Name), (string)(arg.ID), (string)(arg.Manufacturer), EncodePtrToLHChannelParameter(arg.Params), int64(arg.Ntipsloaded), EncodePtrToLHTip(arg.Tiptypeloaded)}
	return &ret
}
func DecodeLHAdaptor(arg *pb.LHAdaptorMessage) i_1.LHAdaptor {
	ret := i_1.LHAdaptor{(string)(arg.Arg_1), (string)(arg.Arg_2), (string)(arg.Arg_3), (*i_1.LHChannelParameter)(DecodePtrToLHChannelParameter(arg.Arg_4)), (int)(arg.Arg_5), (*i_1.LHTip)(DecodePtrToLHTip(arg.Arg_6))}
	return ret
}
func EncodeArrayOfPtrToLHTip(arg []*i_1.LHTip) *pb.ArrayOfPtrToLHTipMessage {
	a := make([]*pb.PtrToLHTipMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodePtrToLHTip(v)
	}
	ret := pb.ArrayOfPtrToLHTipMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfPtrToLHTip(arg *pb.ArrayOfPtrToLHTipMessage) []*i_1.LHTip {
	ret := make(([]*i_1.LHTip), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodePtrToLHTip(v)
	}
	return ret
}
func EncodeCoordinates(arg i_1.Coordinates) *pb.CoordinatesMessage {
	ret := pb.CoordinatesMessage{(float64)(arg.X), (float64)(arg.Y), (float64)(arg.Z)}
	return &ret
}
func DecodeCoordinates(arg *pb.CoordinatesMessage) i_1.Coordinates {
	ret := i_1.Coordinates{(float64)(arg.Arg_1), (float64)(arg.Arg_2), (float64)(arg.Arg_3)}
	return ret
}
func EncodeLHPosition(arg i_1.LHPosition) *pb.LHPositionMessage {
	ret := pb.LHPositionMessage{(string)(arg.ID), (string)(arg.Name), int64(arg.Num), EncodeArrayOfLHDevice(arg.Extra), (float64)(arg.Maxh)}
	return &ret
}
func DecodeLHPosition(arg *pb.LHPositionMessage) i_1.LHPosition {
	ret := i_1.LHPosition{(string)(arg.Arg_1), (string)(arg.Arg_2), (int)(arg.Arg_3), ([]i_1.LHDevice)(DecodeArrayOfLHDevice(arg.Arg_4)), (float64)(arg.Arg_5)}
	return ret
}
func EncodeLHHead(arg i_1.LHHead) *pb.LHHeadMessage {
	ret := pb.LHHeadMessage{(string)(arg.Name), (string)(arg.Manufacturer), (string)(arg.ID), EncodePtrToLHAdaptor(arg.Adaptor), EncodePtrToLHChannelParameter(arg.Params)}
	return &ret
}
func DecodeLHHead(arg *pb.LHHeadMessage) i_1.LHHead {
	ret := i_1.LHHead{(string)(arg.Arg_1), (string)(arg.Arg_2), (string)(arg.Arg_3), (*i_1.LHAdaptor)(DecodePtrToLHAdaptor(arg.Arg_4)), (*i_1.LHChannelParameter)(DecodePtrToLHChannelParameter(arg.Arg_5))}
	return ret
}
func EncodePtrToLHAdaptor(arg *i_1.LHAdaptor) *pb.PtrToLHAdaptorMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHAdaptorMessage{
		EncodeLHAdaptor(*arg),
	}
	return &ret
}
func DecodePtrToLHAdaptor(arg *pb.PtrToLHAdaptorMessage) *i_1.LHAdaptor {
	ret := DecodeLHAdaptor(arg.Arg_1)
	return &ret
}
func EncodePtrToLHChannelParameter(arg *i_1.LHChannelParameter) *pb.PtrToLHChannelParameterMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHChannelParameterMessage{
		EncodeLHChannelParameter(*arg),
	}
	return &ret
}
func DecodePtrToLHChannelParameter(arg *pb.PtrToLHChannelParameterMessage) *i_1.LHChannelParameter {
	ret := DecodeLHChannelParameter(arg.Arg_1)
	// pointers-to-empty-things are actually considered nil
	// this needs addressing correctly in the future
	if ret.ID == "" {
		// this only happens if the structure has not been properly initialized
		return nil
	} else {
		return &ret
	}
}
func EncodePtrToLHPosition(arg *i_1.LHPosition) *pb.PtrToLHPositionMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHPositionMessage{
		EncodeLHPosition(*arg),
	}
	return &ret
}
func DecodePtrToLHPosition(arg *pb.PtrToLHPositionMessage) *i_1.LHPosition {
	ret := DecodeLHPosition(arg.Arg_1)
	return &ret
}
func EncodeMapstringPtrToLHPlateMessage(arg map[string]*i_1.LHPlate) *pb.MapstringPtrToLHPlateMessageMessage {
	a := make([]*pb.MapstringPtrToLHPlateMessageMessageFieldEntry, 0, len(arg))
	for k, v := range arg {
		fe := EncodeMapstringPtrToLHPlateMessageFieldEntry(k, v)
		a = append(a, &fe)
	}
	ret := pb.MapstringPtrToLHPlateMessageMessage{
		a,
	}
	return &ret
}
func EncodeMapstringPtrToLHPlateMessageFieldEntry(k string, v *i_1.LHPlate) pb.MapstringPtrToLHPlateMessageMessageFieldEntry {
	ret := pb.MapstringPtrToLHPlateMessageMessageFieldEntry{
		(string)(k),
		EncodePtrToLHPlate(v),
	}
	return ret
}
func DecodeMapstringPtrToLHPlateMessage(arg *pb.MapstringPtrToLHPlateMessageMessage) map[string]*i_1.LHPlate {
	a := make(map[(string)](*i_1.LHPlate), len(arg.MapField))
	for _, fe := range arg.MapField {
		k, v := DecodeMapstringPtrToLHPlateMessageFieldEntry(fe)
		a[k] = v
	}
	return a
}
func DecodeMapstringPtrToLHPlateMessageFieldEntry(arg *pb.MapstringPtrToLHPlateMessageMessageFieldEntry) (string, *i_1.LHPlate) {
	k := (string)(arg.Key)
	v := DecodePtrToLHPlate(arg.Value)
	return k, v
}
func EncodeMapstringPtrToLHTipboxMessage(arg map[string]*i_1.LHTipbox) *pb.MapstringPtrToLHTipboxMessageMessage {
	a := make([]*pb.MapstringPtrToLHTipboxMessageMessageFieldEntry, 0, len(arg))
	for k, v := range arg {
		fe := EncodeMapstringPtrToLHTipboxMessageFieldEntry(k, v)
		a = append(a, &fe)
	}
	ret := pb.MapstringPtrToLHTipboxMessageMessage{
		a,
	}
	return &ret
}
func EncodeMapstringPtrToLHTipboxMessageFieldEntry(k string, v *i_1.LHTipbox) pb.MapstringPtrToLHTipboxMessageMessageFieldEntry {
	ret := pb.MapstringPtrToLHTipboxMessageMessageFieldEntry{
		(string)(k),
		EncodePtrToLHTipbox(v),
	}
	return ret
}
func DecodeMapstringPtrToLHTipboxMessage(arg *pb.MapstringPtrToLHTipboxMessageMessage) map[string]*i_1.LHTipbox {
	a := make(map[(string)](*i_1.LHTipbox), len(arg.MapField))
	for _, fe := range arg.MapField {
		k, v := DecodeMapstringPtrToLHTipboxMessageFieldEntry(fe)
		a[k] = v
	}
	return a
}
func DecodeMapstringPtrToLHTipboxMessageFieldEntry(arg *pb.MapstringPtrToLHTipboxMessageMessageFieldEntry) (string, *i_1.LHTipbox) {
	k := (string)(arg.Key)
	v := DecodePtrToLHTipbox(arg.Value)
	return k, v
}
func EncodeArrayOfPtrToLHHead(arg []*i_1.LHHead) *pb.ArrayOfPtrToLHHeadMessage {
	a := make([]*pb.PtrToLHHeadMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodePtrToLHHead(v)
	}
	ret := pb.ArrayOfPtrToLHHeadMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfPtrToLHHead(arg *pb.ArrayOfPtrToLHHeadMessage) []*i_1.LHHead {
	ret := make(([]*i_1.LHHead), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodePtrToLHHead(v)
	}
	return ret
}
func EncodeArrayOfPtrToLHAdaptor(arg []*i_1.LHAdaptor) *pb.ArrayOfPtrToLHAdaptorMessage {
	a := make([]*pb.PtrToLHAdaptorMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodePtrToLHAdaptor(v)
	}
	ret := pb.ArrayOfPtrToLHAdaptorMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfPtrToLHAdaptor(arg *pb.ArrayOfPtrToLHAdaptorMessage) []*i_1.LHAdaptor {
	ret := make(([]*i_1.LHAdaptor), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodePtrToLHAdaptor(v)
	}
	return ret
}
func EncodePtrToLHTip(arg *i_1.LHTip) *pb.PtrToLHTipMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHTipMessage{
		EncodeLHTip(*arg),
	}
	return &ret
}
func DecodePtrToLHTip(arg *pb.PtrToLHTipMessage) *i_1.LHTip {
	ret := DecodeLHTip(arg.Arg_1)
	return &ret
}
func EncodeArrayOfArrayOfPtrToLHTip(arg [][]*i_1.LHTip) *pb.ArrayOfArrayOfPtrToLHTipMessage {
	a := make([]*pb.ArrayOfPtrToLHTipMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodeArrayOfPtrToLHTip(v)
	}
	ret := pb.ArrayOfArrayOfPtrToLHTipMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfArrayOfPtrToLHTip(arg *pb.ArrayOfArrayOfPtrToLHTipMessage) [][]*i_1.LHTip {
	ret := make(([][]*i_1.LHTip), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodeArrayOfPtrToLHTip(v)
	}
	return ret
}
func EncodeArrayOfLHDevice(arg []i_1.LHDevice) *pb.ArrayOfLHDeviceMessage {
	a := make([]*pb.LHDeviceMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodeLHDevice(v)
	}
	ret := pb.ArrayOfLHDeviceMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfLHDevice(arg *pb.ArrayOfLHDeviceMessage) []i_1.LHDevice {
	ret := make(([]i_1.LHDevice), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodeLHDevice(v)
	}
	return ret
}
func EncodeVolume(arg i_2.Volume) *pb.VolumeMessage {
	ret := pb.VolumeMessage{EncodeConcreteMeasurement(arg.ConcreteMeasurement)}
	return &ret
}
func DecodeVolume(arg *pb.VolumeMessage) i_2.Volume {
	ret := i_2.Volume{(i_2.ConcreteMeasurement)(DecodeConcreteMeasurement(arg.Arg_1))}
	return ret
}
func EncodePtrToLHWell(arg *i_1.LHWell) *pb.PtrToLHWellMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHWellMessage{
		EncodeLHWell(*arg),
	}
	return &ret
}
func DecodePtrToLHWell(arg *pb.PtrToLHWellMessage) *i_1.LHWell {
	ret := DecodeLHWell(arg.Arg_1)
	return &ret
}
func EncodePtrToVolume(arg *i_2.Volume) *pb.PtrToVolumeMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToVolumeMessage{
		EncodeVolume(*arg),
	}
	return &ret
}
func DecodePtrToVolume(arg *pb.PtrToVolumeMessage) *i_2.Volume {
	ret := DecodeVolume(arg.Arg_1)
	return &ret
}
func EncodeLHDevice(arg i_1.LHDevice) *pb.LHDeviceMessage {
	ret := pb.LHDeviceMessage{(string)(arg.ID), (string)(arg.Name), (string)(arg.Mnfr)}
	return &ret
}
func DecodeLHDevice(arg *pb.LHDeviceMessage) i_1.LHDevice {
	ret := i_1.LHDevice{(string)(arg.Arg_1), (string)(arg.Arg_2), (string)(arg.Arg_3)}
	return ret
}
func EncodePtrToFlowRate(arg *i_2.FlowRate) *pb.PtrToFlowRateMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToFlowRateMessage{
		EncodeFlowRate(*arg),
	}
	return &ret
}
func DecodePtrToFlowRate(arg *pb.PtrToFlowRateMessage) *i_2.FlowRate {
	ret := DecodeFlowRate(arg.Arg_1)
	return &ret
}
func EncodeFlowRate(arg i_2.FlowRate) *pb.FlowRateMessage {
	ret := pb.FlowRateMessage{EncodeConcreteMeasurement(arg.ConcreteMeasurement)}
	return &ret
}
func DecodeFlowRate(arg *pb.FlowRateMessage) i_2.FlowRate {
	ret := i_2.FlowRate{(i_2.ConcreteMeasurement)(DecodeConcreteMeasurement(arg.Arg_1))}
	return ret
}
func EncodeLHWell(arg i_1.LHWell) *pb.LHWellMessage {
	ret := pb.LHWellMessage{
		(string)(arg.ID),
		(string)(arg.Inst),
		(string)(arg.Plateinst),
		(string)(arg.Plateid),
		(string)(arg.Platetype),
		(string)(arg.Crds),
		(float64)(arg.Vol),
		(string)(arg.Vunit),
		EncodeArrayOfPtrToLHComponent(arg.WContents),
		(float64)(arg.Rvol), (float64)(arg.Currvol),
		EncodePtrToShape(arg.WShape),
		int64(arg.Bottom),
		(float64)(arg.Xdim),
		(float64)(arg.Ydim),
		(float64)(arg.Zdim),
		(float64)(arg.Bottomh),
		(string)(arg.Dunit),
		EncodeMapstringinterfaceMessage(arg.Extra),
		//EncodePtrToLHPlate(arg.Plate)//@jmanart gotopb cycle
	}
	return &ret
}
func DecodeLHWell(arg *pb.LHWellMessage) i_1.LHWell {
	ret := i_1.LHWell{
		(string)(arg.Arg_1),
		(string)(arg.Arg_2),
		(string)(arg.Arg_3),
		(string)(arg.Arg_4),
		(string)(arg.Arg_5),
		(string)(arg.Arg_6),
		(float64)(arg.Arg_7),
		(string)(arg.Arg_8),
		([]*i_1.LHComponent)(DecodeArrayOfPtrToLHComponent(arg.Arg_9)),
		(float64)(arg.Arg_10),
		(float64)(arg.Arg_11),
		(*i_1.Shape)(DecodePtrToShape(arg.Arg_12)),
		(int)(arg.Arg_13),
		(float64)(arg.Arg_14),
		(float64)(arg.Arg_15),
		(float64)(arg.Arg_16),
		(float64)(arg.Arg_17),
		(string)(arg.Arg_18),
		(map[string]interface{})(DecodeMapstringinterfaceMessage(arg.Arg_19)),
		//(*i_1.LHPlate)(DecodePtrToLHPlate(arg.Arg_20)),//@jmanart gotopb cycle
		nil,
	}
	return ret
}
func EncodeMapstringPtrToLHWellMessage(arg map[string]*i_1.LHWell) *pb.MapstringPtrToLHWellMessageMessage {
	a := make([]*pb.MapstringPtrToLHWellMessageMessageFieldEntry, 0, len(arg))
	for k, v := range arg {
		fe := EncodeMapstringPtrToLHWellMessageFieldEntry(k, v)
		a = append(a, &fe)
	}
	ret := pb.MapstringPtrToLHWellMessageMessage{
		a,
	}
	return &ret
}
func EncodeMapstringPtrToLHWellMessageFieldEntry(k string, v *i_1.LHWell) pb.MapstringPtrToLHWellMessageMessageFieldEntry {
	ret := pb.MapstringPtrToLHWellMessageMessageFieldEntry{
		(string)(k),
		EncodePtrToLHWell(v),
	}
	return ret
}
func DecodeMapstringPtrToLHWellMessage(arg *pb.MapstringPtrToLHWellMessageMessage) map[string]*i_1.LHWell {
	a := make(map[(string)](*i_1.LHWell), len(arg.MapField))
	for _, fe := range arg.MapField {
		k, v := DecodeMapstringPtrToLHWellMessageFieldEntry(fe)
		a[k] = v
	}
	return a
}
func DecodeMapstringPtrToLHWellMessageFieldEntry(arg *pb.MapstringPtrToLHWellMessageMessageFieldEntry) (string, *i_1.LHWell) {
	k := (string)(arg.Key)
	v := DecodePtrToLHWell(arg.Value)
	return k, v
}
func EncodeArrayOfPtrToLHWell(arg []*i_1.LHWell) *pb.ArrayOfPtrToLHWellMessage {
	a := make([]*pb.PtrToLHWellMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodePtrToLHWell(v)
	}
	ret := pb.ArrayOfPtrToLHWellMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfPtrToLHWell(arg *pb.ArrayOfPtrToLHWellMessage) []*i_1.LHWell {
	ret := make(([]*i_1.LHWell), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodePtrToLHWell(v)
	}
	return ret
}
func EncodeArrayOfArrayOfPtrToLHWell(arg [][]*i_1.LHWell) *pb.ArrayOfArrayOfPtrToLHWellMessage {
	a := make([]*pb.ArrayOfPtrToLHWellMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodeArrayOfPtrToLHWell(v)
	}
	ret := pb.ArrayOfArrayOfPtrToLHWellMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfArrayOfPtrToLHWell(arg *pb.ArrayOfArrayOfPtrToLHWellMessage) [][]*i_1.LHWell {
	ret := make(([][]*i_1.LHWell), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodeArrayOfPtrToLHWell(v)
	}
	return ret
}
func EncodeLHComponent(arg i_1.LHComponent) *pb.LHComponentMessage {
	ret := pb.LHComponentMessage{}
	return &ret
}
func DecodeLHComponent(arg *pb.LHComponentMessage) i_1.LHComponent {
	ret := i_1.LHComponent{}
	return ret
}
func EncodePtrToShape(arg *i_1.Shape) *pb.PtrToShapeMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToShapeMessage{
		EncodeShape(*arg),
	}
	return &ret
}
func DecodePtrToShape(arg *pb.PtrToShapeMessage) *i_1.Shape {
	ret := DecodeShape(arg.Arg_1)
	return &ret
}
func EncodePtrToConcreteLocation(arg *i_1.ConcreteLocation) *pb.PtrToConcreteLocationMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToConcreteLocationMessage{
		EncodeConcreteLocation(*arg),
	}
	return &ret
}
func DecodePtrToConcreteLocation(arg *pb.PtrToConcreteLocationMessage) *i_1.ConcreteLocation {
	ret := DecodeConcreteLocation(arg.Arg_1)
	return &ret
}
func EncodeConcreteMeasurement(arg i_2.ConcreteMeasurement) *pb.ConcreteMeasurementMessage {
	ret := pb.ConcreteMeasurementMessage{(float64)(arg.Mvalue), EncodePtrToGenericPrefixedUnit(arg.Munit)}
	return &ret
}
func DecodeConcreteMeasurement(arg *pb.ConcreteMeasurementMessage) i_2.ConcreteMeasurement {
	ret := i_2.ConcreteMeasurement{(float64)(arg.Arg_1), (*i_2.GenericPrefixedUnit)(DecodePtrToGenericPrefixedUnit(arg.Arg_2))}
	return ret
}
func EncodeArrayOfPtrToLHComponent(arg []*i_1.LHComponent) *pb.ArrayOfPtrToLHComponentMessage {
	a := make([]*pb.PtrToLHComponentMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodePtrToLHComponent(v)
	}
	ret := pb.ArrayOfPtrToLHComponentMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfPtrToLHComponent(arg *pb.ArrayOfPtrToLHComponentMessage) []*i_1.LHComponent {
	ret := make(([]*i_1.LHComponent), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodePtrToLHComponent(v)
	}
	return ret
}
func EncodeShape(arg i_1.Shape) *pb.ShapeMessage {
	ret := pb.ShapeMessage{(string)(arg.ShapeName), (string)(arg.LengthUnit), (float64)(arg.H), (float64)(arg.W), (float64)(arg.D)}
	return &ret
}
func DecodeShape(arg *pb.ShapeMessage) i_1.Shape {
	ret := i_1.Shape{(string)(arg.Arg_1), (string)(arg.Arg_2), (float64)(arg.Arg_3), (float64)(arg.Arg_4), (float64)(arg.Arg_5)}
	return ret
}
func EncodeConcreteLocation(arg i_1.ConcreteLocation) *pb.ConcreteLocationMessage {
	ret := pb.ConcreteLocationMessage{(string)(arg.ID), (string)(arg.Inst), (string)(arg.Name), EncodeArrayOfPtrToConcreteLocation(arg.Psns), EncodePtrToConcreteLocation(arg.Cntr), EncodePtrToShape(arg.Shap)}
	return &ret
}
func DecodeConcreteLocation(arg *pb.ConcreteLocationMessage) i_1.ConcreteLocation {
	ret := i_1.ConcreteLocation{(string)(arg.Arg_1), (string)(arg.Arg_2), (string)(arg.Arg_3), ([]*i_1.ConcreteLocation)(DecodeArrayOfPtrToConcreteLocation(arg.Arg_4)), (*i_1.ConcreteLocation)(DecodePtrToConcreteLocation(arg.Arg_5)), (*i_1.Shape)(DecodePtrToShape(arg.Arg_6))}
	return ret
}
func EncodePtrToLHComponent(arg *i_1.LHComponent) *pb.PtrToLHComponentMessage {
	if arg == nil { //@jmanart fix when arg is nil, try to get a pointer
		return nil
	}
	ret := pb.PtrToLHComponentMessage{
		EncodeLHComponent(*arg),
	}
	return &ret
}
func DecodePtrToLHComponent(arg *pb.PtrToLHComponentMessage) *i_1.LHComponent {
	ret := DecodeLHComponent(arg.Arg_1)
	return &ret
}
func EncodeTemperature(arg i_2.Temperature) *pb.TemperatureMessage {
	ret := pb.TemperatureMessage{EncodeConcreteMeasurement(arg.ConcreteMeasurement)}
	return &ret
}
func DecodeTemperature(arg *pb.TemperatureMessage) i_2.Temperature {
	ret := i_2.Temperature{(i_2.ConcreteMeasurement)(DecodeConcreteMeasurement(arg.Arg_1))}
	return ret
}
func EncodeGenericPrefixedUnit(arg i_2.GenericPrefixedUnit) *pb.GenericPrefixedUnitMessage {
	ret := pb.GenericPrefixedUnitMessage{EncodeGenericUnit(arg.GenericUnit), EncodeSIPrefix(arg.SPrefix)}
	return &ret
}
func DecodeGenericPrefixedUnit(arg *pb.GenericPrefixedUnitMessage) i_2.GenericPrefixedUnit {
	ret := i_2.GenericPrefixedUnit{(i_2.GenericUnit)(DecodeGenericUnit(arg.Arg_1)), (i_2.SIPrefix)(DecodeSIPrefix(arg.Arg_2))}
	return ret
}
func EncodePtrToGenericPrefixedUnit(arg *i_2.GenericPrefixedUnit) *pb.PtrToGenericPrefixedUnitMessage {
	if arg == nil { //@jmanart fixed nil pointer
		return nil
	}
	ret := pb.PtrToGenericPrefixedUnitMessage{
		EncodeGenericPrefixedUnit(*arg),
	}
	return &ret
}
func DecodePtrToGenericPrefixedUnit(arg *pb.PtrToGenericPrefixedUnitMessage) *i_2.GenericPrefixedUnit {
	ret := DecodeGenericPrefixedUnit(arg.Arg_1)
	return &ret
}
func EncodeArrayOfPtrToConcreteLocation(arg []*i_1.ConcreteLocation) *pb.ArrayOfPtrToConcreteLocationMessage {
	a := make([]*pb.PtrToConcreteLocationMessage, len(arg))
	for i, v := range arg {
		a[i] = EncodePtrToConcreteLocation(v)
	}
	ret := pb.ArrayOfPtrToConcreteLocationMessage{
		a,
	}
	return &ret
}
func DecodeArrayOfPtrToConcreteLocation(arg *pb.ArrayOfPtrToConcreteLocationMessage) []*i_1.ConcreteLocation {
	ret := make(([]*i_1.ConcreteLocation), len(arg.Arg_1))
	for i, v := range arg.Arg_1 {
		ret[i] = DecodePtrToConcreteLocation(v)
	}
	return ret
}
func EncodeMass(arg i_2.Mass) *pb.MassMessage {
	ret := pb.MassMessage{EncodeConcreteMeasurement(arg.ConcreteMeasurement)}
	return &ret
}
func DecodeMass(arg *pb.MassMessage) i_2.Mass {
	ret := i_2.Mass{(i_2.ConcreteMeasurement)(DecodeConcreteMeasurement(arg.Arg_1))}
	return ret
}
func EncodeGenericUnit(arg i_2.GenericUnit) *pb.GenericUnitMessage {
	ret := pb.GenericUnitMessage{(string)(arg.StrName), (string)(arg.StrSymbol), (float64)(arg.FltConversionfactor), (string)(arg.StrBaseUnit)}
	return &ret
}
func DecodeGenericUnit(arg *pb.GenericUnitMessage) i_2.GenericUnit {
	ret := i_2.GenericUnit{(string)(arg.Arg_1), (string)(arg.Arg_2), (float64)(arg.Arg_3), (string)(arg.Arg_4)}
	return ret
}
func EncodeSpecificHeatCapacity(arg i_2.SpecificHeatCapacity) *pb.SpecificHeatCapacityMessage {
	ret := pb.SpecificHeatCapacityMessage{EncodeConcreteMeasurement(arg.ConcreteMeasurement)}
	return &ret
}
func DecodeSpecificHeatCapacity(arg *pb.SpecificHeatCapacityMessage) i_2.SpecificHeatCapacity {
	ret := i_2.SpecificHeatCapacity{(i_2.ConcreteMeasurement)(DecodeConcreteMeasurement(arg.Arg_1))}
	return ret
}
func EncodeSIPrefix(arg i_2.SIPrefix) *pb.SIPrefixMessage {
	ret := pb.SIPrefixMessage{(string)(arg.Name), (float64)(arg.Value)}
	return &ret
}
func DecodeSIPrefix(arg *pb.SIPrefixMessage) i_2.SIPrefix {
	ret := i_2.SIPrefix{(string)(arg.Arg_1), (float64)(arg.Arg_2)}
	return ret
}
func (d *Driver) AddPlateTo(arg_1 string, arg_2 interface{}, arg_3 string) i_4.CommandStatus {
	req := pb.AddPlateToRequest{
		(string)(arg_1),
		Encodeinterface(arg_2),
		(string)(arg_3),
	}
	ret, _ := d.C.AddPlateTo(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Aspirate(arg_1 []float64, arg_2 []bool, arg_3 int, arg_4 int, arg_5 []string, arg_6 []string, arg_7 []bool) i_4.CommandStatus {
	req := pb.AspirateRequest{
		EncodeArrayOffloat64(arg_1),
		EncodeArrayOfbool(arg_2),
		int64(arg_3),
		int64(arg_4),
		EncodeArrayOfstring(arg_5),
		EncodeArrayOfstring(arg_6),
		EncodeArrayOfbool(arg_7),
	}
	ret, _ := d.C.Aspirate(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Close() i_4.CommandStatus {
	req := pb.CloseRequest{}
	ret, _ := d.C.Close(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Dispense(arg_1 []float64, arg_2 []bool, arg_3 int, arg_4 int, arg_5 []string, arg_6 []string, arg_7 []bool) i_4.CommandStatus {
	req := pb.DispenseRequest{
		EncodeArrayOffloat64(arg_1),
		EncodeArrayOfbool(arg_2),
		int64(arg_3),
		int64(arg_4),
		EncodeArrayOfstring(arg_5),
		EncodeArrayOfstring(arg_6),
		EncodeArrayOfbool(arg_7),
	}
	ret, _ := d.C.Dispense(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Finalize() i_4.CommandStatus {
	req := pb.FinalizeRequest{}
	ret, _ := d.C.Finalize(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) GetCapabilities() (i_3.LHProperties, i_4.CommandStatus) {
	req := pb.GetCapabilitiesRequest{}
	ret, _ := d.C.GetCapabilities(context.Background(), &req)
	return (i_3.LHProperties)(DecodeLHProperties(ret.Ret_1)), (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) GetCurrentPosition(arg_1 int) (string, i_4.CommandStatus) {
	req := pb.GetCurrentPositionRequest{
		int64(arg_1),
	}
	ret, _ := d.C.GetCurrentPosition(context.Background(), &req)
	return (string)(ret.Ret_1), (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) GetHeadState(arg_1 int) (string, i_4.CommandStatus) {
	req := pb.GetHeadStateRequest{
		int64(arg_1),
	}
	ret, _ := d.C.GetHeadState(context.Background(), &req)
	return (string)(ret.Ret_1), (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) GetPositionState(arg_1 string) (string, i_4.CommandStatus) {
	req := pb.GetPositionStateRequest{
		(string)(arg_1),
	}
	ret, _ := d.C.GetPositionState(context.Background(), &req)
	return (string)(ret.Ret_1), (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) GetStatus() (i_4.Status, i_4.CommandStatus) {
	req := pb.GetStatusRequest{}
	ret, _ := d.C.GetStatus(context.Background(), &req)
	return (i_4.Status)(DecodeMapstringinterfaceMessage(ret.Ret_1)), (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_2))
}
func (d *Driver) Go() i_4.CommandStatus {
	req := pb.GoRequest{}
	ret, _ := d.C.Go(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Initialize() i_4.CommandStatus {
	req := pb.InitializeRequest{}
	ret, _ := d.C.Initialize(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LightsOff() i_4.CommandStatus {
	req := pb.LightsOffRequest{}
	ret, _ := d.C.LightsOff(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LightsOn() i_4.CommandStatus {
	req := pb.LightsOnRequest{}
	ret, _ := d.C.LightsOn(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LoadAdaptor(arg_1 int) i_4.CommandStatus {
	req := pb.LoadAdaptorRequest{
		int64(arg_1),
	}
	ret, _ := d.C.LoadAdaptor(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LoadHead(arg_1 int) i_4.CommandStatus {
	req := pb.LoadHeadRequest{
		int64(arg_1),
	}
	ret, _ := d.C.LoadHead(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) LoadTips(arg_1 []int, arg_2 int, arg_3 int, arg_4 []string, arg_5 []string, arg_6 []string) i_4.CommandStatus {
	req := pb.LoadTipsRequest{
		EncodeArrayOfint(arg_1),
		int64(arg_2),
		int64(arg_3),
		EncodeArrayOfstring(arg_4),
		EncodeArrayOfstring(arg_5),
		EncodeArrayOfstring(arg_6),
	}
	ret, _ := d.C.LoadTips(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Message(arg_1 int, arg_2 string, arg_3 string, arg_4 bool) i_4.CommandStatus {
	req := pb.MessageRequest{
		int64(arg_1),
		(string)(arg_2),
		(string)(arg_3),
		(bool)(arg_4),
	}
	ret, _ := d.C.Message(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Mix(arg_1 int, arg_2 []float64, arg_4 []string, arg_5 []int, arg_6 int, arg_7 []string, arg_8 []bool) i_4.CommandStatus {
	req := pb.MixRequest{
		int64(arg_1),
		EncodeArrayOffloat64(arg_2),
		EncodeArrayOfstring(arg_4),
		EncodeArrayOfint(arg_5),
		int64(arg_6),
		EncodeArrayOfstring(arg_7),
		EncodeArrayOfbool(arg_8),
	}
	ret, _ := d.C.Mix(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Move(arg_1 []string, arg_2 []string, arg_3 []int, arg_4 []float64, arg_5 []float64, arg_6 []float64, arg_7 []string, arg_8 int) i_4.CommandStatus {
	req := pb.MoveRequest{
		EncodeArrayOfstring(arg_1),
		EncodeArrayOfstring(arg_2),
		EncodeArrayOfint(arg_3),
		EncodeArrayOffloat64(arg_4),
		EncodeArrayOffloat64(arg_5),
		EncodeArrayOffloat64(arg_6),
		EncodeArrayOfstring(arg_7),
		int64(arg_8),
	}
	ret, _ := d.C.Move(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) MoveRaw(arg_1 int, arg_2 float64, arg_3 float64, arg_4 float64) i_4.CommandStatus {
	req := pb.MoveRawRequest{
		int64(arg_1),
		(float64)(arg_2),
		(float64)(arg_3),
		(float64)(arg_4),
	}
	ret, _ := d.C.MoveRaw(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Open() i_4.CommandStatus {
	req := pb.OpenRequest{}
	ret, _ := d.C.Open(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) RemoveAllPlates() i_4.CommandStatus {
	req := pb.RemoveAllPlatesRequest{}
	ret, _ := d.C.RemoveAllPlates(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) RemovePlateAt(arg_1 string) i_4.CommandStatus {
	req := pb.RemovePlateAtRequest{
		(string)(arg_1),
	}
	ret, _ := d.C.RemovePlateAt(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) ResetPistons(arg_1 int, arg_2 int) i_4.CommandStatus {
	req := pb.ResetPistonsRequest{
		int64(arg_1),
		int64(arg_2),
	}
	ret, _ := d.C.ResetPistons(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) SetDriveSpeed(arg_1 string, arg_2 float64) i_4.CommandStatus {
	req := pb.SetDriveSpeedRequest{
		(string)(arg_1),
		(float64)(arg_2),
	}
	ret, _ := d.C.SetDriveSpeed(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) SetPipetteSpeed(arg_1 int, arg_2 int, arg_3 float64) i_4.CommandStatus {
	req := pb.SetPipetteSpeedRequest{
		int64(arg_1),
		int64(arg_2),
		(float64)(arg_3),
	}
	ret, _ := d.C.SetPipetteSpeed(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) SetPositionState(arg_1 string, arg_2 i_4.PositionState) i_4.CommandStatus {
	req := pb.SetPositionStateRequest{
		(string)(arg_1),
		EncodeMapstringinterfaceMessage(arg_2),
	}
	ret, _ := d.C.SetPositionState(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Stop() i_4.CommandStatus {
	req := pb.StopRequest{}
	ret, _ := d.C.Stop(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) UnloadAdaptor(arg_1 int) i_4.CommandStatus {
	req := pb.UnloadAdaptorRequest{
		int64(arg_1),
	}
	ret, _ := d.C.UnloadAdaptor(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) UnloadHead(arg_1 int) i_4.CommandStatus {
	req := pb.UnloadHeadRequest{
		int64(arg_1),
	}
	ret, _ := d.C.UnloadHead(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) UnloadTips(arg_1 []int, arg_2 int, arg_3 int, arg_4 []string, arg_5 []string, arg_6 []string) i_4.CommandStatus {
	req := pb.UnloadTipsRequest{
		EncodeArrayOfint(arg_1),
		int64(arg_2),
		int64(arg_3),
		EncodeArrayOfstring(arg_4),
		EncodeArrayOfstring(arg_5),
		EncodeArrayOfstring(arg_6),
	}
	ret, _ := d.C.UnloadTips(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) UpdateMetaData(arg_1 *i_3.LHProperties) i_4.CommandStatus {
	req := pb.UpdateMetaDataRequest{
		EncodePtrToLHProperties(arg_1),
	}
	ret, _ := d.C.UpdateMetaData(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}
func (d *Driver) Wait(arg_1 float64) i_4.CommandStatus {
	req := pb.WaitRequest{
		(float64)(arg_1),
	}
	ret, _ := d.C.Wait(context.Background(), &req)
	return (i_4.CommandStatus)(DecodeCommandStatus(ret.Ret_1))
}

func (d *Driver) asExtendedLiquidhandlingDriver() liquidhandling.ExtendedLiquidhandlingDriver {
	var ret liquidhandling.ExtendedLiquidhandlingDriver
	ret = d
	return ret
}
