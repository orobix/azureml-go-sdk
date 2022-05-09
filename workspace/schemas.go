package workspace

type WriteDatastoreSecretsSchema struct {
	SecretsType     string `json:"secretsType"`
	AccountKey      string `json:"key,omitempty"`
	ClientSecret    string `json:"clientSecret,omitempty"`
	SqlUserPassword string `json:"password,omitempty"`
}

type WriteDatastoreCredentialsSchema struct {
	CredentialsType string                       `json:"credentialsType"`
	Secrets         *WriteDatastoreSecretsSchema `json:"secrets"`
	ClientId        string                       `json:"clientId,omitempty"`
	TenantId        string                       `json:"tenantId,omitempty"`
	SqlUserName     string                       `json:"userId,omitempty"`
}

type WriteDatastoreSchema struct {
	ContentsType         string                           `json:"contentsType"`
	StorageAccountName   string                           `json:"accountName,omitempty"`
	StorageContainerName string                           `json:"containerName,omitempty"`
	Credentials          *WriteDatastoreCredentialsSchema `json:"credentials,omitempty"`
	Endpoint             string                           `json:"endpoint"`
	Protocol             string                           `json:"protocol"`
}

type WriteDatastoreSchemaProperties struct {
	Contents    WriteDatastoreSchema `json:"contents"`
	IsDefault   bool                 `json:"isDefault"`
	Description string               `json:"description"`
}

type DatasetPathsSchema struct {
	FilePath      string `json:"file,omitempty"`
	DirectoryPath string `json:"folder,omitempty"`
}
type WriteDatasetSchema struct {
	Description string               `json:"description,omitempty"`
	Paths       []DatasetPathsSchema `json:"paths"`
}

type SchemaWrapper struct {
	Properties interface{} `json:"properties"`
}
