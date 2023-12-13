package handlers

import (
	"encoding/json"
	"fmt"
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
	fmt.Println("New test!")
	testConfig := types.Configuration{
		ApiPort:    8081,
		ModbusPort: 5502,
		DBPath:     "test/data/test.db",
		Registers:  map[types.InstrumentTag]types.ModbusTag{},
	}

	myDb := types.SqlDb{}
	_ = myDb.Open(testConfig.DBPath)

	testRegisters := []types.ModbusTag{
		{
			Tag:         "TestTagF32",
			Description: "Test",
			Address:     2,
			DataType:    "float32",
		},
		{
			Tag:         "ValidTagF32",
			Description: "Test",
			Address:     4,
			DataType:    "float32",
		},
	}
	for _, register := range testRegisters {
		testConfig.Registers[types.InstrumentTag(register.Tag)] = register
	}

	_ = myDb.CreateTable()
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

func TestGetNullValueTagF32(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 500

	request, _ := http.NewRequest(http.MethodGet, "/tag/TestTagF32", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetTag(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestGetNullRegisterF32(t *testing.T) {
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

	request, _ := http.NewRequest(http.MethodGet, "/tag/ValidTagF32", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetTag(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestPutValidTagF32(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 200

	data := url.Values{}
	data.Add("value", "100")

	request, _ := http.NewRequest(http.MethodPut, "/tag/ValidTagF32", nil)
	request.URL.RawQuery = data.Encode()
	response := httptest.NewRecorder()

	testHandler.handler.GetTag(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestGetValidRegisterF32(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 200

	request, _ := http.NewRequest(http.MethodGet, "/register/3", nil)
	response := httptest.NewRecorder()

	testHandler.handler.GetRegister(response, request)
	res := response.Result().StatusCode

	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestPutValidRegisterF32(t *testing.T) {
	testHandler := setupTestSuite()
	expected := 200

	data := url.Values{}
	data.Add("value", "100")

	request, _ := http.NewRequest(http.MethodPut, "/register/3", nil)
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
	request, _ := http.NewRequest(http.MethodPut, "/register/3", nil)
	request.URL.RawQuery = data.Encode()

	testHandler.handler.GetRegister(response, request)
	response = httptest.NewRecorder()

	request, _ = http.NewRequest(http.MethodGet, "/register/3", nil)
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

func TestModbusGetValidAddressF32(t *testing.T) {
	testHandler := setupTestSuite()
	var expected float32 = 100.0

	mbClient := testHandler.mb_client
	_ = mbClient.Open()

	res, _ := mbClient.ReadFloat32(2, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %.2f, expected %.2f", res, expected)
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

	res, _ := mbClient.ReadRegister(60, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %d, expected %d", res, expected)
	}
	testHandler.cleanUp()
}

func TestModbusSetValidAddressF32(t *testing.T) {
	testHandler := setupTestSuite()
	var expected float32 = 100.0

	mbClient := testHandler.mb_client

	_ = mbClient.WriteFloat32(3, expected)
	res, _ := mbClient.ReadFloat32(3, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %.2f, expected %.2f", res, expected)
	}
	testHandler.cleanUp()
}

func TestModbusReadWritebackF32(t *testing.T) {
	testHandler := setupTestSuite()
	mbClient := testHandler.mb_client

	expected, _ := mbClient.ReadFloat32(3, modbus.HOLDING_REGISTER)
	_ = mbClient.WriteFloat32(3, expected)
	res, _ := mbClient.ReadFloat32(3, modbus.HOLDING_REGISTER)
	if res != expected {
		t.Errorf("Got %.2f, expected %.2f", res, expected)
	}
	testHandler.cleanUp()
}

func TestModbusWriteApiReadF32(t *testing.T) {
	testHandler := setupTestSuite()
	mbClient := testHandler.mb_client

	mbValue, _ := mbClient.ReadFloat32(3, modbus.HOLDING_REGISTER)
	_ = mbClient.WriteFloat32(3, mbValue)

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/register/3", nil)
	testHandler.handler.GetRegister(response, request)

	//_ := response.Result().StatusCode
	dec := json.NewDecoder(response.Body)
	var apiValue types.ModbusResponse
	_ = dec.Decode(&apiValue)

	if float32(apiValue.Value) != mbValue {
		t.Errorf("Api %.2f, Modbus %.2f", apiValue.Value, mbValue)
	}
	testHandler.cleanUp()
}

func TestApiWriteModbusRead(t *testing.T) {
	testHandler := setupTestSuite()
	mbClient := testHandler.mb_client
	expected := float32(200.1234)

	data := url.Values{}
	data.Add("value", strconv.FormatFloat(float64(expected), 'f', -1, 32))

	response := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPut, "/register/3", nil)
	request.URL.RawQuery = data.Encode()

	testHandler.handler.GetRegister(response, request)

	mbValue, _ := mbClient.ReadFloat32(3, modbus.HOLDING_REGISTER)

	if expected != mbValue {
		t.Errorf("Got %.2f, expected %.2f", mbValue, expected)
	}
	testHandler.cleanUp()
}
