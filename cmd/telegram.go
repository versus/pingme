package cmd

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nikoksr/notify"
	"github.com/nikoksr/notify/service/telegram"
	"github.com/urfave/cli/v2"
)

// teleGram struct holds data parsed via flags for telegram service.
type teleGram struct {
	Token   string
	Message string
	Channel string
	Title   string
}

// SendToTelegram parse values from *cli.context and return *cli.Command.
// Values include telegram token, channelIDs, Message and Title.
// If multiple channels are provided they the string is split with "," separator and
// each channelID is added to receiver.
func SendToTelegram() *cli.Command {
	var telegramOpts teleGram
	return &cli.Command{
		Name:  "telegram",
		Usage: "Send message to telegram",
		Description: `Telegram uses bot token to authenticate & send messages to defined channels.
Multiple channel ids can be used separated by comma ','.
All configuration options are also available via environment variables.`,
		UsageText: "pingme telegram --token '123' --channel '-123456' --msg 'some message'",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Destination: &telegramOpts.Token,
				Name:        "token",
				Aliases:     []string{"t"},
				Required:    true,
				Usage:       "Token of telegram bot used for sending message.",
				EnvVars:     []string{"TELEGRAM_TOKEN"},
			},
			&cli.StringFlag{
				Destination: &telegramOpts.Channel,
				Name:        "channel",
				Aliases:     []string{"c"},
				Required:    true,
				Usage:       "Channel ids of telegram.",
				EnvVars:     []string{"TELEGRAM_CHANNELS"},
			},
			&cli.StringFlag{
				Destination: &telegramOpts.Message,
				Name:        "msg",
				Aliases:     []string{"m"},
				Usage:       "Message content.",
				EnvVars:     []string{"TELEGRAM_MESSAGE"},
			},
			&cli.StringFlag{
				Destination: &telegramOpts.Title,
				Name:        "title",
				Value:       TimeValue,
				Usage:       "Title of the message.",
				EnvVars:     []string{"TELEGRAM_TITLE"},
			},
		},
		Action: func(ctx *cli.Context) error {
			notifier := notify.New()

			telegramSvc, err := telegram.New(telegramOpts.Token)
			if err != nil {
				return err
			}
			chn := strings.Split(telegramOpts.Channel, ",")
			for _, v := range chn {
				if len(v) <= 0 {
					return fmt.Errorf(EmptyChannel)
				}
				k, err := strconv.Atoi(v)
				if err != nil {
					log.Println(err)
				}
				telegramSvc.AddReceivers(int64(k))
			}

			notifier.UseServices(telegramSvc)

			if err = notifier.Send(
				context.Background(),
				telegramOpts.Title,
				telegramOpts.Message,
			); err != nil {
				return err
			}
			log.Println("Successfully sent!")
			return nil
		},
	}
}
