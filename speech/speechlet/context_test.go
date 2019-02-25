package speechlet

import (
	"encoding/json"
	"log"
	"testing"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
)

const (
	baseCtxStr = `{
  "system": {
    "skill": {
      "skillId": "rosai.skill.test001"
    },
    "user": {
      "userId": "rosai.user.test001",
      "appId": "rosai.app.test001"
    },
    "device": {
      "deviceId": "rosai.device.test001"
    }
  }
}`

	respEnvelope = `{
  "version": "2.0",
  "status": {
    "code": 0
  },
  "context": {
    "lifespanInMs": 3600000,
    "parameters": {
      "_lastInputSlots": {
        "orgin": null,
        "normType": "Map",
        "norm": {
          "album_name": {
            "confirmationStatus": "NONE",
            "name": "album_name",
            "value": {
              "norm": "FadeOri",
              "normType": "String",
              "orgin": null
            }
          }
        }
      },
      "_lastOutputSlots": {
        "orgin": null,
        "normType": "Map",
        "norm": {
          "album_name": {
            "confirmationStatus": "NONE",
            "name": "album_name",
            "value": {
              "norm": "Fade",
              "normType": "String",
              "orgin": null
            }
          },
          "artist": {
            "confirmationStatus": "NONE",
            "name": "artist",
            "value": {
              "norm": "Alan Walker ",
              "normType": "String",
              "orgin": null
            }
          },
          "name": {
            "confirmationStatus": "NONE",
            "name": "name",
            "value": {
              "norm": "Fade",
              "normType": "String",
              "orgin": null
            }
          }
        }
      },
      "_skillHit": {
        "orgin": null,
        "normType": "String",
        "norm": "Music"
      },
      "_subIntentHit": {
        "orgin": null,
        "normType": "String",
        "norm": "Play"
      }
    }
  },
  "results": []
}`
)

func init() {
	log.SetFlags(log.Ldate | log.Lmicroseconds | log.Lshortfile)
}

func TestNewGetSetContext(t *testing.T) {
	sysInfo := NewCtxSystem().
		WithUser(NewUser(userId, appId)).
		WithSkill(NewSkill(skillId)).
		WithDevice(NewDevice(deviceId))
	ctx := NewContext().WithSystem(sysInfo)
	bytes, _ := json.MarshalIndent(ctx, "", "  ")
	if baseCtxStr != string(bytes) {
		t.Fatalf("baseCtxStr want: %s, got: %s", baseCtxStr, string(bytes))
	}
	ctx.SetSysParameter("location",
		slu.NewValue(slu.MapType,
			map[string]interface{}{
				"latitude":  31.992714,
				"longitude": 118.773946,
				"address": struct {
					Country  string
					Province string
					City     string
					Detail   string
				}{
					Country:  "中国",
					Province: "江苏",
					City:     "南京",
					Detail:   "雨花台风景区",
				},
			}))
	ctx.SetSysParameter("hasScreen", slu.NewStringValue("true"))
	ctx.SetStringValue("city", "北京")
	bytes, _ = json.MarshalIndent(ctx, "", "  ")
	loc, ok := ctx.GetSysParameter("location").GetValue().(map[string]interface{})
	if !ok {
		t.Fatalf("location type assert failed, value: %+v, type: %T",
			ctx.GetSysParameter("location").GetValue(),
			ctx.GetSysParameter("location").GetValue())
	}
	//t.Logf("address value: %+v, type: %T", loc["address"], loc["address"])
	if ctx.GetStringValue("city") != "北京" ||
		ctx.GetSysStringValue("hasScreen") != "true" ||
		loc["latitude"].(float64) != 31.992714 ||
		loc["longitude"].(float64) != 118.773946 {
		//loc["address"].City != "南京" {
		t.Fatalf("context want: %s, got: %+v", string(bytes), ctx)
	}
	// delete "city"
	ctx.DelParameter("city")
	if ctx.GetStringValue("city") != "" {
		t.Fatalf("context delete city string error, got: %s", ctx.GetStringValue("city"))
	}
}

func TestSetGetLastInputSlots(t *testing.T) {
	slots := []*slu.Slot{
		slu.NewSlot("city").WithStringValue("北京"),
		slu.NewSlot("date").WithStringValue("2018-06-05"),
		slu.NewSlot("count").WithIntValue(6),
	}
	ctx := NewContext()
	ctx.SetLastOutputSlots(slots...)
	m := ctx.GetLastOutputSlots()
	if m["city"].GetStringValue() != "北京" ||
		m["date"].GetStringValue() != "2018-06-05" ||
		m["count"].GetIntValue() != 6 {
		bytes, _ := json.MarshalIndent(m, "", "  ")
		t.Fatalf("want: city[北京], date[2018-06-05], count[6], got: %s", bytes)
	}
}

func TestParseContextFromBytes(t *testing.T) {
	var respEnv ResponseEnvelope
	if err := json.Unmarshal([]byte(respEnvelope), &respEnv); err != nil {
		t.Fatal(err)
	}
	outSlots := respEnv.Context.GetLastOutputSlots()
	if len(outSlots) != 3 ||
		outSlots["album_name"].GetStringValue() != "Fade" ||
		outSlots["artist"].GetStringValue() != "Alan Walker " ||
		outSlots["name"].GetStringValue() != "Fade" {
		bytes, _ := json.MarshalIndent(outSlots, "", "  ")
		t.Fatalf("want: Fade, Alan Walker , Fade, got: %s", string(bytes))
	}
	inSlots := respEnv.Context.GetLastInputSlots()
	if len(inSlots) != 1 ||
		inSlots["album_name"].GetStringValue() != "FadeOri" {
		bytes, _ := json.MarshalIndent(inSlots, "", "  ")
		t.Fatalf("want: FadeOri, got: %s", string(bytes))
	}
	// skillHit
	skillHit := respEnv.Context.GetLastSkillHit()
	if skillHit != "Music" {
		t.Fatalf("skillHit want: Music, got: %s", skillHit)
	}
	// subIntentHit
	subIntentHit := respEnv.Context.GetLastSubIntentHit()
	if subIntentHit != "Play" {
		t.Fatalf("subIntentHit want: Play, got: %s", subIntentHit)
	}
	// delete sillHit
	respEnv.Context.DelInternalParameter(CtxParamKeyHitSkill)
	skillHit = respEnv.Context.GetLastSkillHit()
	if skillHit != "" {
		t.Fatalf("delete internal parameters skillHit failed,  want: , got: %s", skillHit)
	}
}
