package pkg

import (
	"bytes"
	"os"
	"sort"
	"testing"
)

func TestAddPkg(t *testing.T) {
	os.Setenv("PKG_LIST_HOME", ".")

	err := os.Remove("base")
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	err = AddPkg("clang", "base")
	if err != nil {
		t.Fatal(err)
	}

	err = AddPkg("clang", "base")
	if err != nil {
		t.Fatal(err)
	}

	err = AddPkg("python", "base")
	if err != nil {
		t.Fatal(err)
	}

	pkgs, err := getPkgFromFile("./base")
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) != 2 {
		t.Fail()
	}

	if _, f := sort.Find(len(pkgs), func(i int) int {
		return bytes.Compare([]byte("clang"), pkgs[i])
	}); !f {
		t.Fail()
	}
	if _, f := sort.Find(len(pkgs), func(i int) int {
		return bytes.Compare([]byte("python"), pkgs[i])
	}); !f {
		t.Fail()
	}
}

func TestRmPkg(t *testing.T) {
	os.Setenv("PKG_LIST_HOME", ".")

	err := os.Remove("base")
	if err != nil && !os.IsNotExist(err) {
		t.Fatal(err)
	}

	err = AddPkg("clang", "base")
	if err != nil {
		t.Fatal(err)
	}

	err = AddPkg("python", "base")
	if err != nil {
		t.Fatal(err)
	}

	err = RmPkg("clang", "base")
	if err != nil {
		t.Fatal(err)
	}

	err = RmPkg("clang", "base")
	if err != nil {
		t.Fatal(err)
	}

	err = RmPkg("python", "base")
	if err != nil {
		t.Fatal(err)
	}

	pkgs, err := getPkgFromFile("./base")
	if err != nil {
		t.Fatal(err)
	}

	if len(pkgs) != 0 {
		t.Fail()
	}
}
