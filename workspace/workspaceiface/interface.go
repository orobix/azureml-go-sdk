package workspaceiface

import "github.com/telemaco019/azureml-workspace-go-sdk/workspace"

type WorkspaceAPI interface {
	GetDatastore() *workspace.Datastore
	GetDatastores() []workspace.Datastore
}
