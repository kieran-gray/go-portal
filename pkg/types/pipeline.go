package types

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
