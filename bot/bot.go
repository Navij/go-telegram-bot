package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const (
	apiURL string = "https://api.telegram.org/bot"
	apiKey string = "6294882859:AAFvgkmPejtGZAP1FZ_cH8w5Z54ddo1ardc"
)

type APIMethod struct {
	Type     string
	Endpoint string
}

type APIResponce struct {
	RawJSON []byte
	Raw     map[string]interface{}
}

type Bot struct {
	lastUpdate int
}

func New() *Bot {
	out := &Bot{}

	return out
}

func (b *Bot) Query(method string, httpMethod string, data map[string]interface{}) (*APIResponce, error) {
	var resp *http.Response
	var body io.ReadCloser
	var err error

	endpoint := fmt.Sprintf("%s%s/%s", apiURL, apiKey, method)
	dataJSON, err := json.Marshal(data)
	dataReader := bytes.NewBuffer(dataJSON)
	if err != nil {
		return nil, fmt.Errorf("Cannot marshal payload to JSON")
	}
	if httpMethod == "GET" {
		endpoint += fmt.Sprintf("?")
		for key, val := range data {
			endpoint += fmt.Sprintf("&%s=%s", key, val)
		}
		resp, err = http.Get(endpoint)
	} else if httpMethod == "POST" {
		resp, err = http.Post(endpoint, "application/json", dataReader)
	} else {
		return nil, fmt.Errorf("Unknown method type")
	}
	if err != nil {
		return nil, err
	}
	body = resp.Body
	bodyDataRaw, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("Cannot read answer: %s", err)
	}
	fmt.Println(string(bodyDataRaw))
	return &APIResponce{
		RawJSON: bodyDataRaw,
	}, nil
}

func getValue[T any](raw interface{}) *T {
	if raw == nil {
		return nil
	}
	val, ok := raw.(T)
	if !ok {
		return nil
	}
	return &val
}

func (b *Bot) GetMe() (*GetMeAnswer, error) {
	resp, err := b.Query("getMe", "GET", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	answer := &GetMeAnswer{}
	err = json.Unmarshal(resp.RawJSON, answer)

	if err != nil {
		return nil, fmt.Errorf("Cannot parse JSON: %s", resp.RawJSON)
	}

	if !answer.Ok {
		return nil, fmt.Errorf("Telegram API said bad request: %s", resp.RawJSON)
	}

	return answer, nil
}

type GetMeAnswer struct {
	Ok     bool `json:"ok"`
	Result GetMeAnswerResult
}

var getMeAnswerSchema map[string]interface{}

func init() {
	getMeAnswerSchema = map[string]interface{}{}
}

type GetMeAnswerResult struct {
	Id                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	Username                string `json:"username"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupprotInlineQueries    bool   `json:"supports_inline_queries"`
}

func (b *Bot) GetUpdates() (*GetUpdatesAnswer, error) {
	data := map[string]interface{}{
		"offset": strconv.Itoa(int(b.lastUpdate)),
	}
	resp, err := b.Query("getUpdates", "GET", data)
	if err != nil {
		return nil, err
	}

	answer := &GetUpdatesAnswer{}
	err = json.Unmarshal(resp.RawJSON, answer)

	if err != nil {
		return nil, fmt.Errorf("Cannot parse JSON: %s", resp.RawJSON)
	}

	if !answer.Ok {
		return nil, fmt.Errorf("Telegram API said bad request: %s", resp.RawJSON)
	}

	if len(answer.Result) > 0 {
		b.lastUpdate = answer.Result[len(answer.Result)-1].UpdateId + 1
	}

	return answer, nil
}

type GetUpdatesAnswer struct {
	Ok     bool
	Result []TelegramUpdate
}

type TelegramUpdate struct {
	UpdateId int             `json:"update_id"`
	Message  TelegramMessage `json:"message"`
}

type TelegramMessage struct {
	MessageId int          `json:"message_id"`
	Text      string       `json:"text"`
	Chat      TelegramChat `json:"chat"`
}

type TelegramChat struct {
	Id int `json:"id"`
}

func (b *Bot) SendMessage(data map[string]interface{}) (*SendMessageAnswer, error) {
	resp, err := b.Query("sendMessage", "POST", data)
	if err != nil {
		return nil, err
	}

	answer := &SendMessageAnswer{}
	err = json.Unmarshal(resp.RawJSON, answer)

	return answer, err
}

type SendMessageAnswer struct {
	Ok bool
}
