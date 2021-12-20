# Azure ML Go SDK

[![Actions Status](https://github.com/telemaco019/azureml-go-sdk/workflows/test/badge.svg)](https://github.com/telemaco019/azureml-go-sdk/actions)
[![codecov](https://codecov.io/gh/telemaco019/azureml-go-sdk/branch/main/graph/badge.svg)](https://codecov.io/gh/telemaco019/azureml-go-sdk)

Go SDK for configuring [Azure Machine Learning](https://azure.microsoft.com/en-us/services/machine-learning/)
workspaces.

**The library is still under development and at the moment it only supports CRUD operations over Datastores of AML
Workspaces.**

## Getting Started

### Installation

Use go get to retrieve the SDK to add it to your GOPATH workspace, or project's Go module dependencies.

```shell
go get github.com/Telemaco019/azureml-go-sdk
```

To update the SDK use go get -u to retrieve the latest version of the SDK.

```shell
go get -u github.com/Telemaco019/azureml-go-sdk
```

## Quick Examples

### Init the client

```go
import (
  "github.com/Telemaco019/azureml-go-sdk/workspace"
)

config := workspace.Config{
  ClientId:       "", // the client ID of the Service Principal used for authenticating with Azure
  ClientSecret:   "", // the client secret of the Service Principal used for authenticating with Azure
  TenantId:       "", // the tenant ID to which the Service Principal used for authenticating with Azure belongs to
  SubscriptionId: "", // the Azure Subscription ID of the subscription containing the AML Workspace
}

ws, err := workspace.New(config, true)
```

### Get all the Datastores of a workspace

```go
datastores, err := ws.GetDatastores( "rg-name", "workspace-name" )
```

### Get a specific Datastore of a workspace

```go
datastore, err := ws.GetDatastores( "rg-name", "workspace-name", "datastore-name" )
```
