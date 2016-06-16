// Part of the Antha language
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

// Defines a platform independent mechanism for storing files generated when using antha.
// A folder  ./antha is produced in the home directory.
package anthapath

import (
	//"bufio"
	//"bytes"
	//"encoding/xml"
	"fmt"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/Parser"
	//"/data"
	//"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/enzymes"
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"io"
	"io/ioutil"
	"log"
	//"net/http"
	"os"
	"os/user"
	//"strings"
	//"time"
	"path/filepath"
)

func HomePath() (pathname string, err error) {
	u, err := user.Current()
	if err != nil {
		log.Panic(err)
	} else {
		pathname := u.HomeDir
		fmt.Println(pathname)
	}
	return
}

func CreatedotAnthafolder() (dpath string, err error) {

	if u, err := user.Current(); err != nil {
		log.Panic(err)
	} else if dpath = filepath.Join(u.HomeDir, ".antha"); false {
	} else if err := os.MkdirAll(dpath, 0755); err != nil {
		log.Panic(err)
	}
	fmt.Printf("Created folder %s\n", dpath)

	return
}

func AddFile(filename string) (f *os.File, err error) {
	_, err = CreatedotAnthafolder()
	if err != nil {
		log.Panic(err)
	}
	if u, err := user.Current(); err != nil {
		log.Panic(err)
	} else if dpath := filepath.Join(u.HomeDir, ".antha"); false {
	} else if err := os.MkdirAll(dpath, 0755); err != nil {
		log.Panic(err)
	} else if f, err := os.Create(filepath.Join(dpath, filename)); err != nil {
		log.Panic(err)
	} else {
		defer f.Close()
		fmt.Printf("Created file %s\n", f.Name())
	}
	return
}

func ExporttoFile(filename string, contents []byte) (err error) {

	err = ioutil.WriteFile(filepath.Join(Dirpath(), filename), contents, os.ModeDir)
	fmt.Println("lining")
	if err != nil {
		fmt.Println("3")
		log.Fatal(err)
	}

	return err
}

func ExportTextFile(filename string, contents string) (err error) {
	_, err = AddFile(filename)
	if err != nil {
		log.Panic(err)
	}
	f, err := os.Open(filepath.Join(Dirpath(), filename))
	if err != nil {
		log.Panic(err)
	}
	defer f.Close()
	_, err = fmt.Fprint(f, contents)

	if err != nil {
		log.Fatal(err)
	}

	return err
}

func Dirpath() (dirpath string) {
	u, err := user.Current()
	if err != nil {
		log.Panic(err)
	}
	dirpath = filepath.Join(u.HomeDir, ".antha")
	return
}

func AnthaFile(filename string) (anthapathandfilename string) {

	anthapathandfilename = filepath.Join(Dirpath(), filename)

	return
}

func Anthafileexists(filename string) bool {
	fullpath := filepath.Join(Dirpath(), filename)
	if Exists(fullpath) {
		return true
	}
	return false
}

func Exists(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
