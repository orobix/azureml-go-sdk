package workspaceiface

import "github.com/Telemaco019/azureml-go-sdk/workspace"

type WorkspaceAPI interface {
	GetDatastore() *workspace.Datastore
	GetDatastores() []workspace.Datastore
}
