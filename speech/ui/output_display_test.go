package ui

import (
	"encoding/json"
	"testing"
)

const (
	textCardOutput = `{
  "type": "Text",
  "title": "simple text card",
  "content": "this is a simple text card"
}`

	standardCardOutput = `{
  "type": "Standard",
  "title": "standard card",
  "content": "hello standard card",
  "image": {
    "url": "images.ros.ai/bot/skills/sdk/card/standardcard.jpg",
    "bulletScreen": {
      "imageUrl": "pic.ros.ai/bot/skills/sdk/bullet.jpg",
      "text": "BulletScreen",
      "textPosition": "bottom",
      "imagePosition": "top"
    }
  }
}`

	imagesCardOutput = `{
  "type": "Images",
  "list": [
    {
      "url": "images.ros.ai/bot/skills/sdk/card/image1.jpg"
    },
    {
      "url": "images.ros.ai/bot/skills/sdk/card/image2.jpg",
      "bulletScreen": {
        "imageUrl": "pic.ros.ai/bot/skills/sdk/bullet.jpg",
        "text": "5 stars",
        "textPosition": "bottom",
        "imagePosition": "bottom"
      }
    },
    {
      "url": "images.ros.ai/bot/skills/sdk/card/image3.jpg"
    }
  ]
}`

	listCardOutput = `{
  "type": "List",
  "list": [
    {
      "type": "Standard",
      "title": "card1",
      "content": "content1",
      "image": {
        "url": "images.ros.ai/bot/skills/sdk/card/image1.jpg"
      }
    },
    {
      "type": "Standard",
      "title": "card2",
      "content": "content2",
      "image": {
        "url": "images.ros.ai/bot/skills/sdk/card/image2.jpg"
      }
    },
    {
      "type": "Standard",
      "title": "card3",
      "content": "content3",
      "image": {
        "url": "images.ros.ai/bot/skills/sdk/card/image3.jpg"
      }
    }
  ]
}`
)

func TestMakeCards(t *testing.T) {
	var card CardInterface
	// text card
	card = NewTextCard("simple text card", "this is a simple text card")
	bytes, _ := json.MarshalIndent(card, "", "  ")
	if textCardOutput != string(bytes) {
		t.Fatalf("StandardCard want: %s, got: %s", textCardOutput, string(bytes))
	}
	// standard card
	card = NewStandardCard("standard card", "hello standard card").
		WithBulletImage("images.ros.ai/bot/skills/sdk/card/standardcard.jpg",
			NewBulletScreen().WithImage("pic.ros.ai/bot/skills/sdk/bullet.jpg", Top).
				WithText("BulletScreen", Bottom))
	bytes, _ = json.MarshalIndent(card, "", "  ")
	if standardCardOutput != string(bytes) {
		t.Fatalf("StandardCard want: %s, got: %s", standardCardOutput, string(bytes))
	}
	// images card
	card = NewImagesCard(NewImage("images.ros.ai/bot/skills/sdk/card/image1.jpg"),
		NewImage("images.ros.ai/bot/skills/sdk/card/image2.jpg").
			WithBullet(NewBulletScreen().
				WithImage("pic.ros.ai/bot/skills/sdk/bullet.jpg", Bottom).
				WithText("5 stars", Bottom))).
		AppendImages(NewImage("images.ros.ai/bot/skills/sdk/card/image3.jpg"))
	bytes, _ = json.MarshalIndent(card, "", "  ")
	if imagesCardOutput != string(bytes) {
		t.Fatalf("images card want: %s, got: %s", imagesCardOutput, string(bytes))
	}
	// list card
	card = NewListCard(
		NewStandardCard("card1", "content1").WithImage("images.ros.ai/bot/skills/sdk/card/image1.jpg"),
		NewStandardCard("card2", "content2").WithImage("images.ros.ai/bot/skills/sdk/card/image2.jpg")).
		AppendCards(NewStandardCard("card3", "content3").WithImage("images.ros.ai/bot/skills/sdk/card/image3.jpg"))
	bytes, _ = json.MarshalIndent(card, "", "  ")
	if listCardOutput != string(bytes) {
		t.Fatalf("list card want: %s, got: %s", listCardOutput, string(bytes))
	}
}
