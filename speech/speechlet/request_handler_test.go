package speechlet

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/dialog/model"
)

var (
	rh *RequestHandler
)

func init() {
	dm, err := getDialogModel()
	if err != nil {
		log.Fatal(err)
	}
	rh = &RequestHandler{
		AppId:       "rosai1.ask.skill.test.12345",
		DialogModel: dm,
	}
}

func TestDispatch(t *testing.T) {

}

func getDialogModel() (*model.DialogModel, error) {
	data, err := ioutil.ReadFile("../conf/dialog_test.json")
	if err != nil {
		return nil, err
	}
	dm := new(model.DialogModel)
	if err := json.Unmarshal(data, dm); err != nil {
		return nil, err
	}
	return dm, nil
}

func TestTryResolveParams(t *testing.T) {
	var unresolved string = "niaho{$city}12{$date}你好123"
	var paramMap = make(map[string]string)
	paramMap["city"] = "北京"
	paramMap["date"] = "明天"
	_tryResolveParams(&unresolved, paramMap)
}