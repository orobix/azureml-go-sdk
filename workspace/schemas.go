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
	AuthorityUrl    string                       `json:"authorityUrl,omitempty"`
	SqlUserName     string                       `json:"userId,omitempty"`
}

type WriteDatastoreSchema struct {
	ContentsType         string                           `json:"contentsType"`
	StorageAccountName   string                           `json:"accountName,omitempty"`
	StorageContainerName string                           `json:"containerName,omitempty"`
	Credentials          *WriteDatastoreCredentialsSchema `json:"credentials"`
}

type WriteDatastoreSchemaProperties struct {
	Contents    WriteDatastoreSchema `json:"contents"`
	IsDefault   bool                 `json:"isDefault"`
	Description string               `json:"description"`
}

type SchemaWrapper struct {
	Properties interface{} `json:"properties"`
}
