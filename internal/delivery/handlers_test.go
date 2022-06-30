package delivery

import (
	"fmt"
	"github/mysql-dbmanager/internal/controller/mock_controller"
	"github/mysql-dbmanager/internal/model"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/magiconair/properties/assert"
)

func TestGetTablesHandler(t *testing.T) {
	type mockBehavior = func(c *mock_controller.MockIController)
	testCases := []struct {
		name                 string
		mockBehavior         mockBehavior
		httpMethod           string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "test_ok",
			mockBehavior: func(c *mock_controller.MockIController) {
				tables := []model.Table{
					{
						Name:    "animals",
						Columns: []string{},
					},
					{
						Name:    "cities",
						Columns: []string{},
					},
				}
				c.EXPECT().GetTables().Return(tables)
			},
			httpMethod:           "GET",
			expectedStatusCode:   200,
			expectedResponseBody: `[{"name":"animals","columns":[]},{"name":"cities","columns":[]}]`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			controller := mock_controller.NewMockIController(c)
			testCase.mockBehavior(controller)

			handlers := &Handlers{Controller: controller}
			router := mux.NewRouter()
			router.HandleFunc("/", handlers.GetTablesHandler).Methods(testCase.httpMethod)
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(testCase.httpMethod, "/", nil)
			router.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, testCase.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestGetEntriesHandler(t *testing.T) {
	type mockBehavior = func(c *mock_controller.MockIController, tableName string, limit int64, offset int64)
	testCases := []struct {
		name                 string
		mockBehavior         mockBehavior
		tableName            string
		limit                int64
		offset               int64
		httpMethod           string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "test_ok",
			mockBehavior: func(c *mock_controller.MockIController, tableName string, limit int64, offset int64) {
				rows := []model.Row{
					{
						ColumnsByValues: map[string]interface{}{
							"id":   1,
							"kind": "cat",
						},
					},
					{
						ColumnsByValues: map[string]interface{}{
							"id":   2,
							"kind": "lion",
						},
					},
				}
				c.EXPECT().GetRows(tableName, limit, offset).Return(rows, nil)
			},
			tableName:            "animals",
			limit:                2,
			offset:               0,
			httpMethod:           "GET",
			expectedStatusCode:   200,
			expectedResponseBody: `[{"columnsByValues":{"id":1,"kind":"cat"}},{"columnsByValues":{"id":2,"kind":"lion"}}]`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			controller := mock_controller.NewMockIController(c)
			testCase.mockBehavior(controller, testCase.tableName, testCase.limit, testCase.offset)

			handlers := &Handlers{Controller: controller}
			router := mux.NewRouter()
			router.HandleFunc("/{table}", handlers.GetEntriesHandler).Methods(testCase.httpMethod)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s?limit=%d&offset=%d", testCase.tableName, testCase.limit, testCase.offset)
			request := httptest.NewRequest(testCase.httpMethod, url, nil)

			router.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, testCase.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestGetEntryHandler(t *testing.T) {
	type mockBehavior = func(c *mock_controller.MockIController, tableName string, id int64)
	testCases := []struct {
		name                 string
		mockBehavior         mockBehavior
		tableName            string
		id                   int64
		httpMethod           string
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "test_ok",
			mockBehavior: func(c *mock_controller.MockIController, tableName string, id int64) {
				row := model.Row{
					ColumnsByValues: map[string]interface{}{
						"id":   1,
						"kind": "cat",
					},
				}
				c.EXPECT().GetRow(tableName, id).Return(row, nil)
			},
			tableName:            "animals",
			id:                   1,
			httpMethod:           "GET",
			expectedStatusCode:   200,
			expectedResponseBody: `{"columnsByValues":{"id":1,"kind":"cat"}}`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			controller := mock_controller.NewMockIController(c)
			testCase.mockBehavior(controller, testCase.tableName, testCase.id)

			handlers := &Handlers{Controller: controller}
			router := mux.NewRouter()
			router.HandleFunc("/{table}/{id}", handlers.GetEntryHandler).Methods(testCase.httpMethod)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/%d", testCase.tableName, testCase.id)
			request := httptest.NewRequest(testCase.httpMethod, url, nil)

			router.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, testCase.expectedStatusCode)
			assert.Equal(t, recorder.Body.String(), testCase.expectedResponseBody)
		})
	}
}

func TestCreateEntryHandler(t *testing.T) {
	type mockBehavior = func(c *mock_controller.MockIController, tableName string, row model.Row)
	testCases := []struct {
		name               string
		mockBehavior       mockBehavior
		tableName          string
		inputRow           model.Row
		httpMethod         string
		expectedStatusCode int
	}{
		{
			name: "test_ok",
			mockBehavior: func(c *mock_controller.MockIController, tableName string, row model.Row) {
				c.EXPECT().CreateRow(tableName, row).Return(nil)
			},
			tableName: "animals",
			inputRow: model.Row{
				ColumnsByValues: map[string]interface{}{
					"id":   1,
					"kind": "cat",
				},
			},
			httpMethod:         "PUT",
			expectedStatusCode: 200,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			controller := mock_controller.NewMockIController(c)
			testCase.mockBehavior(controller, testCase.tableName, testCase.inputRow)

			handlers := &Handlers{Controller: controller}
			router := mux.NewRouter()
			router.HandleFunc("/{table}", handlers.CreateEntryHandler).Methods(testCase.httpMethod)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s", testCase.tableName)
			request := httptest.NewRequest(testCase.httpMethod, url, nil)
			request.Form = make(map[string][]string, 0)
			for column, value := range testCase.inputRow.ColumnsByValues {
				switch castedValue := value.(type) {
				case string:
					request.Form.Add(column, castedValue)
				case int:
					request.Form.Add(column, strconv.Itoa(castedValue))
				case float32:
				case float64:
					request.Form.Add(column, fmt.Sprintf("%f", castedValue))
				default:
				}
			}
			router.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, testCase.expectedStatusCode)
		})
	}
}

