package utils

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestGetVariable(t *testing.T) {
	testCases := []struct {
		name                string
		httpMethod          string
		urlRouterFormat     string
		urlRequestFormat    string
		inputVariablesNames []interface{}
		expectedVariables   []interface{}
	}{
		{
			name:                "test_one_param",
			httpMethod:          "GET",
			urlRouterFormat:     "/{%s}",
			urlRequestFormat:    "/%s",
			inputVariablesNames: []interface{}{"table"},
			expectedVariables:   []interface{}{"animals"},
		},
		{
			name:                "test_two_params",
			httpMethod:          "GET",
			urlRouterFormat:     "/{%s}/{%s}",
			urlRequestFormat:    "/%s/%s",
			inputVariablesNames: []interface{}{"table", "id"},
			expectedVariables:   []interface{}{"animals", "1"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			router := mux.NewRouter()
			handler := func(resp http.ResponseWriter, req *http.Request) {
				variables := make([]interface{}, 0, len(testCase.inputVariablesNames))
				for _, inputVariableName := range testCase.inputVariablesNames {
					inputVariableNameStr, _ := inputVariableName.(string)
					variables = append(variables, GetVariable(inputVariableNameStr, req))
				}
				assert.Equal(t, variables, testCase.expectedVariables)
			}
			urlForRouter := fmt.Sprintf(testCase.urlRouterFormat, testCase.inputVariablesNames...)
			router.HandleFunc(urlForRouter, handler).Methods(testCase.httpMethod)
			urlForRequest := fmt.Sprintf(testCase.urlRequestFormat, testCase.expectedVariables...)
			request := httptest.NewRequest(testCase.httpMethod, urlForRequest, nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, request)
		})
	}
}
