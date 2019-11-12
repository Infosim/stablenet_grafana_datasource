package request

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"github.com/grafana/grafana-plugin-model/go/datasource"
)

type Content struct {
	datasource.DatasourceRequest
}

func (c *Content) GetCustomField(name string) (string, error) {
	if len(c.Queries) != 1 {
		return "", fmt.Errorf("can only extract custom field \"%s\" if number of queries is 1, but was %d", name, len(c.Queries))
	}
	firstQuery := c.Queries[0]
	queryJson, err := simplejson.NewJson([]byte(firstQuery.ModelJson))
	if err != nil {
		return "", err
	}
	return queryJson.Get(name).String()
}

func (c *Content) GetCustomIntField(name string) (int, error) {
	if len(c.Queries) != 1 {
		return 0, fmt.Errorf("can only extract custom field \"%s\" if number of queries is 1, but was %d", name, len(c.Queries))
	}
	firstQuery := c.Queries[0]
	queryJson, err := simplejson.NewJson([]byte(firstQuery.ModelJson))
	if err != nil {
		return 0, err
	}
	return queryJson.Get(name).Int()
}