func TestUpdateEntryHandler(t *testing.T) {
	type mockBehavior = func(c *mock_controller.MockIController, tableName string, row model.Row)
	testCases := []struct {
		name               string
		mockBehavior       mockBehavior
		tableName          string
		inputRow           model.Row
		id                 int64
		httpMethod         string
		expectedStatusCode int
	}{
		{
			name: "test_ok",
			mockBehavior: func(c *mock_controller.MockIController, tableName string, row model.Row) {
				c.EXPECT().UpdateRow(tableName, row).Return(nil)
			},
			tableName: "animals",
			inputRow: model.Row{
				ColumnsByValues: map[string]interface{}{
					"id":   1,
					"kind": "cat",
				},
			},
			id:                 1,
			httpMethod:         "POST",
			expectedStatusCode: 200,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			controller := mock_controller.NewMockIController(c)
			testCase.mockBehavior(controller, testCase.tableName, testCase.inputRow)

			handlers := &Handlers{Controller: controller}
			router := mux.NewRouter()
			router.HandleFunc("/{table}/{id}", handlers.UpdateEntryHandler).Methods(testCase.httpMethod)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/%d", testCase.tableName, testCase.id)
			request := httptest.NewRequest(testCase.httpMethod, url, nil)
			request.Form = make(map[string][]string, 0)
			for column, value := range testCase.inputRow.ColumnsByValues {
				switch castedValue := value.(type) {
				case string:
					request.Form.Add(column, castedValue)
				case int:
					request.Form.Add(column, strconv.Itoa(castedValue))
				case float32:
				case float64:
					request.Form.Add(column, fmt.Sprintf("%f", castedValue))
				default:
				}
			}
			router.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, testCase.expectedStatusCode)
		})
	}
}

func TestDeleteEntryHandler(t *testing.T) {
	type mockBehavior = func(c *mock_controller.MockIController, tableName string, id int64)
	testCases := []struct {
		name               string
		mockBehavior       mockBehavior
		tableName          string
		id                 int64
		httpMethod         string
		expectedStatusCode int
	}{
		{
			name: "test_ok",
			mockBehavior: func(c *mock_controller.MockIController, tableName string, id int64) {
				c.EXPECT().DeleteRow(tableName, id).Return(nil)
			},
			tableName:          "animals",
			id:                 1,
			httpMethod:         "DELETE",
			expectedStatusCode: 200,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()
			controller := mock_controller.NewMockIController(c)
			testCase.mockBehavior(controller, testCase.tableName, testCase.id)

			handlers := &Handlers{Controller: controller}
			router := mux.NewRouter()
			router.HandleFunc("/{table}/{id}", handlers.DeleteEntryHandler).Methods(testCase.httpMethod)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/%s/%d", testCase.tableName, testCase.id)
			request := httptest.NewRequest(testCase.httpMethod, url, nil)
			router.ServeHTTP(recorder, request)

			assert.Equal(t, recorder.Code, testCase.expectedStatusCode)
		})
	}
}
