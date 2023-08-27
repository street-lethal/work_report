package usecase

import (
	"work_report/model"
	"work_report/service"
)

type MainUseCase interface {
	GenerateReport(inputFilePath, outputFilePath string) error
	SendReport(reportFilePath string) error
	ClearReport(outputFilePath string) error
}

type mainUseCase struct {
	service.GenerateReportDataService
	service.IOService
	service.ParseHTMLService
	service.ParseJiraService
	service.SendReportService
	model.Setting
}

func NewMainUseCase(
	gs service.GenerateReportDataService,
	is service.IOService,
	hs service.ParseHTMLService,
	js service.ParseJiraService,
	rs service.SendReportService,
	setting model.Setting,
) MainUseCase {
	return &mainUseCase{
		gs, is, hs, js, rs, setting,
	}
}

func (u mainUseCase) GenerateReport(inputFilePath, outputFilePath string) error {
	htm, err := u.IOService.Input(inputFilePath)
	if err != nil {
		return err
	}

	node, err := u.ParseHTMLService.Parse(htm)
	if err != nil {
		return err
	}

	works := u.ParseJiraService.Parse(node)
	report := u.Generate(works)
	return report.ToFile(outputFilePath)
}

func (u mainUseCase) SendReport(reportFilePath string) error {
	report, err := model.FileToReport(reportFilePath)
	if err != nil {
		return err
	}

	if err := u.SendReportService.Send(*report); err != nil {
		return err
	}

	return nil
}

func (u mainUseCase) ClearReport(outputFilePath string) error {
	report := u.Generate(nil)
	return report.ToFile(outputFilePath)
}
