package commandst

type CommandsSender interface {
	SendCommand(cmd CommandMessage, queue string) error
}

type CommandsMetrics interface {
	CommandUsed(command string)
	CommandFailed(command string)
}

type CommandsReceiver interface {
	ReceiveCommands(command string) (<-chan CommandMessage, error)
	AddHandler(command string, handler ReceiverHandler)
	SetErrorMetrics(metrics CommandsMetrics)
	Start(command string) error
	Close()
}

type CommandMessage struct {
	GuildID   string `json:"guild_id"`
	ChannelID string `json:"channel_id"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
}

type ReceiverHandler func(m CommandMessage) error

type PingStatus int

const (
	RCommandPing = PingStatus(0)
	RCommandPong = PingStatus(1)
)

type PingMessage struct {
	Status PingStatus
}
