package main

import (
	"testing"
)

func TestNoExpand(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha": { "component": "Alpha" } } }`)
	pdata := []byte(`{"Parameters": { "TheAlpha": { "Key": "Value" } } }`)
	if _, _, err := tryExpand(wdata, pdata); err != nil {
		t.Fatal(err)
	}
}

func TestExpand(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha": { "component": "Alpha" } } }`)
	pdata := []byte(`[ 
		{"Parameters": { "TheAlpha": { "Key": "Value" } } }, 
		{"Parameters": { "TheAlpha": { "Key": "Value" } } }
	]`)
	if desc, params, err := tryExpand(wdata, pdata); err != nil {
		t.Fatal(err)
	} else if l := len(desc.Processes); l != 2 {
		t.Errorf("expecting workflow of %d processes but found only %d", 2, l)
	} else if l := len(params.Parameters); l != 2 {
		t.Errorf("expecting parameters for %d processes but found only %d", 2, l)
	}
}

func TestBadExpand(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha": { "component": "Alpha" } } }`)
	pdata := []byte(`[ 
		{"Parameters": { "TheAlpha": { "Key": "Value" } }, "Config": { "A": "B" } }, 
		{"Parameters": { "TheAlpha": { "Key": "Value" } } }
	]`)
	if _, _, err := tryExpand(wdata, pdata); err == nil {
		t.Errorf("expecting error but got success")
	}
}

func TestSExpand(t *testing.T) {
	wdata := []byte(`{"processes": { "TheAlpha": { "component": "Alpha" } } }`)
	pdata := []byte(`[ 
		{ "TheAlpha": { "Key": "Value" } }, 
		{ "TheAlpha": { "Key": "Value" } }
	]`)
	if desc, params, err := tryExpand(wdata, pdata); err != nil {
		t.Fatal(err)
	} else if l := len(desc.Processes); l != 2 {
		t.Errorf("expecting workflow of %d processes but found only %d", 2, l)
	} else if l := len(params.Parameters); l != 2 {
		t.Errorf("expecting parameters for %d processes but found only %d", 2, l)
	}
}
