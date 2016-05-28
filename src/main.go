package main

import (
	"fmt"
	"time"
	"github.com/gin-gonic/gin"
	expand "github.com/openvenues/gopostal/expand"
	parser "github.com/openvenues/gopostal/parser"
	"github.com/itsjamie/gin-cors"
)

func parseAddrToMap(v string) (m map[string]string) {
	parsedComponents := parser.ParseAddress(v)
	m = make(map[string]string)
	for _, v := range parsedComponents {
		m[v.Label] = v.Value
	}
	return
}

func parseAddress(c *gin.Context) {
	address := c.Query("address")
	m := parseAddrToMap(address)
	c.JSON(200, m)
}

func expandAddress(c *gin.Context) {
	expansions := expand.ExpandAddress(c.Query("address"))
	c.JSON(200, expansions)
}

type MultiAddrs struct {
	Addresses []string `json:"addresses"`
}

type ParsingResult struct {
	Label string `json:label`
	Value string `json:value`
}

type MultiAddrResult struct {
	Input string `json:input`
	Result []ParsingResult
}

func parseMultipleAddresses(c *gin.Context) {
	//	fmt.Println(c.Request.Body)
	addrs := &MultiAddrs{}
	c.Bind(addrs)
	if len(addrs.Addresses) > 500 {
		c.JSON(422, gin.H{
			"error": "Limit of 500 addresses per call",
		})
	} else {
		result := make([]MultiAddrResult, len(addrs.Addresses))
		for i, v := range addrs.Addresses {
			result[i].Input = v
			parsed := parser.ParseAddress(v)
			result[i].Result = make([]ParsingResult, len(parsed))
			for j, component := range parsed {
				result[i].Result[j].Label = component.Label
				result[i].Result[j].Value = component.Value
			}
		}
		c.JSON(200, result)
	}
}

func main() {
	fmt.Println("Starting the server")
	r := gin.Default()
	// Apply the middleware to the router (works with groups too)
	r.Use(cors.Middleware(cors.Config{
    Origins:        "*",
    Methods:        "GET, PUT, POST, DELETE",
    RequestHeaders: "Origin, Authorization, Content-Type",
    ExposedHeaders: "",
    MaxAge: 50 * time.Second,
    Credentials: true,
    ValidateHeaders: false,
	}))
	r.GET("/parse", parseAddress)
	r.GET("/expand", expandAddress)
	r.POST("/parse_multi", parseMultipleAddresses)
	r.Run() // listen and server on 0.0.0.0:8080
}
