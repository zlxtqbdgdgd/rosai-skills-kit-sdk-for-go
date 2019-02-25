package speechlet

/*
const (
	resultWithDisplay = `{
  "directives": [
    {
      "type": "Display.Customized",
      "hint": "hint text",
      "card": {
        "type": "Images",
        "list": [
          {
            "url": "pic.ros.ai/bot/sdk/standard.jpg",
            "bulletScreen": {
              "imageUrl": "pic.ros.ai/bot/sdk/bullet.jpg",
              "text": "bullet",
              "textPosition": "top",
              "imagePosition": "bottom"
            }
          }
        ]
      }
    }
  ]
}`

	resultWithDisplayAndEvent = `{
  "directives": [
    {
      "type": "Display.Customized",
      "hint": "hint text",
      "card": {
        "type": "Images",
        "list": [
          {
            "url": "pic.ros.ai/bot/sdk/standard.jpg",
            "bulletScreen": {
              "imageUrl": "pic.ros.ai/bot/sdk/bullet.jpg",
              "text": "bullet",
              "textPosition": "top",
              "imagePosition": "bottom"
            }
          }
        ]
      }
    },
    {
      "type": "ROSAI.EVENT",
      "event": {
        "name": "ROSAI.ContinueEvent",
        "period": 5000
      }
    }
  ]
}`

	resultWithTwoDisplayAndEvent = `{
  "directives": [
    {
      "type": "Display.Customized",
      "hint": "hint text",
      "card": {
        "type": "Images",
        "list": [
          {
            "url": "pic.ros.ai/bot/sdk/standard.jpg",
            "bulletScreen": {
              "imageUrl": "pic.ros.ai/bot/sdk/bullet.jpg",
              "text": "bullet",
              "textPosition": "top",
              "imagePosition": "bottom"
            }
          }
        ]
      }
    },
    {
      "type": "ROSAI.EVENT",
      "event": {
        "name": "ROSAI.ContinueEvent",
        "period": 5000
      }
    },
    {
      "type": "Display.Customized",
      "hint": "hint text",
      "card": {
        "type": "Images",
        "list": [
          {
            "url": "pic.ros.ai/bot/sdk/standard.jpg",
            "bulletScreen": {
              "imageUrl": "pic.ros.ai/bot/sdk/bullet.jpg",
              "text": "bullet",
              "textPosition": "top",
              "imagePosition": "bottom"
            }
          }
        ]
      }
    }
  ]
}`

	dispRaw = `{
  "type": "Display.Customized",
  "hint": "多云，气温23度到35度，东南风2级",
  "card": {
    "type": "Standard",
    "title": "北京天气",
    "content": "多云，气温23度到35度，东南风2级",
    "image": {
      "url": "www.roobo.com/aicloud/skills/weather/cloudy.jpg"
    }
  },
  "suggestions": [
    "明天",
    "后天"
  ]
}`
)

func TestDisplayDirective(t *testing.T) {
	var dir Directive
	result := NewResult()
	// test DisplayDirective
	dir = NewDisplayDirective().WithHint("hint text").
		WithCard(ui.NewImagesCard(ui.NewImage("pic.ros.ai/bot/sdk/standard.jpg").
			WithBullet(ui.NewBulletScreen().WithText("bullet", ui.Top).
				WithImage("pic.ros.ai/bot/sdk/bullet.jpg", ui.Bottom))))
	result.WithDisplayDirective(dir.(*DisplayDirective))
	bytes, _ := json.MarshalIndent(result, "", "  ")
	if resultWithDisplay != string(bytes) {
		t.Fatalf("images card want: %s, got: %s", resultWithDisplay, string(bytes))
	}
	displayDir := dir.(*DisplayDirective)
	// test EventDirective
	dir = NewEventDirective("ROSAI.ContinueEvent", 5000)
	result.WithEventDirective(dir.(*EventDirective))
	bytes, _ = json.MarshalIndent(result, "", "  ")
	if resultWithDisplayAndEvent != string(bytes) {
		t.Fatalf("images card want: %s, got: %s", resultWithDisplayAndEvent, string(bytes))
	}
	// test two DisplayDirective
	result.AppendDisplayDirective(displayDir)
	bytes, _ = json.MarshalIndent(result, "", "  ")
	if resultWithTwoDisplayAndEvent != string(bytes) {
		t.Fatalf("images card want: %s, got: %s", resultWithTwoDisplayAndEvent, string(bytes))
	}
}

func TestUnmarshalDisplayDirective(t *testing.T) {
	var ddr DisplayDirectiveRaw
	if err := json.Unmarshal([]byte(dispRaw), &ddr); err != nil {
		t.Fatal(err)
	}
	if ddr.Type != "Display.Customized" || ddr.Hint != "多云，气温23度到35度，东南风2级" ||
		ddr.Suggestions[0] != "明天" || ddr.Suggestions[1] != "后天" {
		t.Fatalf("want: %+v, got: %+v", dispRaw, ddr)
	}
	bytes, _ := json.Marshal(ddr)
	t.Logf("dispRaw: %s", string(bytes))
	card := ddr.GetCard()
	c := card.(*ui.StandardCard)
	if c.Type != "Standard" || c.Title != "北京天气" ||
		c.Content != "多云，气温23度到35度，东南风2级" ||
		c.Image.URL != "www.roobo.com/aicloud/skills/weather/cloudy.jpg" {
		t.Fatalf("got: %+v", c)
	}
	bytes, _ = json.Marshal(card)
	t.Logf("card: %s", string(bytes))
}
*/
