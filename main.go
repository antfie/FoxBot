package main

import (
	_ "embed"
	"fmt"
	"foxbot/config"
	"foxbot/db"
	"foxbot/tasks"
	"foxbot/utils"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

//go:embed config.yaml
var defaultConfigData []byte

func main() {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			stackSize := runtime.Stack(buf, false)
			stackTrace := string(buf[:stackSize])
			log.Printf("Main panic recovered: %v\nStack trace:\n%s", r, stackTrace)
		}
	}()

	c := config.Load(defaultConfigData)
	task := &tasks.Context{
		Config: c,
		DB:     db.NewDB(c.DBPath),
	}

	var tasksToRun []*tasks.Task

	if c.Reminders != nil {
		if len(c.Reminders.Reminders) < 1 {
			log.Print("No reminders configured.")
		} else {
			tasksToRun = append(tasksToRun, tasks.NewTask(c.Reminders.Check.Frequency, task.Reminders))
		}
	}

	if c.Countdown != nil {
		if len(c.Countdown.Timers) < 1 {
			log.Print("No countdown timers configured.")
		} else {
			tasksToRun = append(tasksToRun, tasks.NewTask(c.Countdown.Check.Frequency, task.Countdown))
		}
	}

	if c.RSS != nil {
		if len(c.RSS.Feeds) < 1 {
			log.Print("No RSS feeds configured.")
		} else {
			tasksToRun = append(tasksToRun, tasks.NewTask(c.RSS.Check.Frequency, task.RSS))
		}
	}

	if c.SiteChanges != nil {
		if len(c.SiteChanges.Sites) < 1 {
			log.Print("No sites to monitor configured.")
		} else {
			tasksToRun = append(tasksToRun, tasks.NewTask(c.SiteChanges.Check.Frequency, task.SiteChanges))
		}
	}

	if len(tasksToRun) == 0 {
		log.Print("Error: No tasks to run.")
		os.Exit(1)
	}

	task.Notify(fmt.Sprintf("🦊🤖 Running with %s.", utils.Pluralize("task", len(tasksToRun))))
	go tasks.Run(tasksToRun)

	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-shutdownSignal

	// Clear out any "^C" from the console
	print("\r")

	task.Notify("🦊🤖 Stopped")
}
