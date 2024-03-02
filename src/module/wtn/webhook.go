package wtn

import (
	"WTN-Sniper/src/packages/file_manager"
	"WTN-Sniper/src/packages/notifier"
	"strconv"
	"time"
)

func TestWebhook(settings file_manager.SettingsData) bool {
	succResult := successWebhook(settings)
	errResult := errorWebhook(settings)

	if succResult && errResult {
		return succResult
	} else {
		return false
	}
}

func successWebhook(settings file_manager.SettingsData) bool {
	if settings.SuccessWebhook == "" {
		return false
	}

	var whData notifier.Discord_Params
	var tempEmbed notifier.Embed

	tempEmbed.Title = "Webhook tester"
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Type",
		Value:  "Success",
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Delay",
		Value:  strconv.Itoa(settings.Delay),
		Inline: true,
	})
	tempEmbed.Footer = notifier.Footer{
		Text:    "Kitsune WTN",
		IconURL: "https://media.discordapp.net/attachments/1198216317870817292/1199790308074979370/kitsune.png?ex=65c3d2cc&is=65b15dcc&hm=f3c166eae7f0b33c9af82656d85c73e5390c47cc8f984ebefb4ec74c3f0c65e2&=&format=webp&quality=lossless&width=1228&height=1228",
	}
	tempEmbed.Color = 366185
	tempEmbed.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	whData.Webhook_url = settings.SuccessWebhook
	whData.Webhook_data.Embeds = []notifier.Embed{tempEmbed}

	err := notifier.SendWebhook(whData)
	if err != nil {
		l.Error(err.Error())
	}
	return true
}

func errorWebhook(settings file_manager.SettingsData) bool {
	if settings.ErrorWebhook == "" {
		return false
	}

	var whData notifier.Discord_Params
	var tempEmbed notifier.Embed

	tempEmbed.Title = "Webhook tester"
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Type",
		Value:  "Error",
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Delay",
		Value:  strconv.Itoa(settings.Delay),
		Inline: true,
	})
	tempEmbed.Footer = notifier.Footer{
		Text:    "Kitsune WTN",
		IconURL: "https://media.discordapp.net/attachments/1198216317870817292/1199790308074979370/kitsune.png?ex=65c3d2cc&is=65b15dcc&hm=f3c166eae7f0b33c9af82656d85c73e5390c47cc8f984ebefb4ec74c3f0c65e2&=&format=webp&quality=lossless&width=1228&height=1228",
	}
	tempEmbed.Color = 14427686
	tempEmbed.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	whData.Webhook_url = settings.SuccessWebhook
	whData.Webhook_data.Embeds = []notifier.Embed{tempEmbed}

	err := notifier.SendWebhook(whData)
	if err != nil {
		l.Error(err.Error())
	}
	return true
}

func (s *WTNSession) listingWebhook() bool {
	var whData notifier.Discord_Params
	var tempEmbed notifier.Embed

	tempEmbed.Title = "Listings have been exported"
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Info",
		Value:  "Check the settings bot folder",
		Inline: true,
	})
	tempEmbed.Footer = notifier.Footer{
		Text:    "Kitsune WTN",
		IconURL: "https://media.discordapp.net/attachments/1198216317870817292/1199790308074979370/kitsune.png?ex=65c3d2cc&is=65b15dcc&hm=f3c166eae7f0b33c9af82656d85c73e5390c47cc8f984ebefb4ec74c3f0c65e2&=&format=webp&quality=lossless&width=1228&height=1228",
	}
	tempEmbed.Color = 366185
	tempEmbed.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	whData.Webhook_url = s.settingsData.SuccessWebhook
	whData.Webhook_data.Embeds = []notifier.Embed{tempEmbed}

	err := notifier.SendWebhook(whData)
	if err != nil {
		l.Error(err.Error())
	}
	return true
}

func (s *WTNSession) labelWebhook(sale Label_Result) bool {
	var whData notifier.Discord_Params
	var tempEmbed notifier.Embed

	tempEmbed.Title = "Label downloaded"
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Name",
		Value:  sale.Label_Product.Name,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Sku",
		Value:  sale.Label_Product.Sku,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Size",
		Value:  sale.Label_Product.Size,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Price",
		Value:  strconv.Itoa(int(sale.Label_Product.Price)),
		Inline: true,
	})
	tempEmbed.Thumbnail.URL = sale.Label_Product.Image
	tempEmbed.Footer = notifier.Footer{
		Text:    "Kitsune WTN",
		IconURL: "https://media.discordapp.net/attachments/1198216317870817292/1199790308074979370/kitsune.png?ex=65c3d2cc&is=65b15dcc&hm=f3c166eae7f0b33c9af82656d85c73e5390c47cc8f984ebefb4ec74c3f0c65e2&=&format=webp&quality=lossless&width=1228&height=1228",
	}
	tempEmbed.Color = 366185
	tempEmbed.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	whData.Webhook_url = s.settingsData.SuccessWebhook
	whData.Webhook_data.Embeds = []notifier.Embed{tempEmbed}

	err := notifier.SendWebhook(whData)
	if err != nil {
		l.Error(err.Error())
	}
	return true
}

