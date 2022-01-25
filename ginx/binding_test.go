package ginx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type testRequest struct {
	Id      int64   `uri:"id" binding:"required,min=10,max=100"`
	Name    string  `query:"name" binding:"required,max=10"`
	Amount  float64 `json:"amount" binding:"required,max=100"`
	Age     int     `json:"age" binding:"required,min=1,max=130"`
	Address string  `json:"address" binding:"required"`
}

func TestBindAndValidate(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{
		{
			Key:   "id",
			Value: "99",
		},
	}

	m := make(map[string]interface{})
	m["amount"] = 100
	m["age"] = 18
	m["address"] = "test address"
	marshal, _ := json.Marshal(m)
	r := httptest.NewRequest(http.MethodGet, "/hello", bytes.NewReader(marshal))

	query := r.URL.Query()
	query.Add("name", "test")
	r.URL.RawQuery = query.Encode()
	c.Request = r

	req := testRequest{}
	if err := BindAndValidate(c, &req); err != nil {
		t.Fatal(err)
	}
	data, _ := json.Marshal(&req)
	fmt.Println(string(data))
}

func init() {
	gin.SetMode(gin.TestMode)
}
