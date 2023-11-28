package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/dshargool/go-mbslave-api.git/pkg/types"
	"github.com/simonvetter/modbus"
)

type testHandler struct {
	handler   Handler
	mb_client *modbus.ModbusClient
}

func setupTestSuite() testHandler {
	testConfig := types.Configuration{
		ApiPort:    8081,
		ModbusPort: 5502,
		DBPath:     "test/data/test.db",
		Registers:  map[types.OpcTag]types.ModbusTag{},
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
	_ = myDb.SetAddressValue(2, 100.0)

	myHandler := New(testConfig, &myDb)

	myHandler.MbSlave = myHandler.MbInit(testConfig.ModbusPort)
	myHandler.MbStart()

	client, _ := modbus.NewClient(&modbus.ClientConfiguration{
		URL:     "tcp://localhost:" + strconv.Itoa(testConfig.ModbusPort),
		Timeout: 1 * time.Second,
	})
    _ = client.Open()

	var retHandler testHandler
	retHandler.handler = myHandler
	retHandler.mb_client = client

	return retHandler
}

func (h *testHandler) cleanUp() {
	h.handler.db.DB.Close()
	h.handler.MbStop()
	h.mb_client.Close()
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

func TestPutGetWriteback(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 200.0

	data := url.Values{}
	data.Add("value", strconv.FormatFloat(expected, 'f', -1, 64))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, "/register/2", nil)
	request.URL.RawQuery = data.Encode()

	testHandler.handler.GetRegister(response, request)
	response = httptest.NewRecorder()

	request, _ = http.NewRequest(http.MethodGet, "/register/2", nil)
	testHandler.handler.GetRegister(response, request)

	//_ := response.Result().StatusCode
	dec := json.NewDecoder(response.Body)
	var respValue types.ModbusResponse
	_ = dec.Decode(&respValue)

	if respValue.Value != expected {
		t.Errorf("Got %.2f, expected %.2f", respValue.Value, expected)
	}
	testHandler.cleanUp()
}

func TestModbusGetValidAddress(t *testing.T) {
	testHandler := setupTestSuite()
	var expected uint16 = 1000

	mbClient := testHandler.mb_client
	_ = mbClient.Open()

	res, _ := mbClient.ReadRegister(2, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestModbusGetNullValueAddress(t *testing.T) {
	testHandler := setupTestSuite()
	var expected uint16 = 0

	mbClient := testHandler.mb_client

	res, _ := mbClient.ReadRegister(1, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestModbusGetUnknownAddress(t *testing.T) {
	testHandler := setupTestSuite()
	var expected uint16 = 0

	mbClient := testHandler.mb_client

	res, _ := mbClient.ReadRegister(3, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestModbusSetValidAddress(t *testing.T) {
	testHandler := setupTestSuite()
	var expected uint16 = 100

	mbClient := testHandler.mb_client

	_ = mbClient.WriteRegister(2, expected)
	res, _ := mbClient.ReadRegister(2, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestModbusReadWriteback(t *testing.T) {
	testHandler := setupTestSuite()
	mbClient := testHandler.mb_client

	expected, _ := mbClient.ReadRegister(2, modbus.HOLDING_REGISTER)
	_ = mbClient.WriteRegister(2, expected)
	res, _ := mbClient.ReadRegister(2, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestModbusWriteApiRead(t *testing.T) {
	testHandler := setupTestSuite()
	mbClient := testHandler.mb_client

	mbValue, _ := mbClient.ReadRegister(2, modbus.HOLDING_REGISTER)
	_ = mbClient.WriteRegister(2, mbValue)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/register/2", nil)
	testHandler.handler.GetRegister(response, request)

	//_ := response.Result().StatusCode
	dec := json.NewDecoder(response.Body)
	var apiValue types.ModbusResponse
	_ = dec.Decode(&apiValue)

	if uint16(apiValue.Value) != mbValue/uint16(apiValue.Divisor) {
		t.Errorf("Got %.2f, expected %d", apiValue.Value, mbValue/uint16(apiValue.Divisor))
	}
	testHandler.cleanUp()
}

func TestApiWriteModbusRead(t *testing.T) {
	testHandler := setupTestSuite()
	mbClient := testHandler.mb_client
	expected := 200.0
	var divisor uint16 = 10

	data := url.Values{}
	data.Add("value", strconv.FormatFloat(expected, 'f', -1, 64))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, "/register/2", nil)
	request.URL.RawQuery = data.Encode()

	testHandler.handler.GetRegister(response, request)

	mbValue, _ := mbClient.ReadRegister(2, modbus.HOLDING_REGISTER)

	if uint16(expected) != mbValue/divisor {
		t.Errorf("Got %d, expected %.2f", mbValue, expected*float64(divisor))
	}
	testHandler.cleanUp()
}
