package registry

import (
	"work_report/model"
	"work_report/service"
	"work_report/usecase"
)

type Registry struct {
	setting model.Setting
}

func NewRegistry(settingFilePath string) (*Registry, error) {
	setting, err := model.FileToSetting(settingFilePath)
	if err != nil {
		return nil, err
	}

	return &Registry{*setting}, nil
}

func (i *Registry) NewGenerateReportService() service.GenerateReportDataService {
	return service.NewGenerateReportDataService(i.setting)
}

func (i *Registry) NewIOService() service.IOService {
	return service.NewIOService(
		i.NewGenerateReportService(),
	)
}

func (i *Registry) NewParseHTMLService() service.ParseHTMLService {
	return service.NewParseHTMLService()
}

func (i *Registry) NewParseJiraService() service.ParseJiraService {
	return service.NewParseJiraService(
		i.NewParseHTMLService(),
	)
}

func (i *Registry) NewSendReportService() service.SendReportService {
	return service.NewSendReportService(i.setting)
}

func (i *Registry) NewMainUseCase() usecase.MainUseCase {
	return usecase.NewMainUseCase(
		i.NewGenerateReportService(),
		i.NewIOService(),
		i.NewParseHTMLService(),
		i.NewParseJiraService(),
		i.NewSendReportService(),
		i.setting,
	)
}
