// SPDX-FileCopyrightText: 2022 Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// This function is not autogenerated
func AddService(engine *gin.Engine) *gin.RouterGroup {
	group := engine.Group("app/v1")

	for _, route := range routes {
		switch route.Method {
		case "GET":
			group.GET(route.Pattern, route.HandlerFunc)
		case "POST":
			group.POST(route.Pattern, route.HandlerFunc)
		case "PUT":
			group.PUT(route.Pattern, route.HandlerFunc)
		case "DELETE":
			group.DELETE(route.Pattern, route.HandlerFunc)
		}
	}
	return group

}

var routes = Routes{
	{
		"Index",
		http.MethodGet,
		"/",
		Index,
	},

	{
		"AcquireIntegerResource",
		http.MethodPost,
		"/integer-resource/:resource-name", // e.g. "/integer-resource/resourceids?number=1"
		IntegerResourceNamePost,
	},

	{
		"ReleaseIntegerResource",
		http.MethodDelete,
		"/integer-resource/:resource-name/:resource-id", // e.g. "/integer-resource/resourceids/10"
		IntegerResourceNameDelete,
	},

	{
		"AcquireIpv4Resource",
		http.MethodPost,
		"/ipv4-resource/:resource-name", // "/ipv4-resource"
		Ipv4ResourceNamePost,
	},

	{
		"ReleaseIpv4Resource",
		http.MethodDelete,
		"/ipv4-resource/:resource-name/:resource-id", // "/integer-resource/pool1/1.1.1.1"
		Ipv4ResourceNameDelete,
	},

	{
		"studentRecordTest",
		http.MethodPost,
		"/student-record-test",
		StudentRecordTest,
	},

	{
		"getUniqueIdentityTest",
		http.MethodPost,
		"/unique-identity/:pool",
		GetUniqueIdentityTest,
	},

	{
		"getUniqueIdentityWithinRangeTest",
		http.MethodPost,
		"/unique-identity-range/",
		GetUniqueIdentityWithinRangeTest,
	},

	{
		"getIdFromPoolTest",
		http.MethodPost,
		"/id-from-pool/:pool/",
		GetIdFromPoolTest,
	},

	{
		"getIdFromInsertPoolTest",
		http.MethodPost,
		"/id-from-insert-pool/:pool/",
		GetIdFromInsertPoolTest,
	},

	{
		"getChunkFromPoolTest",
		http.MethodPost,
		"/chunk-from-pool/:chunk/",
		GetChunkFromPoolTest,
	},
}

