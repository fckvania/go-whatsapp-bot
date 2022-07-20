package lib

import (
	"context"
	"fmt"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

type SimpleImpl struct {
	VClient *whatsmeow.Client
	Msg     *events.Message
}

func NewSimpleImpl(Cli *whatsmeow.Client, m *events.Message) *SimpleImpl {
	return &SimpleImpl{
		VClient: Cli,
		Msg:     m,
	}
}

func (simp *SimpleImpl) Reply(teks string) {
	simp.VClient.SendMessage(context.Background(), simp.Msg.Info.Chat, "", &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(teks),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &simp.Msg.Info.ID,
				Participant:   proto.String(simp.Msg.Info.Sender.String()),
				QuotedMessage: simp.Msg.Message,
			},
		},
	})
}

func (simp *SimpleImpl) SendHydratedBtn(jid types.JID, teks string, foter string, buttons []*waProto.HydratedTemplateButton) {
	simp.VClient.SendMessage(context.Background(), jid, "", &waProto.Message{
		TemplateMessage: &waProto.TemplateMessage{
			HydratedTemplate: &waProto.HydratedFourRowTemplate{
				HydratedContentText: proto.String(teks),
				HydratedFooterText:  proto.String(foter),
				HydratedButtons:     buttons,
			},
		},
	})
}

func (simp *SimpleImpl) SendContact(jid types.JID, number string, nama string) {
	simp.VClient.SendMessage(context.Background(), jid, "", &waProto.Message{
		ContactMessage: &waProto.ContactMessage{
			DisplayName: proto.String(nama),
			Vcard:       proto.String(fmt.Sprintf("BEGIN:VCARD\nVERSION:3.0\nN:%s;;;\nFN:%s\nitem1.TEL;waid=%s:+%s\nitem1.X-ABLabel:Mobile\nEND:VCARD", nama, nama, number, number)),
			ContextInfo: &waProto.ContextInfo{
				StanzaId:      &simp.Msg.Info.ID,
				Participant:   proto.String(simp.Msg.Info.Sender.String()),
				QuotedMessage: simp.Msg.Message,
			},
		},
	})
}

func (simp *SimpleImpl) FetchGroupAdmin(Jid types.JID) ([]string, error) {
	var Admin []string
	resp, err := simp.VClient.GetGroupInfo(Jid)
	if err != nil {
		return Admin, err
	} else {
		for _, group := range resp.Participants {
			if group.IsAdmin || group.IsSuperAdmin {
				Admin = append(Admin, group.JID.String())
			}
		}
	}
	return Admin, nil
}

func (simp *SimpleImpl) GetGroupAdmin(jid types.JID, sender string) bool {
	if !simp.Msg.Info.IsGroup {
		return false
	}
	admin, err := simp.FetchGroupAdmin(jid)
	if err != nil {
		return false
	}
	for _, v := range admin {
		if v == sender {
			return true
		}
	}
	return false
}

func (simp *SimpleImpl) GetCMD() string {
	extended := simp.Msg.Message.GetExtendedTextMessage().GetText()
	text := simp.Msg.Message.GetConversation()
	imageMatch := simp.Msg.Message.GetImageMessage().GetCaption()
	videoMatch := simp.Msg.Message.GetVideoMessage().GetCaption()
	tempBtnId := simp.Msg.Message.GetTemplateButtonReplyMessage().GetSelectedId()
	var command string
	if text != "" {
		command = text
	} else if imageMatch != "" {
		command = imageMatch
	} else if videoMatch != "" {
		command = videoMatch
	} else if extended != "" {
		command = extended
	} else if tempBtnId != "" {
		command = tempBtnId
	}
	return command
}
