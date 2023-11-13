package mapstater

import (
	"fmt"

	"github.com/GrishaSkurikhin/DivanBot/internal/models"
)

type mapStater struct {
	states map[int]map[int]models.State // dialogType->userID->State
}

func New() *mapStater {
	return &mapStater{
		states: make(map[int]map[int]models.State),
	}
}

func (ms *mapStater) RegState(dialogType int) {
	ms.states[dialogType] = make(map[int]models.State)
}

func (ms *mapStater) AddState(dialogType int, chatID int) error {
	ms.states[dialogType][chatID] = models.State{
		Val:  1,
		Info: make(map[string]string),
	}
	return nil
}

func (ms *mapStater) GetState(dialogType int, chatID int) (int, error) {
	if state, isExist := ms.states[dialogType][chatID]; isExist {
		return state.Val, nil
	}
	return 0, nil
}

func (ms *mapStater) GetInfo(dialogType int, chatID int) (map[string]string, error) {
	if state, isExist := ms.states[dialogType][chatID]; isExist {
		return state.Info, nil
	}
	return nil, fmt.Errorf("No dialogs (type: %d) with user (id: %d)", dialogType, chatID)
}

func (ms *mapStater) NextState(dialogType int, chatID int, NewInfo map[string]string) error {
	if state, isExist := ms.states[dialogType][chatID]; isExist {
		state.Val += 1
		for infoName, info := range NewInfo {
			state.Info[infoName] = info
		}
		ms.states[dialogType][chatID] = state
	} else {
		return fmt.Errorf("No dialog (type: %d) with user (id: %d)", dialogType, chatID)
	}
	return nil
}

func (ms *mapStater) DelState(dialogType int, chatID int) error {
	if _, isExist := ms.states[dialogType][chatID]; !isExist {
		return fmt.Errorf("No dialog (type: %d) with user (id: %d)", dialogType, chatID)
	}
	delete(ms.states[dialogType], chatID)
	return nil
}
