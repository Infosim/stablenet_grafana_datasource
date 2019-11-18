package query

import (
	"github.com/grafana/grafana-plugin-model/go/datasource"
	testify "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildErrorResult(t *testing.T) {
	limerick := `Ein Limerickdichter aus Aachen,
nicht ahnend, was Limericks versprachen,
der trieb es zu bunt,
und das war der Grund,
dass Freunde zuletzt mit ihm brachen.`
	result := BuildErrorResult(limerick, "XYZ")
	assert := testify.New(t)
	assert.Equal(result.Error, limerick, "result error message wrong")
	assert.Equal(result.RefId, "XYZ", "result refId wrong")
	assert.Nil(result.Series, "series should be nil")
	assert.Empty(result.MetaJson, "meta json should be empty")
	assert.Nil(result.Tables, "tables should be nil")
}

func TestQuery_GetCustomField(t *testing.T) {
	rawquery := datasource.Query{
		ModelJson: "{\"favouriteDish\": \"all that is tasty\"}",
	}
	query := Query{rawquery}
	t.Run("test successful", func(t *testing.T) {
		actual, err := query.GetCustomField("favouriteDish")
		require.NoError(t, err)
		testify.Equal(t, "all that is tasty", actual)
	})
	t.Run("test error", func(t *testing.T) {
		_, err := query.GetCustomField("favouriteMeal")
		testify.EqualError(t, err, "type assertion to string failed")
	})
}

func TestQuery_GetCustomFieldNoJson(t *testing.T) {
	rawquery := datasource.Query{
		ModelJson: "{\"favouriteDish\": \"all that is tasty\"",
	}
	query := Query{rawquery}
	t.Run("test successful", func(t *testing.T) {
		_, err := query.GetCustomField("favouriteDish")
		require.EqualError(t, err, "unexpected EOF")
	})
}

func TestQuery_GetCustomIntField(t *testing.T) {
	rawquery := datasource.Query{
		ModelJson: "{\"age\": 5}",
	}
	query := Query{rawquery}
	t.Run("test successful", func(t *testing.T) {
		actual, err := query.GetCustomIntField("age")
		require.NoError(t, err)
		testify.Equal(t, 5, actual)
	})
	t.Run("test error", func(t *testing.T) {
		_, err := query.GetCustomIntField("birthYear")
		testify.EqualError(t, err, "invalid value type")
	})
}

func TestQuery_GetCustomIntFieldNoJson(t *testing.T) {
	rawquery := datasource.Query{
		ModelJson: "{\"favouriteDish\": \"all that is tasty\"",
	}
	query := Query{rawquery}
	t.Run("test successful", func(t *testing.T) {
		_, err := query.GetCustomIntField("favouriteDish")
		require.EqualError(t, err, "unexpected EOF")
	})
}

func Test_queryHandlerImpl_HandleWithNoQueryType(t *testing.T) {
	handler := queryHandlerImpl{}
	result := handler.Handle(Query{datasource.Query{}})
	assert := testify.New(t)
	assert.Equal("could not retrieve query type: EOF", result.Error)
	assert.Nil(result.Series)
}
