package app

type Config struct {
	IntegrationID      string
	WorkflowInstanceID string
	InputDirectory     string
	OutputDirectory    string
	ConsentGroup       string
	AnalyteType        string
	IsTumor            bool
}
