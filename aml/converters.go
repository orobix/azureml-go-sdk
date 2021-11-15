package aml

import (
	"fmt"
	"github.com/tidwall/gjson"
)

func toDatastore(json []byte) *Datastore {
	return &Datastore{
		Id:                   gjson.GetBytes(json, "id").Str,
		Name:                 gjson.GetBytes(json, "name").Str,
		Description:          gjson.GetBytes(json, "properties.description").Str,
		IsDefault:            gjson.GetBytes(json, "properties.isDefault").Bool(),
		StorageAccountName:   gjson.GetBytes(json, "properties.contents.accountName").Str,
		StorageContainerName: gjson.GetBytes(json, "properties.contents.containerName").Str,
		StorageContainerType: gjson.GetBytes(json, "properties.contents.contentsType").Str,
		CreationDate:         gjson.GetBytes(json, "systemData.createdAt").Time(),
		LastModifiedDate:     gjson.GetBytes(json, "systemData.lastModifiedAt").Time(),
	}
}

func toDatastoreArray(json []byte) []Datastore {
	jsonDatastoreArray := gjson.GetBytes(json, "value").Array()
	datastoreSlice := make([]Datastore, gjson.GetBytes(json, "value.#").Int())
	for i, jsonDatastore := range jsonDatastoreArray {
		datastore := toDatastore([]byte(jsonDatastore.Raw))
		datastoreSlice[i] = *datastore
		datastore.Id = "pippo"
		fmt.Println(datastore)
	}
	return datastoreSlice
}
