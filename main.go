package main

import (
	"flag"
	"work_report/registry"
	"work_report/usecase"
)

var (
	mode           = flag.String("mode", "gen", `"gen", "send" or "clear" (default: "gen")`)
	inputFilePath  = "./data/jira.html"
	outputFilePath = "./data/report.json"
)

func main() {
	flag.Parse()

	switch *mode {
	case "send":
		send()
	case "clear":
		clear()
	default:
		gen()
	}
}

func mainUseCase() usecase.MainUseCase {
	r, err := registry.NewRegistry("./config/settings.json")
	if err != nil {
		panic(err)
	}

	return r.NewMainUseCase()
}

func send() {
	if err := mainUseCase().SendReport(outputFilePath); err != nil {
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
