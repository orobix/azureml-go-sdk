package workspace

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func unmarshalDatastore(json []byte) *Datastore {
	sysData := SystemData{
		CreationDate:         gjson.GetBytes(json, "systemData.createdAt").Time(),
		CreationUserType:     gjson.GetBytes(json, "systemData.createdByType").Str,
		CreationUser:         gjson.GetBytes(json, "systemData.createdBy").Str,
		LastModifiedDate:     gjson.GetBytes(json, "systemData.lastModifiedAt").Time(),
		LastModifiedUserType: gjson.GetBytes(json, "systemData.lastModifiedByType").Str,
		LastModifiedUser:     gjson.GetBytes(json, "systemData.lastModifiedBy").Str,
	}
	auth := DatastoreAuth{
		CredentialsType: gjson.GetBytes(json, "properties.contents.credentials.credentialsType").Str,
		ClientId:        gjson.GetBytes(json, "properties.contents.credentials.clientId").Str,
		ClientSecret:    gjson.GetBytes(json, "properties.contents.credentials.secret.clientSecret").Str,
		AccountKey:      gjson.GetBytes(json, "properties.contents.credentials.secret.accountKey").Str,
		SqlUserName:     gjson.GetBytes(json, "properties.contents.credentials.secret.userId").Str,
		SqlUserPassword: gjson.GetBytes(json, "properties.contents.credentials.secret.password").Str,
	}
	return &Datastore{
		Id:                   gjson.GetBytes(json, "id").Str,
		Name:                 gjson.GetBytes(json, "name").Str,
		Type:                 gjson.GetBytes(json, "properties.contents.contentsType").Str,
		Description:          gjson.GetBytes(json, "properties.description").Str,
		IsDefault:            gjson.GetBytes(json, "properties.isDefault").Bool(),
		StorageAccountName:   gjson.GetBytes(json, "properties.contents.accountName").Str,
		StorageContainerName: gjson.GetBytes(json, "properties.contents.containerName").Str,
		StorageType:          gjson.GetBytes(json, "properties.contents.contentsType").Str,

		SystemData: &sysData,
		Auth:       &auth,
	}
}

func unmarshalDatastoreArray(json []byte) []Datastore {
	jsonDatastoreArray := gjson.GetBytes(json, "value").Array()
	datastoreSlice := make([]Datastore, gjson.GetBytes(json, "value.#").Int())
	for i, jsonDatastore := range jsonDatastoreArray {
		datastore := unmarshalDatastore([]byte(jsonDatastore.Raw))
		datastoreSlice[i] = *datastore
		fmt.Println(datastore)
	}
	return datastoreSlice
}
