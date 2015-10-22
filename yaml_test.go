package swagger2

import (
	"io/ioutil"
	"testing"
)

func TestYamlFiles(t *testing.T) {
	files, err := getTestFiles("yaml")
	if err != nil {
		t.Log(err)
		t.FailNow()
	}

	for _, file := range files {
		buf, err := ioutil.ReadFile(file)
		swag, err := LoadYaml(buf)
		if err != nil {
			t.Errorf("Unable to parse file \"%s\": %s", file, err)
		} else {
			buf2, err := swag.Yaml()
			if err != nil {
				t.Errorf("Unable to parse file \"%s\": %s", file, err)
			} else {
				errs := swag.Validate()
				if len(errs) > 0 {
					t.Error("Swagger does not validate: ", ErrorList(errs).String())
				}
				//  YAML file compare is a mess
				if cheesyCompare("yaml", buf, buf2) != true {
					t.Log("Reserialized data does not match original. See", file+".new")
					ioutil.WriteFile(file+".new", buf2, 0666)
				}
			}
		}
	}
}
