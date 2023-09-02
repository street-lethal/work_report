package main

import (
	"flag"
	"work_report/registry"
	"work_report/usecase"
)

var (
	mode                    = flag.String("mode", "gen", `"gen", "send", "send-m" or "clear" (default: "gen")`)
	inputFilePath           = "./data/jira.html"
	outputFilePath          = "./data/report.json"
	settingsFilePath        = "./config/settings.json"
	platformSessionFilePath = "./data/platform_session.json"
	platformIDFilePath      = "./config/platform_id.json"
)

func main() {
	flag.Parse()

	switch *mode {
	case "send":
		send()
	case "send-m":
		sendManually()
	case "clear":
		clear()
	default:
		gen()
	}
}

func mainUseCase() usecase.MainUseCase {
	r, err := registry.NewRegistry(settingsFilePath)
	if err != nil {
		panic(err)
	}

	return r.NewMainUseCase()
}

func send() {
	err := mainUseCase().LogInAndSendReport(platformIDFilePath, outputFilePath)
	if err != nil {
		panic(err)
	}
}

func sendManually() {
	err := mainUseCase().SendReport(outputFilePath, platformSessionFilePath)
	if err != nil {
		panic(err)
	}
}

func gen() {
	if err := mainUseCase().GenerateReport(inputFilePath, outputFilePath); err != nil {
		panic(err)
	}
}

func clear() {
	if err := mainUseCase().ClearReport(outputFilePath); err != nil {
		panic(err)
	}
}
