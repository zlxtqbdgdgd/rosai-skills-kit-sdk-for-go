package slu

import (
	"encoding/json"
	"log"
	"testing"
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func genTestIntent() *Intent {
	return NewIntent("PlanMyTrip").
		WithSlot(NewSlot("fromCity")).
		WithSlot(NewSlot("travelDate").WithStringValue("2018-04-11").WithStatus(CONFIRMED))
}

func TestIntentMerge(t *testing.T) {
	intent := genTestIntent()
	objIntent := NewIntent("PlanMyTrip").
		WithSlot(NewSlot("fromCity").WithStringValue("Beijing").WithStatus(CONFIRMED)).
		WithSlot(NewSlot("travelDate"))
	intent.Merge(objIntent)
	//bytes, _ := json.MarshalIndent(intent, "", "  ")
	//t.Log("intent merged: ", string(bytes))
	if intent.Slots["fromCity"].GetStringValue() != "Beijing" ||
		intent.Slots["fromCity"].ConfirmationStatus != CONFIRMED ||
		intent.Slots["travelDate"].GetStringValue() != "2018-04-11" ||
		intent.Slots["travelDate"].ConfirmationStatus != CONFIRMED {
		bytes, _ := json.MarshalIndent(intent, "", "  ")
		t.Fatalf("merge error, got: %+v", string(bytes))
	}
}

func TestIntentCanElicitConfirm(t *testing.T) {
	intent := genTestIntent()
	if !intent.CanElicit("fromCity") {
		t.Fatal("fromCity should be able to elicited")
	}
	if intent.CanConfirm("fromCity") {
		t.Fatal("fromCity should not be able to confirmed")
	}
	if intent.CanElicit("travelDate") {
		t.Fatal("travelDate should not be able to elicited")
	}
	if intent.CanConfirm("travelDate") {
		t.Fatal("travelDate should not be able to confirmed")
	}
	intent.GetSlot("travelDate").SetStatus(NONE)
	if !intent.CanConfirm("travelDate") {
		t.Fatal("travelDate should be able to confirmed")
	}
}

func TestGetStrMapValue(t *testing.T) {
	s := `{
  "latitude":39.919815,
  "longitude":116.43324,
  "address": {
    "country":"中国",
    "province":"北京市",
    "city":"北京市",
    "detail":"北京市 东城区 小牌坊胡同 靠近银河SOHO"
  }
}`
	v := NewStrMapValue(s)
	m, err := v.GetMapValue()
	if err != nil {
		t.Fatal(err)
	}
	if m["latitude"].(float64) != 39.919815 || m["longitude"] != 116.43324 {
		t.Fatalf("got: %#v, want: %s", m, s)
	}
}

func TestNewSetGetValue(t *testing.T) {
	// test int
	v := NewIntValue(10)
	r1, err := v.GetIntValue()
	if err != nil {
		t.Fatal(err)
	}
	if r1 != 10 {
		t.Fatalf("want int vlaue: 10, got: %d", r1)
	}
	// test bool
	v = NewBoolValue(true)
	r2, err := v.GetBoolValue()
	if err != nil {
		t.Fatal(err)
	}
	if r2 != true {
		t.Fatalf("want bool vlaue: true, got: %t", r2)
	}
	v = NewBoolValue(false)
	r2, err = v.GetBoolValue()
	if err != nil {
		t.Fatal(err)
	}
	if r2 != false {
		t.Fatalf("want bool vlaue: false, got: %t", r2)
	}
	// test float
	v = NewFloatValue(1.23)
	r3, err := v.GetFloatValue()
	if err != nil {
		t.Fatal(err)
	}
	if r3 != 1.23 {
		t.Fatalf("want int vlaue: 1.23, got: %f", r3)
	}
	v = NewFloatValue(-1.23)
	r3, err = v.GetFloatValue()
	if err != nil {
		t.Fatal(err)
	}
	if r3 != -1.23 {
		t.Fatalf("want int vlaue: -1.23, got: %f", r3)
	}
	v = NewFloatValue(0)
	r3, err = v.GetFloatValue()
	if err != nil {
		t.Fatal(err)
	}
	if r3 != 0 {
		t.Fatalf("want int vlaue: 0, got: %f", r3)
	}
	// test array
	v = NewArrayValue([]interface{}{1, 2, 3, "abc"})
	r4, err := v.GetArrayValue()
	if err != nil {
		t.Fatal(err)
	}
	if r4[0].(int) != 1 || r4[1].(int) != 2 || r4[2].(int) != 3 || r4[3].(string) != "abc" {
		t.Fatalf("want array %#v, got: %#v", v)
	}
	// test string array
	v = NewStrArrayValue([]string{"123", "aaa", "abc"})
	r5, err := v.GetStrArrayValue()
	if err != nil {
		t.Fatal(err)
	}
	if r5[0] != "123" || r5[1] != "aaa" || r5[2] != "abc" {
		t.Fatalf("want array %#v, got: %#v", "[123 aaa abc]", v)
	}
}

func TestGetSlotValue(t *testing.T) {
	s := `{
  "name": "artist",
  "value": {
    "orgin": [
      "刘德华"
    ],
    "normType": "StrArray",
    "norm": [
      "刘德华",
      "张学友"
    ],
    "priority": 99999902,
    "tag": "ARTIST"
  },
  "confirmationStatus": "NONE"
}`
	var slot Slot
	if err := json.Unmarshal([]byte(s), &slot); err != nil {
		t.Fatal(err)
	}
	vv := slot.GetStrArrayValue()
	if len(vv) != 2 || vv[0] != "刘德华" || vv[1] != "张学友" {
		t.Fatalf("want vv[刘德华, 张学友], got: %+v", vv)
	}
}
