package pkg

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path"
	"slices"
	"sort"
)

func gePkgListHome() (string, error) {
	pkgListHome := os.Getenv("PKG_LIST_HOME")
	if pkgListHome == "" {
		return "", errors.New("Environmental variable PKG_LIST_HOME isn't specified.")
	}

	_, lErr := os.Lstat(pkgListHome)
	if lErr != nil {
		if os.IsNotExist(lErr) {
			mErr := os.Mkdir(pkgListHome, 0666)
			if mErr != nil {
				return "", mErr
			}
		} else {
			return "", lErr
		}
	}
	return pkgListHome, nil
}

func getPkgFromFile(file string) ([][]byte, error) {
	buf, rErr := os.ReadFile(file)
	if rErr != nil {
		if os.IsNotExist(rErr) {
			return [][]byte{}, nil
		} else {
			return nil, rErr
		}
	}

	if len(buf) > 0 {
		return bytes.Split(buf, []byte{ '\n' }), nil
	}
	return [][]byte{}, nil
}

func AddPkg(pkg string, category string) (error) {
	if len(category) <= 0 {
		category = "other"
	}

	pkgListHome, cErr := gePkgListHome()
	if cErr != nil {
		return cErr
	}

	categoryFile := path.Join(pkgListHome, category)

	in, aErr := addPkg([]byte(pkg), categoryFile)
	if aErr != nil {
		return aErr
	}
	if in {
		fmt.Printf("Package %s was already in %s.\n", pkg, category)
	}
	return nil
}

func addPkg(pkg []byte, file string) (bool, error) {
	pkgs, err := getPkgFromFile(file)
	if err != nil {
		return false, err
	}

	pkgs, in := addPkgList([]byte(pkg), pkgs)
	if in {
		return in, nil
	}

	buf := bytes.Join(pkgs, []byte{ '\n' })
	tErr := os.WriteFile(file, buf, 0666)
	if tErr != nil {
		return false, tErr
	}

	return in, nil
}

func addPkgList(pkg []byte, pkgs [][]byte) ([][]byte, bool) {
	insertPos, found := sort.Find(len(pkgs), func(i int) int {
		return bytes.Compare([]byte(pkg), pkgs[i])
	})
	if !found {
		pkgs = slices.Insert(pkgs, insertPos, pkg)
	}
	return pkgs, found
}

func RmPkg(pkg string, category string) (error) {
	if len(category) <= 0 {
		category = "other"
	}

	pkgListHome, cErr := gePkgListHome()
	if cErr != nil {
		return cErr
	}

	categoryFile := path.Join(pkgListHome, category)

	in, rErr := rmPkg([]byte(pkg), categoryFile)
	if rErr != nil {
		return rErr
	}
	if !in {
		fmt.Printf("Package %s was not in %s.\n", pkg, category)
	}
	return nil
}

func rmPkg(pkg []byte, file string) (bool, error) {
	pkgs, err := getPkgFromFile(file)
	if err != nil {
		return false, err
	}

	pkgs, in := rmPkgList([]byte(pkg), pkgs)
	if !in {
		return in, nil
	}
	tErr := os.WriteFile(file, bytes.Join(pkgs, []byte{ '\n' }), 0666)
	if tErr != nil {
		return false, tErr
	}

	return in, nil
}

func rmPkgList(pkg []byte, pkgs [][]byte) ([][]byte, bool) {
	rmPos, found := sort.Find(len(pkgs), func(i int) int {
		return bytes.Compare([]byte(pkg), pkgs[i])
	})
	if found {
		pkgs = slices.Delete(pkgs, rmPos, rmPos + 1)
	}
	return pkgs, found
}
