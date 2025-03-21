package tools

import (
	"EuroprotocolTGBot/internal/loggin"
	"encoding/json"
	"os"
	"strconv"

	"go.uber.org/zap"
)

// структура текста с номером для хранения вопросов и ответов в бот
type TextWithID struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

// структура цепочки вопросов и ответов
type MsgChain struct {
	AnswerList map[int]TextWithID
	AskList    map[int]TextWithID
	CurrID     int
	Start      bool
}

// функция создания новой цепочки
func NewMsgChain() MsgChain {
	return MsgChain{
		AnswerList: make(map[int]TextWithID),
		AskList:    make(map[int]TextWithID),
		CurrID:     1,
		Start:      false,
	}
}

// функуция загрузки вопросов из файла конфигурации бота
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

	return nil
}

// функция получения текущего вопроса при сквозном опросе
func (chain *MsgChain) GetCurrentAsk() (TextWithID, bool) {
	chain.Start = true
	v, ok := chain.AskList[chain.CurrID]
	loggin.Log.Info("AskList LEN", zap.String("", strconv.Itoa(len(chain.AskList))))
	return v, ok
}

// функция сохранения ответа на текущий вопрос
func (chain *MsgChain) SetCurrentAnswer(answer string) {
	currAns := TextWithID{
		ID:   chain.CurrID,
		Text: answer,
	}
	chain.AnswerList[chain.CurrID] = currAns
	chain.CurrID++
}

// функция сброса ответов и текущей цепочки
func (chain *MsgChain) Reset() {
	chain.CurrID = 1
	chain.Start = false
	chain.AnswerList = make(map[int]TextWithID)
}

func (chain *MsgChain) PrintAnswer(filePath string) error {
	// Создаем файл для записи
	file, err := os.Create(filePath)
	if err != nil {
		loggin.Log.Error("PrintAnswer:", zap.String("Ошибка создания файла", err.Error()))
		return err
	}
	defer file.Close()

	// Конвертируем map в JSON
	jsonData, err := json.MarshalIndent(chain.AnswerList, "", "    ")
	if err != nil {
		loggin.Log.Error("PrintAnswer:", zap.String("Ошибка сериализации данных ответов", err.Error()))
		return err
	}

	// Записываем JSON в файл
	_, err = file.Write(jsonData)
	if err != nil {
		loggin.Log.Error("PrintAnswer:", zap.String("Ошибка записи в файл", err.Error()))
		return err
	}

	// Добавляем перевод строки в конец файла
	_, err = file.WriteString("\n")
	if err != nil {
		loggin.Log.Error("PrintAnswer:", zap.String("Ошибка записи в файл", err.Error()))
		return err
	}
	return nil
}
