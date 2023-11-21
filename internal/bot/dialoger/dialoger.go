package dialoger

import (
	"context"
	"fmt"

	"github.com/GrishaSkurikhin/DivanBot/internal/bot/commands"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

const (
	UnknownDialog       = -1
	RegDialog           = 0
	ChangeDataDialog    = 1
	LeaveFeedbackDialog = 2
)

type DialogHandler func(ctx context.Context, b *bot.Bot, msg *models.Message, state int, info map[string]string) (newInfo map[string]string, isErr bool)

type Stater interface {
	AddState(dialogType int, chatID int64) error
	GetState(chatID int64) (int, int, error)
	GetInfo(dialogType int, chatID int64) (map[string]string, error)
	NextState(dialogType int, chatID int64, NewInfo map[string]string) error
	DelState(dialogType int, chatID int64) error
}

type Dialoger struct {
	st               Stater
	dialogs          map[int]DialogHandler // dialogType -> dialog handler
	dialogsStatesNum map[int]int           // dialogType -> dialog states num
}

func New(st Stater) *Dialoger {
	return &Dialoger{
		st:               st,
		dialogs:          make(map[int]DialogHandler),
		dialogsStatesNum: make(map[int]int),
	}
}

func (d *Dialoger) AddDialog(dialogType int, handler DialogHandler, statesNum int) {
	d.dialogs[dialogType] = handler
	d.dialogsStatesNum[dialogType] = statesNum
}

// return dialogType and state of this dialog dialog if exist
func (d *Dialoger) CheckDialog(chatID int64) (int, int, error) {
	dialogType, state, err := d.st.GetState(chatID)
	if err != nil {
		return UnknownDialog, 0, fmt.Errorf("failed to get state: %v", err)
	}
	if state != 0 {
		return dialogType, state, nil
	}

	return UnknownDialog, 0, nil
}

func (d *Dialoger) StartDialog(ctx context.Context, b *bot.Bot, msg *models.Message, dialogType int, chatID int64, startInfo map[string]string) error {
	if _, isExist := d.dialogs[dialogType]; !isExist {
		return fmt.Errorf("no dialog of this type")
	}

	err := d.st.AddState(dialogType, chatID)
	if err != nil {
		return fmt.Errorf("failed to add state: %v", err)
	}

	handler := d.dialogs[dialogType]
	_, isErr := handler(ctx, b, msg, 1, nil)

	if isErr {
		return nil
	}

	err = d.st.NextState(dialogType, chatID, startInfo)
	if err != nil {
		return fmt.Errorf("failed to next state: %v", err)
	}

	return nil
}

func (d *Dialoger) ServeMessage(ctx context.Context, b *bot.Bot, msg *models.Message, dialogType int, state int) error {
	var (
		inputMsg = msg.Text
		chatID   = msg.Chat.ID
	)

	if _, isExist := d.dialogs[dialogType]; !isExist {
		return fmt.Errorf("no dialog of this type")
	}

	if inputMsg == commands.Cancel {
		err := d.st.DelState(dialogType, chatID)
		if err != nil {
			return fmt.Errorf("failed to del state: %v", err)
		}
		return nil
	}

	info, err := d.st.GetInfo(dialogType, chatID)
	if err != nil {
		return fmt.Errorf("no dialog of this type")
	}

	handler := d.dialogs[dialogType]
	newInfo, isErr := handler(ctx, b, msg, state, info)

	if isErr {
		return nil
	}

	if state == d.dialogsStatesNum[dialogType] {
		err := d.st.DelState(dialogType, chatID)
		if err != nil {
			return fmt.Errorf("failed to del state: %v", err)
		}
		return nil
	}

	err = d.st.NextState(dialogType, chatID, newInfo)
	if err != nil {
		return fmt.Errorf("failed to next state: %v", err)
	}

	return nil
}