package datastore

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/datastore"
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"
)

func Put(db *datastore.Client, ctx context.Context, kind_name string, src interface{}) (*datastore.Key, error) {
	key := datastore.NameKey(kind_name, "", nil)
	if key, err := db.Put(ctx, key, src); err != nil {
		return nil, err
	} else {
		return key, nil
	}
}

func Get(db *datastore.Client, ctx context.Context, kind_name string, dst interface{}, filters []DatastoreFilter, limit int) error {
	query := datastore.NewQuery(kind_name)
	for _, i := range filters {
		query = query.Filter(i.FilterString, i.Value)
	}
	if _, err := db.GetAll(ctx, query.Limit(limit), dst); err != nil {
		return err
	}
	return nil
}

func GetByID(db *datastore.Client, ctx context.Context, dst interface{}, id string) error {
	key, err := datastore.DecodeKey(id)
	fmt.Println("\nid: ", id)
	fmt.Println("\nkey: ", key)
	if err != nil {
		return err
	}
	if err := db.Get(ctx, key, dst); err != nil {
		return err
	}
	return nil
}

func Update(db *datastore.Client, ctx context.Context, dst interface{}, id string) error {
	key, err := datastore.DecodeKey(id)
	if err != nil {
		return err
	}
	if _, err := db.Put(ctx, key, dst); err != nil {
		return err
	}
	return nil
}

func UpdateByMap(db *datastore.Client, ctx context.Context, old_interface interface{}, data map[string]interface{}, id string) error {
	err := GetByID(db, ctx, old_interface, id)
	if err != nil {
		return err
	}
	var inInterface map[string]interface{}
	inrec, _ := json.Marshal(old_interface)
	json.Unmarshal(inrec, &inInterface)
	if err := mergo.Merge(&inInterface, data, mergo.WithOverwriteWithEmptyValue); err != nil {
		return err
	}
	mapstructure.Decode(inInterface, old_interface)

	key, err := datastore.DecodeKey(id)
	if err != nil {
		return err
	}
	if _, err := db.Put(ctx, key, old_interface); err != nil {
		return err
	}
	return nil
}
