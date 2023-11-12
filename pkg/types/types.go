package types

type ServicesFile struct {
	Metadata Metadata  `json:"metadata"`
	Services []Service `json:"services"`
}

type Metadata struct {
	LastUpdatedOn string `json:"last_updated_on"`
	LastUpdatedBy string `json:"last_updated_by"`
	Version       string `json:"version"`
}

type Service struct {
	Metadata ServiceMetadata `json:"metadata"`
	Ui       ServiceDetails  `json:"ui"`
	Api      ServiceDetails  `json:"api"`
}

type ServiceMetadata struct {
	Name    string `json:"name"`
	Aliases string `json:"aliases"`
	DevOnly bool   `json:"dev_only"`
}

type ServiceDetails struct {
	Description   string        `json:"description"`
	RepositoryUrl string        `json:"repository_url"`
	Environments  []Environment `json:"environments"`
}

type Environment struct {
	Name     string `json:"name"`
	Priority int16  `json:"priority"`
	Url      string `json:"url"`
	LogsUrl  string `json:"logs_url"`
}

type Message struct {
	Status  string
	Message string
}

type PipelineFile struct {
	Pipelines map[string]Pipeline `json:"pipelines"`
}

type Pipeline struct {
	Status        string `json:"status"`
	UpdatedAt     string `json:"updated_at"`
	PipelineUrl   string `json:"pipeline_url"`
	ProjectId     int64  `json:"project_id"`
	DefaultBranch string `json:"default_branch"`
}

type PipelineData struct {
	Pipeline    Pipeline
	HasPipeline bool
}

type DisplayServices struct {
	Favourites []Service
	Services   []Service
}

type IndexData struct {
	Favourites []Service
	Services   []Service
	Pipelines  map[string]Pipeline
}

type AdminData struct {
	ServicesFile ServicesFile
	Messages     []Message
}

type ServiceCardData struct {
	Service       Service
	Id            string
	SearchAliases string
	Favourite     bool
	HasUi         bool
	HasApi        bool
	Pipelines     map[string]Pipeline
}

type ServiceCardTabData struct {
	ServiceDetails ServiceDetails
	Id             string
	Name           string
	ServiceType    string
	HasLogs        bool
	Favourite      bool
	Pipeline       PipelineData
	ShowTab        bool
}

type ServiceCardFooterButtonData struct {
	HasLogs         bool
	EnvironmentName string
	Url             string
	LogsUrl         string
}

type AdminCardData struct {
	Service       Service
	Id            string
	SearchAliases string
}

type AdminCardTabData struct {
	ServiceDetails ServiceDetails
	Name           string
	ServiceType    string
	Id             string
}
