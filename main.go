package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

type Api struct {
	Endpoint string `yaml:"endpoint"`
	Method   string `yaml:"method"`
	Response string `yaml:"response"`
	Status   int    `yaml:"status"`
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	renderBackendConfig(r, "./backend.yaml")
	r.Run()
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func renderBackendConfig(r *gin.Engine, configFilePath string) *gin.Engine {
	data, err := os.ReadFile(configFilePath)
	check(err)
	var apis []Api
	err = yaml.Unmarshal(data, &apis)
	check(err)
	for _, api := range apis {
		_endpoint := api.Endpoint
		_handlerFunc := handlerFile(api.Response, api.Status)
		switch api.Method {
		case "get":
			r.GET(_endpoint, _handlerFunc)
		case "post":
			r.POST(_endpoint, _handlerFunc)
		case "delete":
			r.DELETE(_endpoint, _handlerFunc)
		case "put":
			r.DELETE(_endpoint, _handlerFunc)
		}
	}
	return r
}

func handlerFile(sampleResponseFilePath string, status int) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := os.ReadFile(sampleResponseFilePath)
		check(err)
		var resp interface{}
		err = json.Unmarshal(data, &resp)
		check(err)
		c.JSON(status, resp)
	}
}
