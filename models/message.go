package models

const (
	MESSAGE_COMMAND_OPEN      = "OPEN"
	MESSAGE_COMMAND_CLOSE     = "CLOSE"
	MESSAGE_COMMAND_NEW_STOCK = "NEW_STOCK"
	MESSAGE_COMMAND_BUY       = "BUY"
	MESSAGE_COMMAND_SELL      = "SELL"
	MESSAGE_COMMAND_ERROR     = "ERROR"
	MESSAGE_COMMAND_SUMMARY   = "SUMMARY"
)

type Message struct {
	Command   string     `json:"command"`
	Order     *Order     `json:"order"`
	Stock     *Stock     `json:"stock"`
	Summaries []*Summary `json:"summaries"`
	Error     string     `json:"error"`
}

func (msg *Message) GetCommand() string {
	return msg.Command
}

func NewErrorMessage(msg string) *Message {
	return &Message{
		Command: MESSAGE_COMMAND_ERROR,
		Error:   msg,
	}
}
