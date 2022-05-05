package main

import (
	"context"
	"github.com/hectorgimenez/koolo/api"
	zapLogger "github.com/hectorgimenez/koolo/cmd/koolo/log"
	koolo "github.com/hectorgimenez/koolo/internal"
	"github.com/hectorgimenez/koolo/internal/action"
	"github.com/hectorgimenez/koolo/internal/character"
	"github.com/hectorgimenez/koolo/internal/config"
	"github.com/hectorgimenez/koolo/internal/health"
	"github.com/hectorgimenez/koolo/internal/remote"
	"github.com/hectorgimenez/koolo/internal/run"
	"github.com/hectorgimenez/koolo/internal/stats"
	"github.com/hectorgimenez/koolo/internal/town"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err.Error())
	}

	logger, err := zapLogger.NewLogger(config.Config.Debug, config.Config.LogFilePath)
	if err != nil {
		log.Fatalf("Error starting logger: %s", err.Error())
	}
	defer logger.Sync()

	ctx, cancel := context.WithCancel(context.Background())
	ma := api.NewMapAssistClient(logger)
	controller := remote.New(config.Config.Controller.Port)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-ch
		logger.Info("Shutting down...")
		signal.Stop(ch)
		cancel()
		ma.Stop()
		if config.Config.Controller.Webserver {
			controller.Stop(ctx)
		}
	}()
	g, ctx := errgroup.WithContext(ctx)

	err = ma.StartAndConfigure(ctx)
	if err != nil {
		logger.Fatal("error starting MapAssist", zap.Error(err))
	}

	bm := health.NewBeltManager(logger)
	hm := health.NewHealthManager(logger, bm)
	sm := town.NewShopManager(logger, bm)
	char, err := character.BuildCharacter()
	if err != nil {
		logger.Fatal("Error creating character", zap.Error(err))
	}

	ab := action.NewBuilder(logger, sm, bm)
	bot := koolo.NewBot(logger, hm, ab)
	supervisor := koolo.NewSupervisor(logger, bot)

	g.Go(func() error {
		return supervisor.Start(ctx, run.BuildRuns(ab, char))
	})

	if config.Config.Controller.Webserver {
		g.Go(func() error {
			return controller.Run()
		})
	}

	discordBot, err := stats.NewDiscordBot(config.Config.Discord.Token, config.Config.Discord.ChannelID)
	if err != nil {
		logger.Fatal("Discord could not been initialized", zap.Error(err))
	}
	eventListener := stats.NewEventListener(discordBot, logger)
	if config.Config.Discord.Enabled {
		g.Go(func() error {
			return discordBot.Start(ctx)
		})
	}

	g.Go(func() error {
		return eventListener.Listen(ctx)
	})

	err = g.Wait()
	if err != nil {
		log.Fatalf("Error running Koolo: %s", err.Error())
	}
}
