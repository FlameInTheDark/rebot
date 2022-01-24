package service

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"

	"github.com/FlameInTheDark/rebot/business/transport/commandst"
)

type ExchangeService struct {
	cr     commandst.CommandsReceiver
	sess   *discordgo.Session
	logger *zap.Logger
}
