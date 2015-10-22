package swagger2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// getJsonFiles returns a list of files from the test folder
func getTestFiles(ext string) ([]string, error) {
	fil, err := os.Open("examples/" + ext)
	if err != nil {
		return nil, fmt.Errorf("Unable to open test folder: %s", err)
	}
	defer fil.Close()

	var info []os.FileInfo
	result := make([]string, 0)

	info, err = fil.Readdir(256)

	for err == nil {
		for i := range info {
			name := info[i].Name()
			if strings.HasSuffix(name, "."+ext) {
				result = append(result, filepath.Join("examples/"+ext, name))
			}
		}
		info, err = fil.Readdir(256)
	}

	return result, nil
}

func cheesyCompare(ext string, buf1, buf2 []byte) bool {
	// remove CR bytes
	buf1 = bytes.Replace(buf1, []byte("\r"), []byte(""), -1)
	buf2 = bytes.Replace(buf2, []byte("\r"), []byte(""), -1)

	// remove commas, which won't be consistent
	buf1 = bytes.Replace(buf1, []byte(","), []byte(""), -1)
	buf2 = bytes.Replace(buf2, []byte(","), []byte(""), -1)

	// split into strings
	arr1 := strings.Split(string(buf1), "\n")
	arr2 := strings.Split(string(buf2), "\n")

	// rmeove empty lines
	xarr1 := make([]string, 0)
	for _, s := range arr1 {
		st := strings.TrimSpace(s)
		if st != "" && st[0] != '#' {
			if ext == "yaml" {
				s = st
			}
			xarr1 = append(xarr1, s)
		}
	}
	xarr2 := make([]string, 0)
	for _, s := range arr2 {
		st := strings.TrimSpace(s)
		if st != "" && st[0] != '#' {
			if ext == "yaml" {
				s = st
			}
			xarr2 = append(xarr2, s)
		}
	}

	// sort lines
	sort.Strings(xarr1)
	sort.Strings(xarr2)

	// reassmeble
	str1 := strings.Join(xarr1, "\n")
	str2 := strings.Join(xarr2, "\n")

	// byte compare
	x := bytes.Compare([]byte(str1), []byte(str2))
	if x == 0 {
		return true
	}
	return false
}

func TestJsonFiles(t *testing.T) {
	files, err := getTestFiles("json")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, file := range files {
		buf, err := ioutil.ReadFile(file)
		swag, err := LoadJson(buf)
		if err != nil {
			t.Errorf("Unable to parse file \"%s\": %s", file, err)
		} else {
			buf2, err := swag.Json()
			if err != nil {
				t.Errorf("Unable to parse file \"%s\": %s", file, err)
			} else {
				errs := swag.Validate()
				if len(errs) > 0 {
					t.Error("Swagger does not validate: ", ErrorList(errs).String())
				}
				if cheesyCompare("json", buf, buf2) != true {
					t.Error("Reserialized data does not match original. See", file+".new")
					ioutil.WriteFile(file+".new", buf2, 0666)
				}
			}
		}
	}
}
