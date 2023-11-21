package ydb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/ydb-platform/ydb-go-sdk/v3/table"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/result/named"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

func (ydb *yandexDatabase) AddState(dialogType int, chatID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $ID AS String;
	DECLARE $ChatID AS Uint64;
	DECLARE $DialogType AS int64;
	DECLARE $State AS int64;
	DECLARE $Info AS JSON;
	UPSERT INTO States ( id, chat_id, dialog_type, state, info )
	VALUES ( $ID, $ChatID, $DialogType, $State, $Info );`

	params := table.NewQueryParameters(
		table.ValueParam("ID", types.BytesValueFromString(uuid.New().String())),
		table.ValueParam("ChatID", types.Uint64Value(uint64(chatID))),
		table.ValueParam("DialogType", types.Int64Value(int64(dialogType))),
		table.ValueParam("State", types.Int64Value(1)),
		table.ValueParam("Info", types.JSONValue("{}")),
	)

	return ydb.CUDoperations(ctx, query, params)
}

func (ydb *yandexDatabase) GetState(chatID int64) (int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $ChatID AS Uint64;
	SELECT dialog_type, state FROM States
	WHERE chat_id = $ChatID;`

	params := table.NewQueryParameters(
		table.ValueParam("ChatID", types.Uint64Value(uint64(chatID))),
	)

	var dialogType, state int64
	readTx := table.TxControl(table.BeginTx(table.WithOnlineReadOnly()), table.CommitTx())
	err := ydb.conn.Table().Do(ctx,
		func(ctx context.Context, s table.Session) (err error) {
			_, res, err := s.Execute(ctx, readTx, query, params)
			if err != nil {
				return fmt.Errorf("execute error: %v", err)
			}
			defer res.Close()

			if res.NextResultSet(ctx) {
				if res.NextRow() {
					err := res.ScanNamed(
						named.Required("dialog_type", &dialogType),
						named.Required("state", &state),
					)
					if err != nil {
						return fmt.Errorf("scan res error: %v", err)
					}
				}
			}
			return res.Err()
		},
	)
	return int(dialogType), int(state), err
}

func (ydb *yandexDatabase) GetInfo(dialogType int, chatID int64) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $ChatID AS Uint64;
	DECLARE $DialogType AS int64;
	SELECT info FROM States
	WHERE chat_id = $ChatID AND dialog_type = $DialogType;`

	params := table.NewQueryParameters(
		table.ValueParam("ChatID", types.Uint64Value(uint64(chatID))),
		table.ValueParam("DialogType", types.Int64Value(int64(dialogType))),
	)

	info, err := ydb.selectSingleValue(ctx, query, params)
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, fmt.Errorf("No dialogs (type: %d) with user (id: %d)", dialogType, chatID)
	}

	var infoMap map[string]string
	json.Unmarshal(info.([]byte), &infoMap)
	return infoMap, nil
}

func (ydb *yandexDatabase) NextState(dialogType int, chatID int64, NewInfo map[string]string) error {
	info, err := ydb.GetInfo(dialogType, chatID)
	if err != nil {
		return err
	}
	for name, val := range NewInfo {
		info[name] = val
	}

	infoJson, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal info: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $ChatID AS Uint64;
	DECLARE $DialogType AS int64;
	DECLARE $Info AS json;
	UPDATE States SET state = state + 1, info = $Info
	WHERE chat_id = $ChatID AND dialog_type = $DialogType;`

	params := table.NewQueryParameters(
		table.ValueParam("ChatID", types.Uint64Value(uint64(chatID))),
		table.ValueParam("DialogType", types.Int64Value(int64(dialogType))),
		table.ValueParam("Info", types.JSONValue(string(infoJson))),
	)
	return ydb.CUDoperations(ctx, query, params)
}

func (ydb *yandexDatabase) DelState(dialogType int, chatID int64) error {
	ctx, cancel := context.WithTimeout(context.Background(), operationDeadline)
	defer cancel()

	query := `DECLARE $ChatID AS Uint64;
	DECLARE $DialogType AS int64;
	DELETE FROM States
	WHERE chat_id = $ChatID AND dialog_type = $DialogType;`

	params := table.NewQueryParameters(
		table.ValueParam("ChatID", types.Uint64Value(uint64(chatID))),
		table.ValueParam("DialogType", types.Int64Value(int64(dialogType))),
	)

	return ydb.CUDoperations(ctx, query, params)
}
