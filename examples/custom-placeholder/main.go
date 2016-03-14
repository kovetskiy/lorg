package main

import "github.com/kovetskiy/lorg"

func main() {
	format := lorg.NewFormat(`[${level}] ${file}:${line} ${custom} %s`)
	format.SetPlaceholder("custom", func(_ lorg.Level, _ string) string {
		return "~custom~"
	})

	log := lorg.NewLog()
	log.SetFormat(format)

	log.Print("Print: info message")
	log.Printf("Printf: %s %s", "info", "message")

	log.Info("Info: info message")
	log.Infof("Infof: %s %s", "info", "message")

	log.Warning("Warning: warning message")
	log.Warningf("Warningf: %s %s", "warning", "message")

	log.Error("Error: error message")
	log.Errorf("Errorf: %s %s", "error", "message")

	// these messages will not be printed, because default level is LevelInfo
	log.Debug("Debug: debug message")
	log.Debugf("Debugf: %s %s", "debug", "message")

	log.SetLevel(lorg.LevelDebug)

	log.Debug("Debug: debug message")
	log.Debugf("Debugf: %s %s", "debug", "message")
}