func http_server() {

	engine := gin.New()
	AddService(engine)

	engine.Use(cors.New(cors.Config{
		AllowMethods: []string{"GET", "POST", "OPTIONS", "PUT", "PATCH", "DELETE"},
		AllowHeaders: []string{
			"Origin", "Content-Length", "Content-Type", "User-Agent",
			"Referrer", "Host", "Token", "X-Requested-With",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	httpAddr := ":" + strconv.Itoa(8000)
	engine.Run(httpAddr)
	log.Println("Webserver stopped/terminated/not-started ")
}

func Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

func IntegerResourceNamePost(c *gin.Context) {
	c.String(http.StatusOK, "IntegerResourceNamePost!")
	resName, exists := c.Params.Get("resource-name")
	if exists == false {
		log.Printf("Received resource delete. resource-name not found ")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	number, ok := c.GetQuery("number")
	if ok == true {
		n1, _ := strconv.Atoi(number)
		var n int32
		n = int32(n1)
		resId := AllocateInt32Many(resName, n)
		if len(resId) == 0 {
			log.Println("Id allocation error ")
			c.JSON(http.StatusBadRequest, gin.H{})
		}
		log.Printf("Received resource create. Pool name %v, Pool Id %v ", resName, resId)
		c.JSON(http.StatusOK, gin.H{})

	} else {
		resId := AllocateInt32One(resName)
		if resId == 0 {
			log.Println("Id allocation error ")
			c.JSON(http.StatusBadRequest, gin.H{})
		}
		log.Printf("Received resource create. Pool name %v, Pool Id %v ", resName, resId)
		c.JSON(http.StatusOK, gin.H{})
	}
	return
}

func IntegerResourceNameDelete(c *gin.Context) {
	c.String(http.StatusOK, "IntegerResourceNameDelete!")
	resName, exists := c.Params.Get("resource-name")
	if exists == false {
		log.Printf("Received resource delete. resource-name not found ")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	log.Printf("Received resource delete. Pool name %v ", resName)

	resId, exists := c.Params.Get("resource-id")
	if exists == false {
		log.Printf("resource-id param not found ")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	log.Printf("Received resource delete. Res Id %v ", resId)
	r, _ := strconv.Atoi(resId)
	rid := int32(r)

	err := ReleaseInt32One(resName, rid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
	return
}

func Ipv4ResourceNamePost(c *gin.Context) {
	c.String(http.StatusOK, "Ipv4ResourceNamePost!")
	resName, exists := c.Params.Get("resource-name")
	if exists == false {
		log.Printf("Received resource delete. resource-name not found ")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	number, ok := c.GetQuery("number")
	if ok == true {
		n1, _ := strconv.Atoi(number)
		var n int32
		n = int32(n1)
		resId := IpAddressAllocMany(resName, n)
		if len(resId) == 0 {
			log.Println("Id allocation error ")
			c.JSON(http.StatusBadRequest, gin.H{})
		}
		log.Printf("Received resource create. Pool name %v, Pool Id %v ", resName, resId)
		c.JSON(http.StatusOK, gin.H{})

	} else {
		resId, err := IpAddressAllocOne(resName)
		if err != nil {
			log.Println("Id allocation error ")
			c.JSON(http.StatusBadRequest, gin.H{})
		}
		log.Printf("Received resource create. Pool name %v, Pool Id %v ", resName, resId)
		c.JSON(http.StatusOK, gin.H{})
	}
	return
}

func Ipv4ResourceNameDelete(c *gin.Context) {
	c.String(http.StatusOK, "Ipv4ResourceNameDelete!")
	resName, exists := c.Params.Get("resource-name")
	if exists == false {
		log.Printf("Received resource delete. resource-name not found ")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	log.Printf("Received resource delete. Pool name %v ", resName)

	resId, exists := c.Params.Get("resource-id")
	if exists == false {
		log.Printf("resource-id param not found ")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	log.Printf("Received resource delete. Res Id %v ", resId)

	err := IpAddressRelease(resName, resId)
	if err != nil {
		log.Printf("IP address %v release failed -  %v ", resId, err)
		c.JSON(http.StatusBadRequest, gin.H{})
	} else {
		log.Printf("IP address %v release success ", resId)
		c.JSON(http.StatusOK, gin.H{})
	}
	return
}

func GetUniqueIdentityTest(c *gin.Context) {
	c.String(http.StatusOK, "GetUniqueIdentityTest!")
	resName, exists := c.Params.Get("pool")
	if exists == false {
		log.Printf("pool param missing in URI")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	uniqueId := mongoHndl.GetUniqueIdentity(resName)
	log.Println(uniqueId)
	c.JSON(http.StatusOK, gin.H{})
}

func GetUniqueIdentityWithinRangeTest(c *gin.Context) {
	c.String(http.StatusOK, "GetUniqueIdentityWithinRangeTest!")
	resName, exists := c.Params.Get("pool")
	if exists == false {
		log.Printf("pool param missing in URI")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	min, ok := GetQuery("min", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	max, ok := GetQuery("max", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	uniqueId := mongoHndl.GetUniqueIdentityWithinRange(resName, min, max)
	log.Println(uniqueId)
	c.JSON(http.StatusOK, gin.H{})
}

func GetIdFromPoolTest(c *gin.Context) {
	log.Println("TESTING POOL OF IDS")

	poolName, exists := c.Params.Get("pool")
	if exists == false {
		log.Printf("pool param missing in URI")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	min, ok := GetQuery("min", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	max, ok := GetQuery("max", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	mongoHndl.InitializePool(poolName, min, max)

	uniqueId, err := mongoHndl.GetIDFromPool(poolName)
	log.Println(uniqueId, err)

	mongoHndl.ReleaseIDToPool(poolName, uniqueId)

	uniqueId, err = mongoHndl.GetIDFromPool(poolName)
	log.Println(uniqueId, err)

	uniqueId, err = mongoHndl.GetIDFromPool(poolName)
	log.Println(uniqueId, err)
	c.JSON(http.StatusOK, gin.H{})

}

func GetIdFromInsertPoolTest(c *gin.Context) {
	log.Println("TESTING INSERT APPROACH")
	var randomId int32

	pool, exists := c.Params.Get("pool")
	if exists == false {
		log.Printf("pool param missing in URI")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	min, ok := GetQuery("min", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	max, ok := GetQuery("max", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	retry, ok := GetQuery("retry", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	mongoHndl.InitializeInsertPool(pool, min, max, retry)

	randomId, err := mongoHndl.GetIDFromInsertPool(pool)
	log.Println(randomId)
	if err != nil {
		log.Println(err.Error())
	}

	randomId, err = mongoHndl.GetIDFromInsertPool(pool)
	log.Println(randomId)
	if err != nil {
		log.Println(err.Error())
	}

	randomId, err = mongoHndl.GetIDFromInsertPool(pool)
	log.Println(randomId)
	if err != nil {
		log.Println(err.Error())
	}

	mongoHndl.ReleaseIDToInsertPool(pool, randomId)
	c.JSON(http.StatusOK, gin.H{})
}

func GetChunkFromPoolTest(c *gin.Context) {
	log.Println("TESTING CHUNK APPROACH")
	var lower int32
	var upper int32

	resName, exists := c.Params.Get("chunk")
	if exists == false {
		log.Printf("chunk param missing in URI")
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	min, ok := GetQuery("min", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	max, ok := GetQuery("max", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	retry, ok := GetQuery("retry", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	csize, ok := GetQuery("csize", c)
	if ok == false {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}

	mongoHndl.InitializeChunkPool(resName, min, max, retry, csize) // min, max, retries, chunkSize

	randomId, lower, upper, err := mongoHndl.GetChunkFromPool(resName)
	log.Println(randomId, lower, upper)
	if err != nil {
		log.Println(err.Error())
	}

	randomId, lower, upper, err = mongoHndl.GetChunkFromPool(resName)
	log.Println(randomId, lower, upper)
	if err != nil {
		log.Println(err.Error())
	}

	randomId, lower, upper, err = mongoHndl.GetChunkFromPool(resName)
	log.Println(randomId, lower, upper)
	if err != nil {
		log.Println(err.Error())
	}

	mongoHndl.ReleaseChunkToPool(resName, randomId)

	c.JSON(http.StatusOK, gin.H{})

}

func GetQuery(param string, c *gin.Context) (int32, bool) {
	p1, ok := c.GetQuery(param)
	if ok == false {
		return 0, false
	}
	p2, _ := strconv.Atoi(p1)
	p := int32(p2)
	return p, true
}
