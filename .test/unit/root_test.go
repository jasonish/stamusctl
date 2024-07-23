package tests

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"stamus-ctl/pkg"
)

func TestPing(t *testing.T) {
	res, _ := newRequest("GET", "/api/v1/ping", nil)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "{\"message\":\"pong\"}", res.Body.String())
}

func TestComposeInit(t *testing.T) {
	initRequest := pkg.InitRequest{
		IsDefault: true,
		Values: map[string]string{
			"nginx.exec": "nginx",
		},
	}

	res, _ := newRequest("POST", "/api/v1/compose/init", initRequest)
	if t != nil {
		assert.Equal(t, 200, res.Code)
		assert.Equal(t, "{\"message\":\"ok\"}", res.Body.String())
		compareDirs(t, "./tmp", ".././outputs/compose-init")
	}
}

func TestComposeInitSet(t *testing.T) {
	initRequest := pkg.InitRequest{
		IsDefault: true,
		Values: map[string]string{
			"nginx.exec":         "nginx",
			"websocket.response": "lel",
		},
	}

	res, _ := newRequest("POST", "/api/v1/compose/init", initRequest)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "{\"message\":\"ok\"}", res.Body.String())

	compareDirs(t, "./tmp", "../outputs/compose-init-set")
}

func TestComposeInitOptional(t *testing.T) {
	initRequest := pkg.InitRequest{
		IsDefault: true,
		Values: map[string]string{
			"nginx": "false",
		},
	}

	res, _ := newRequest("POST", "/api/v1/compose/init", initRequest)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "{\"message\":\"ok\"}", res.Body.String())

	compareDirs(t, "./tmp", "../outputs/compose-init-optional")
}

func TestComposeInitArbitrary(t *testing.T) {
	initRequest := pkg.InitRequest{
		IsDefault: true,
		Values: map[string]string{
			"nginx.exec":     "nginx",
			"websocket.port": "6969",
		},
	}

	res, _ := newRequest("POST", "/api/v1/compose/init", initRequest)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "{\"message\":\"ok\"}", res.Body.String())

	compareDirs(t, "./tmp", "../outputs/compose-init-arbitrary")

	err := os.RemoveAll("./tmp")
	assert.NoError(t, err)
}

func TestConfigSet(t *testing.T) {
	// Setup
	TestComposeInit(nil)
	// Test
	setRequest := pkg.SetRequest{
		Values: map[string]string{
			"websocket.response": "lel",
		},
	}
	res, _ := newRequest("POST", "/api/v1/config", setRequest)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "{\"message\":\"ok\"}", res.Body.String())
	// Compare
	compareDirs(t, "./tmp", "../outputs/compose-init-set")
}

func TestConfigReload(t *testing.T) {
	// Setup
	TestComposeInit(nil)
	setRequest := pkg.SetRequest{
		Values: map[string]string{
			"websocket.port": "6969",
		},
	}
	newRequest("POST", "/api/v1/config", setRequest)
	// Test
	setRequest = pkg.SetRequest{
		Reload: true,
	}
	res, _ := newRequest("POST", "/api/v1/config", setRequest)
	assert.Equal(t, 200, res.Code)
	assert.Equal(t, "{\"message\":\"ok\"}", res.Body.String())
	// Compare
	compareDirs(t, "./tmp", "../outputs/compose-init")
}
