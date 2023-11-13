package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/dshargool/go-mbslave-api.git/pkg/types"
)

type testHandler struct {
	handler Handler
}

func setupTestSuite() testHandler {
	testConfig := types.Configuration{
		Port:      8081,
		DBPath:    "test/data/test.db",
		Registers: map[types.OpcTag]types.ModbusTag{},
	}

	myDb := types.SqlDb{}
	myDb.Open(testConfig.DBPath)

	testRegisters := []types.ModbusTag{
		{
			Tag:         "TestTag",
			Description: "Test",
			Address:     1,
			Divisor:     10,
		},
		{
			Tag:         "ValidTag",
			Description: "Test",
			Address:     2,
			Divisor:     10,
		},
	}
	for _, register := range testRegisters {
		testConfig.Registers[types.OpcTag(register.Tag)] = register
	}

	myDb.CreateTable()
	myDb.UpdateTableTags(testConfig.Registers)
	// Set a valid value to our 'ValidTag' address in the test db
	myDb.SetAddressValue(2, 100.0)

	myHandler := New(testConfig, &myDb)

	return testHandler{
		handler: myHandler,
	}
}

func (h testHandler) cleanUp() {
	h.handler.db.DB.Close()
}

func TestGetRegisters(t *testing.T) {
	expected := 200

	testHandler := setupTestSuite()

	request, _ := http.NewRequest(http.MethodGet, "/all_registers", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetRegisters(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, want %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestGetNullValueTag(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 500

	request, _ := http.NewRequest(http.MethodGet, "/tag/TestTag", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetTag(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestGetNullRegister(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 500

	request, _ := http.NewRequest(http.MethodGet, "/register/1", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetRegister(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestGetUnknownTag(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 404

	request, _ := http.NewRequest(http.MethodGet, "/tag/UnknownTag", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetTag(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestGetUnknownRegister(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 404

	request, _ := http.NewRequest(http.MethodGet, "/register/0", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetRegister(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestGetValidTag(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 200

	request, _ := http.NewRequest(http.MethodGet, "/tag/ValidTag", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetTag(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestPutValidTag(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 200

	data := url.Values{}
	data.Add("value", "100")

	request, _ := http.NewRequest(http.MethodPut, "/tag/ValidTag", nil)
	request.URL.RawQuery = data.Encode()
	response := httptest.NewRecorder()

	testHandler.handler.GetTag(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestGetValidRegister(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 200

	request, _ := http.NewRequest(http.MethodGet, "/register/2", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetRegister(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestPutValidRegister(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 200

	data := url.Values{}
	data.Add("value", "100")

	request, _ := http.NewRequest(http.MethodPut, "/register/2", nil)
	request.URL.RawQuery = data.Encode()
	response := httptest.NewRecorder()

	testHandler.handler.GetRegister(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}
