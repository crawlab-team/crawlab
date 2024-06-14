package interfaces

type ExportService interface {
	GenerateId() (exportId string, err error)
	Export(exportType, target string, filter Filter) (exportId string, err error)
	GetExport(exportId string) (export Export, err error)
}
