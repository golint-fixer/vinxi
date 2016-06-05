package config

import (
	"encoding/json"
	"github.com/nbio/st"
	"testing"
)

const (
	key = "key"
)

func TestGetString(t *testing.T) {
	c := Config{}
	c.Set(key, "1")
	st.Expect(t, "1", c.GetString(key))
}

func TestGetStringForMissingKey(t *testing.T) {
	c := Config{}
	st.Expect(t, "", c.GetString(key))
}

func TestGetBool(t *testing.T) {
	c := Config{}
	c.Set(key, true)
	st.Expect(t, true, c.GetBool(key))
}

func TestGetBoolForMissingKey(t *testing.T) {
	c := Config{}
	st.Expect(t, false, c.GetBool(key))
}

func TestGetInt(t *testing.T) {
	c := Config{}
	c.Set(key, 1)
	st.Expect(t, 1, c.GetInt(key))
}

func TestGetIntForMissingKey(t *testing.T) {
	c := Config{}
	st.Expect(t, 0, c.GetInt(key))
}

func TestGetInt64(t *testing.T) {
	c := Config{}
	c.Set(key, int64(1))
	st.Expect(t, int64(1), c.GetInt64(key))
}

func TestGetInt64ForMissingKey(t *testing.T) {
	c := Config{}
	st.Expect(t, int64(0), c.GetInt64(key))
}

func TestGetFloat(t *testing.T) {
	c := Config{}
	c.Set(key, float64(1.1))
	st.Expect(t, float64(1.1), c.GetFloat(key))
}

func TestGetFloatForMissingKey(t *testing.T) {
	c := Config{}
	st.Expect(t, float64(0), c.GetFloat(key))
}

func TestGet(t *testing.T) {
	c := Config{}
	c.Set(key, 1)
	st.Expect(t, 1, c.Get(key))
}

func TestGetForMissingKey(t *testing.T) {
	c := Config{}
	st.Expect(t, nil, c.Get(key))
}

func TestJSON(t *testing.T) {
	c := Config{}
	c.Set(key, 1)
	c.Set("another_key", "1")
	expected := []byte(`{"another_key":"1","key":1}`)
	res, err := c.JSON()
	st.Expect(t, nil, err)
	st.Expect(t, expected, res)
}

func TestJSONError(t *testing.T) {
	c := Config{}
	c.Set(key, make(chan bool))
	res, err := c.JSON()
	err, ok := err.(*json.UnsupportedTypeError)
	st.Expect(t, true, ok)
	st.Expect(t, "json: unsupported type: chan bool", err.Error())
	st.Expect(t, 0, len(res))
}
