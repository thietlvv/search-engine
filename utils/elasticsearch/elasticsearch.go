package elasticsearch

import (
	"auth/config"
	"context"
	"os"
	"time"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

const mapping = `
{
   "settings":{
      "number_of_shards":3,
      "number_of_replicas":2
   },
   "mappings":{
      "properties":{
         "birth_day":{
            "type":"text"
         },
         "cmnd_issue_date":{
            "type":"text"
         }
      }
   }
}`

func InitElasticsearch(index_name string) (*elastic.Client, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "localhost"
	}
	index_name_env := env + "_" + index_name
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var options []elastic.ClientOptionFunc
	if config.C.Elasticsearch.IsAuth {
		options = append(options, elastic.SetSniff(false), elastic.SetURL(config.C.Elasticsearch.URL), elastic.SetBasicAuth(config.C.Elasticsearch.User, config.C.Elasticsearch.Secret))
	} else {
		options = append(options, elastic.SetSniff(false), elastic.SetURL(config.C.Elasticsearch.URL))
	}
	client, err := elastic.NewClient(options...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists(index_name_env).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex(index_name_env).BodyString(mapping).Do(ctx)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	return client, nil
}

func Put(client *elastic.Client, index_name string, type_name string, data interface{}, id string) (string, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "localhost"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	index_name_env := env + "_" + index_name
	res, err := client.Index().
		Index(index_name_env).
		Id(id).
		BodyJson(data).
		Do(ctx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	_, err = client.Flush().Index(index_name_env).Do(ctx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}
	return res.Result, nil
}

func Update(client *elastic.Client, index_name string, type_name string, data map[string]interface{}, id string) (string, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "localhost"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	index_name_env := env + "_" + index_name
	res, err := client.Update().
		Index(index_name_env).
		Id(id).
		Doc(data).
		Do(ctx)
	if err != nil {
		logrus.Error(err)
		return "", err
	}

	return res.Id, nil
}
