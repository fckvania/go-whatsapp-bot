package message

import (
	"context"
	"strings"

	"go.vnia.dev/helper"
	"go.vnia.dev/lib"
	"google.golang.org/protobuf/proto"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types/events"
)

// Config
var (
	prefix = "!"
	self   = true
	owner  = "6281236031617"
)

func Msg(client *whatsmeow.Client, msg *events.Message) {
	// simple
	simp := lib.NewSimpleImpl(client, msg)
	// dll
	from := msg.Info.Chat
	sender := msg.Info.Sender.String()
	args := strings.Split(simp.GetCMD(), " ")
	command := strings.ToLower(args[0])
	//query := strings.Join(args[1:], ` `)
	pushName := msg.Info.PushName
	isOwner := strings.Contains(sender, owner)
	//isAdmin := simp.GetGroupAdmin(from, sender)
	//isGroup := msg.Info.IsGroup
	extended := msg.Message.GetExtendedTextMessage()
	quotedMsg := extended.GetContextInfo().GetQuotedMessage()
	quotedImage := quotedMsg.GetImageMessage()
	//quotedVideo := quotedMsg.GetVideoMessage()
	//quotedSticker := quotedMsg.GetStickerMessage()
	// Self
	if self && !isOwner {
		return
	}
	// Switch Cmd
	switch command {
	case prefix + "menu":
		buttons := []*waProto.HydratedTemplateButton{
			{
				HydratedButton: &waProto.HydratedTemplateButton_QuickReplyButton{
					QuickReplyButton: &waProto.HydratedQuickReplyButton{
						DisplayText: proto.String("OWNER"),
						Id:          proto.String(prefix + "owner"),
					},
				},
			},
		}
		simp.SendHydratedBtn(from, helper.Menu(pushName, prefix), "Author : Vnia\nLibrary : Whatsmeow", buttons)
	case prefix + "owner":
		simp.SendContact(from, owner, "vnia")
	case prefix + "sticker":
		if quotedImage != nil {
			data, _ := client.Download(quotedImage)
			stc := simp.CreateStickerIMG(data)
			client.SendMessage(context.Background(), from, "", stc)
		} else if msg.Message.GetImageMessage() != nil {
			data, _ := client.Download(msg.Message.GetImageMessage())
			stc := simp.CreateStickerIMG(data)
			client.SendMessage(context.Background(), from, "", stc)
		}
	}
}
