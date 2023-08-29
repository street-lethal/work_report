package usecase

import (
	"work_report/model"
	"work_report/service"
)

type MainUseCase interface {
	GenerateReport(inputFilePath, outputFilePath string) error
	SendReport(reportFilePath, sessionFilePath string) error
	LogInAndSendReport(
		platformIDFilePath, reportFilePath string,
	) error
	ClearReport(outputFilePath string) error
}

type mainUseCase struct {
	service.GenerateReportDataService
	service.IOService
	service.ParseHTMLService
	service.ParseJiraService
	service.SendReportService
	service.LoginPlatformService
	model.Setting
}

func NewMainUseCase(
	gs service.GenerateReportDataService,
	is service.IOService,
	hs service.ParseHTMLService,
	js service.ParseJiraService,
	rs service.SendReportService,
	ls service.LoginPlatformService,
	setting model.Setting,
) MainUseCase {
	return &mainUseCase{
		gs, is, hs, js, rs, ls, setting,
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

func (u mainUseCase) SendReport(reportFilePath, sessionFilePath string) error {
	report, err := model.FileToReport(reportFilePath)
	if err != nil {
		return err
	}

	session, err := model.FileToPlatformSession(sessionFilePath)
	if err != nil {
		return err
	}

	if err := u.SendReportService.Send(*report, session); err != nil {
		return err
	}

	return nil
}

func (u mainUseCase) LogInAndSendReport(
	platformIDFilePath, reportFilePath string,
) error {
	id, err := model.FileToPlatformID(platformIDFilePath)
	if err != nil {
		return err
	}

	report, err := model.FileToReport(reportFilePath)
	if err != nil {
		return err
	}

	session, err := u.LoginPlatformService.Login(*id)
	if err != nil {
		return err
	}

	if err := u.SendReportService.Send(*report, session); err != nil {
		return err
	}

	return nil
}

func (u mainUseCase) ClearReport(outputFilePath string) error {
	report := u.Generate(nil)
	return report.ToFile(outputFilePath)
}
