package types

type WorkflowFile struct {
	Workflows map[string]Workflow `json:"workflows"`
}

type WebhookWorkflow struct {
	Action      string             `json:"action"`
	WorkflowRun WorkflowRun        `json:"workflow_run"`
	Repository  WorkflowRepository `json:"repository"`
}

type WorkflowRun struct {
	Name       string `json:"name"`
	HeadBranch string `json:"head_branch"`
	Event      string `json:"event"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
	UpdatedAt  string `json:"updated_at"`
	Url        string `json:"html_url"`
}

type WorkflowRepository struct {
	Name          string `json:"name"`
	Url           string `json:"html_url"`
	DefaultBranch string `json:"default_branch"`
}

type WorkflowData struct {
	Workflow    Workflow
	HasWorkflow bool
}

type Workflow struct {
	Status     string `json:"Status"`
	Conclusion string `json:"Conclusion"`
	UpdatedAt  string `json:"UpdatedAt"`
	Url        string `json:"Url"`
	Name       string `json:"Name"`
	Branch     string `json:"Branch"`
}