func (s *WTNSession) extendWebhook(product Result) bool {
	var whData notifier.Discord_Params
	var tempEmbed notifier.Embed

	tempEmbed.Title = "List time extended to 60 days"
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Name",
		Value:  product.Product.Name,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Sku",
		Value:  product.Product.Sku,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Size",
		Value:  product.Product.EuropeanSize,
		Inline: true,
	})
	tempEmbed.Thumbnail.URL = product.Product.Image
	tempEmbed.Footer = notifier.Footer{
		Text:    "Kitsune WTN",
		IconURL: "https://media.discordapp.net/attachments/1198216317870817292/1199790308074979370/kitsune.png?ex=65c3d2cc&is=65b15dcc&hm=f3c166eae7f0b33c9af82656d85c73e5390c47cc8f984ebefb4ec74c3f0c65e2&=&format=webp&quality=lossless&width=1228&height=1228",
	}
	tempEmbed.Color = 366185
	tempEmbed.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	whData.Webhook_url = s.settingsData.SuccessWebhook
	whData.Webhook_data.Embeds = []notifier.Embed{tempEmbed}

	err := notifier.SendWebhook(whData)
	if err != nil {
		l.Error(err.Error())
	}
	return true
}

func (s *WTNSession) offerFoundWebhook(offer Offer_Result) bool {
	var whData notifier.Discord_Params
	var tempEmbed notifier.Embed

	tempEmbed.Title = "Offer found!"
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Name",
		Value:  offer.Name,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Sku",
		Value:  offer.Sku,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Size",
		Value:  offer.EuropeanSize,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Price",
		Value:  strconv.FormatInt(offer.Price, 10),
		Inline: true,
	})
	tempEmbed.Thumbnail.URL = offer.Image
	tempEmbed.Footer = notifier.Footer{
		Text:    "Kitsune WTN",
		IconURL: "https://media.discordapp.net/attachments/1198216317870817292/1199790308074979370/kitsune.png?ex=65c3d2cc&is=65b15dcc&hm=f3c166eae7f0b33c9af82656d85c73e5390c47cc8f984ebefb4ec74c3f0c65e2&=&format=webp&quality=lossless&width=1228&height=1228",
	}
	tempEmbed.Color = 366185
	tempEmbed.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	whData.Webhook_url = s.settingsData.SuccessWebhook
	whData.Webhook_data.Embeds = []notifier.Embed{tempEmbed}

	err := notifier.SendWebhook(whData)
	if err != nil {
		l.Error(err.Error())
	}
	return true
}

func (s *WTNSession) acceptWebhook(offer Offer_Result) bool {
	var whData notifier.Discord_Params
	var tempEmbed notifier.Embed

	tempEmbed.Title = "Offer accepted!"
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Name",
		Value:  offer.Name,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Sku",
		Value:  offer.Sku,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Size",
		Value:  offer.EuropeanSize,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Price",
		Value:  strconv.FormatInt(offer.Price, 10),
		Inline: true,
	})
	tempEmbed.Thumbnail.URL = offer.Image
	tempEmbed.Footer = notifier.Footer{
		Text:    "Kitsune WTN",
		IconURL: "https://media.discordapp.net/attachments/1198216317870817292/1199790308074979370/kitsune.png?ex=65c3d2cc&is=65b15dcc&hm=f3c166eae7f0b33c9af82656d85c73e5390c47cc8f984ebefb4ec74c3f0c65e2&=&format=webp&quality=lossless&width=1228&height=1228",
	}
	tempEmbed.Color = 366185
	tempEmbed.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	whData.Webhook_url = s.settingsData.SuccessWebhook
	whData.Webhook_data.Embeds = []notifier.Embed{tempEmbed}

	err := notifier.SendWebhook(whData)
	if err != nil {
		l.Error(err.Error())
	}
	return true
}

func (s *WTNSession) acceptErrorWebhook(offer Offer_Result, err_msg string) bool {
	var whData notifier.Discord_Params
	var tempEmbed notifier.Embed

	tempEmbed.Title = "Offer accepted!"
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Name",
		Value:  offer.Name,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Sku",
		Value:  offer.Sku,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Size",
		Value:  offer.EuropeanSize,
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Price",
		Value:  strconv.FormatInt(offer.Price, 10),
		Inline: true,
	})
	tempEmbed.Fields = append(tempEmbed.Fields, notifier.Field{
		Name:   "Error",
		Value:  err_msg,
		Inline: true,
	})
	tempEmbed.Thumbnail.URL = offer.Image
	tempEmbed.Footer = notifier.Footer{
		Text:    "Kitsune WTN",
		IconURL: "https://media.discordapp.net/attachments/1198216317870817292/1199790308074979370/kitsune.png?ex=65c3d2cc&is=65b15dcc&hm=f3c166eae7f0b33c9af82656d85c73e5390c47cc8f984ebefb4ec74c3f0c65e2&=&format=webp&quality=lossless&width=1228&height=1228",
	}
	tempEmbed.Color = 366185
	tempEmbed.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
	whData.Webhook_url = s.settingsData.SuccessWebhook
	whData.Webhook_data.Embeds = []notifier.Embed{tempEmbed}

	err := notifier.SendWebhook(whData)
	if err != nil {
		l.Error(err.Error())
	}
	return true
}
