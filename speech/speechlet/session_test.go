package speechlet

import (
	"testing"

	"roobo.com/rosai-skills-kit-sdk-for-go/speech/slu"
)

var (
	userId   = "rosai.user.test001"
	appId    = "rosai.app.test001"
	deviceId = "rosai.device.test001"
	skillId  = "rosai.skill.test001"

	intent = slu.NewIntent("PlanMyTrip").WithStatus(slu.CONFIRMED).
		WithSlot(slu.NewSlot("toCity").WithStringValue("Sanya").WithStatus(slu.CONFIRMED)).
		WithSlot(slu.NewSlot("travelDate").WithStringValue("2018-04-05"))
)

func TestNewSession(t *testing.T) {
	ss := NewSession(userId, appId, deviceId, skillId)
	t.Logf("session: %+v", ss)
}

func TestNewRediSession(t *testing.T) {
	GetRediSession()
}

func TestSessionOperate(t *testing.T) {
	// test New
	rediSS := GetRediSession()
	// save
	ss := NewSession(userId, appId, deviceId, skillId).WithUpdatedIntent(intent)
	ss.SetNew(false)
	t.Logf("session: %+v", ss)
	if err := rediSS.Save(ss); err != nil {
		t.Fatal(err)
	}
	// fetch
	ssGot, _ := rediSS.Fetch(userId, appId, deviceId, skillId)
	if ssGot.New {
		t.Fatal("session got field new should be false, now got true")
	}
	t.Logf("session got: %+v", ssGot)
	intentPre := ss.GetUpdatedIntent("PlanMyTrip")
	intentPost := ssGot.GetUpdatedIntent("PlanMyTrip")
	if ss.ID != ssGot.ID || ss.User.UserId != ssGot.User.UserId ||
		intentPre.GetSlot("travelDate").GetStringValue() !=
			intentPost.GetSlot("travelDate").GetStringValue() {
		t.Fatalf("want: %+v, fetch got: %+v", ss, ssGot)
	}
	if ssGot.New {
		t.Fatal("ssesion got field New should be false, now got true")
	}
	// drop
	if err := rediSS.Drop(userId, appId, deviceId, skillId); err != nil {
		t.Fatal(err)
	}
	// test fetch
	ssGot, _ = rediSS.Fetch(userId, appId, deviceId, skillId)
	if !ssGot.New {
		t.Fatal("session got field new should be true after drop, now got false")
	}
	if len(ssGot.Attributes) > 0 {
		t.Fatalf("session got after drop Attributes should be empty, now got: %+v",
			ssGot.Attributes)
	}
	t.Logf("session got: %+v", ssGot)
}
