package message

import (
	"strings"

	"go.vnia.dev/helper"
	"go.vnia.dev/lib"

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
	pushName := msg.Info.PushName
	isOwner := strings.Contains(sender, owner)
	// Self
	if self && !isOwner {
		return
	}
	// Switch Cmd
	switch strings.ToLower(simp.GetCMD()) {
	case prefix + "menu":
		simp.SendHydratedBtn(from, helper.Menu(pushName, prefix), "Author : Vnia\nLibrary : Whatsmeow", []*waProto.HydratedTemplateButton{})
	}
}
