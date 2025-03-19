package tools

import (
	"EuroprotocolTGBot/internal/loggin"
	"encoding/json"
	"os"
	"strconv"

	"go.uber.org/zap"
)

type TextWithID struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type MsgChain struct {
	AnswerList map[int]TextWithID
	AskList    map[int]TextWithID
	CurrID     int
	Start      bool
}

func NewMsgChain() *MsgChain {
	return &MsgChain{
		AnswerList: make(map[int]TextWithID),
		AskList:    make(map[int]TextWithID),
		CurrID:     1,
		Start:      false,
	}
}

func (chain *MsgChain) LoadAsks(file string) error {
	content, err := os.ReadFile(file)
	if err != nil {
		loggin.Log.Debug("", zap.String("Ошибка при чтении файла: %s", err.Error()))
		return err
	}

	// Создаем срез для хранения объектов
	var texts []TextWithID

	// Разбираем JSON в наш срез структур
	err = json.Unmarshal(content, &texts)
	if err != nil {
		loggin.Log.Debug("", zap.String("Ошибка при десериализации JSON: %s", err.Error()))
		return err
	}

	for _, v := range texts {
		chain.AskList[v.ID] = v
		loggin.Log.Info("AskList", zap.String("ID: ", strconv.Itoa(v.ID)+" Text:"+v.Text))
	}

	return nil
}

func (chain *MsgChain) GetCurrentAsk() (TextWithID, bool) {
	chain.Start = true
	v, ok := chain.AskList[chain.CurrID]
	loggin.Log.Info("AskList LEN", zap.String("", strconv.Itoa(len(chain.AskList))))
	return v, ok
}

func (chain *MsgChain) SetCurrentAnswer(answer string) {
	currAns := TextWithID{
		ID:   chain.CurrID,
		Text: answer,
	}
	chain.AnswerList[chain.CurrID] = currAns
	chain.CurrID++
}

func (chain *MsgChain) Reset() {
	chain.CurrID = 1
	chain.Start = false
	chain.AnswerList = make(map[int]TextWithID)
}
