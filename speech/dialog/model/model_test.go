package model

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestLoadModelJson(t *testing.T) {
	data, err := ioutil.ReadFile("./test_dialog.json")
	if err != nil {
		t.Fatal(err)
	}
	var dm DialogModel
	if err := json.Unmarshal(data, &dm); err != nil {
		t.Fatal(err)
	}
	_, err = json.MarshalIndent(dm, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	/*if strings.TrimSpace(string(data)) != strings.TrimSpace(string(data2)) {
		t.Fatalf("bytes load and marshal model json is not same, origin:\n%sgot:\n%s",
			data, data2)
	}*/
}
