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

type DialogHandler func(ctx context.Context, b *bot.Bot, update *models.Update, state int, info map[string]string) (newInfo map[string]string, isErr bool)

type stater interface {
	RegState(dialogType int)
	AddState(dialogType int, chatID int) error
	GetState(dialogType int, chatID int) (int, error) // return -1 if no dialog
	GetInfo(dialogType int, chatID int) (map[string]string, error)
	NextState(dialogType int, chatID int, NewInfo map[string]string) error
	DelState(dialogType int, chatID int) error
}

type Dialoger struct {
	st               stater
	dialogs          map[int]DialogHandler // dialogType -> dialog handler
	dialogsStatesNum map[int]int           // dialogType -> dialog states num
}

func New(st stater) *Dialoger {
	return &Dialoger{
		st:               st,
		dialogs:          make(map[int]DialogHandler),
		dialogsStatesNum: make(map[int]int),
	}
}

func (d *Dialoger) AddDialog(dialogType int, handler DialogHandler, statesNum int) {
	d.dialogs[dialogType] = handler
	d.dialogsStatesNum[dialogType] = statesNum
	d.st.RegState(dialogType)
}

// return dialogType and state of this dialog dialog if exist
func (d *Dialoger) CheckDialog(chatID int) (int, int, error) {
	for dialogType := range d.dialogs {
		state, err := d.st.GetState(dialogType, chatID)
		if err != nil {
			return UnknownDialog, 0, fmt.Errorf("failed to get state: %v", err)
		}
		if state != 0 {
			return dialogType, state, nil
		}
	}
	return UnknownDialog, 0, nil
}

func (d *Dialoger) StartDialog(ctx context.Context, b *bot.Bot, update *models.Update, dialogType int, chatID int, startInfo map[string]string) error {
	if _, isExist := d.dialogs[dialogType]; !isExist {
		return fmt.Errorf("no dialog of this type")
	}

	err := d.st.AddState(dialogType, chatID)
	if err != nil {
		return fmt.Errorf("failed to add state: %v", err)
	}

	handler := d.dialogs[dialogType]
	_, isErr := handler(ctx, b, update, 1, nil)

	if isErr {
		return nil
	}

	err = d.st.NextState(dialogType, chatID, startInfo)
	if err != nil {
		return fmt.Errorf("failed to next state: %v", err)
	}

	return nil
}

func (d *Dialoger) ServeMessage(ctx context.Context, b *bot.Bot, update *models.Update, dialogType int, state int) error {
	var (
		inputMsg = update.Message.Text
		chatID   = int(update.Message.Chat.ID)
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
	newInfo, isErr := handler(ctx, b, update, state, info)

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

// func DialogTypeIfExist(stater DialogStater, userID int) (bool, int, error) {
// 	reg, err := stater.GetState(RegDialog, userID)
// 	if err != nil {
// 		return false, 0, err
// 	}
// 	if reg != 0 {
// 		return true, RegDialog, nil
// 	}

// 	changeData, err := stater.GetState(ChangeDataDialog, userID)
// 	if err != nil {
// 		return false, 0, err
// 	}
// 	if changeData != 0 {
// 		return true, ChangeDataDialog, nil
// 	}

// 	leaveFeedback, err := stater.GetState(LeaveFeedbackDialog, userID)
// 	if err != nil {
// 		return false, 0, err
// 	}
// 	if leaveFeedback != 0 {
// 		return true, LeaveFeedbackDialog, nil
// 	}

// 	return false, 0, nil
// }
