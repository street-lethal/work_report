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
	return service.NewIOService()
}

func (i *Registry) NewParseHTMLService() service.ParseHTMLService {
	return service.NewParseHTMLService()
}

func (i *Registry) NewParseCSVService() service.ParseCSVService {
	return service.NewParseCSVService()
}

func (i *Registry) NewParseJiraHTMLService() service.ParseJiraHTMLService {
	return service.NewParseJiraHTMLService(
		i.NewParseHTMLService(),
	)
}

func (i *Registry) NewParseJiraCSVService() service.ParseJiraCSVService {
	return service.NewParseJiraCSVService()
}

func (i *Registry) NewSendReportService() service.SendReportService {
	return service.NewSendReportService(i.setting)
}

func (i *Registry) NewLoginPlatformService() service.LoginPlatformService {
	return service.NewLoginPlatformService()
}

func (i *Registry) NewFetchPlatformService() service.FetchPlatformWorkService {
	return service.NewFetchPlatformWorkService(
		i.setting,
		i.NewParseHTMLService(),
	)
}

func (i *Registry) NewMainUseCase() usecase.MainUseCase {
	return usecase.NewMainUseCase(
		i.NewGenerateReportService(),
		i.NewIOService(),
		i.NewParseHTMLService(),
		i.NewParseCSVService(),
		i.NewParseJiraCSVService(),
		i.NewSendReportService(),
		i.NewLoginPlatformService(),
		i.NewFetchPlatformService(),
		i.setting,
	)
}
