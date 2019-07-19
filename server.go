package main

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type (
	holder struct {
		Id_        bson.ObjectId `json:"id" bson:"_id,omitempty"`
		Address    string        `json:address`
		Quantity   float64       `json:quantity`
		Percentage float64       `json:percentage`
		Tag        string        `json:tag`
	}

	asset_holder struct {
		Address    string  `json:address`
		Quantity   float64 `json:quantity`
		Percentage float64 `json:percentage`
		Tag        string  `json:tag`
	}

	asset_holders struct {
		TotalNum       int            `json:"totalNum"`
		AddressHolders []asset_holder `json:"addressHolders"`
	}
)

func connect(cName string) *mgo.Collection {
	// Database connection
	session, err := mgo.Dial("localhost:27017")
	if err != nil {
		panic(err)
	}
	// defer session.Close()
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("spider").C(cName)
	return c
}

// func createHolder(c echo.Context) error {
// 	u := &holder{
// 		ID: c.Param("id"),
// 	}
// 	if err := c.Bind(u); err != nil {
// 		return err
// 	}
// 	return c.JSON(http.StatusCreated, u)
// }

// func updateHolder(c echo.Context) error {
// 	u := new(holder)
// 	if err := c.Bind(u); err != nil {
// 		return err
// 	}
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	// holders[id].Name = u.Name
// 	return c.JSON(http.StatusOK, holders[id])
// }

// func deleteHolder(c echo.Context) error {
// 	id, _ := strconv.Atoi(c.Param("id"))
// 	delete(holders, id)
// 	return c.NoContent(http.StatusNoContent)
// }

func getHolder(c echo.Context) error {
	address := c.Param("address")

	db := connect("holder")
	result := holder{}

	err := db.Find(bson.M{"address": address}).One(&result)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: err}
	}
	return c.JSON(http.StatusOK, result)
}

func getHolders(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	rows, _ := strconv.Atoi(c.QueryParam("rows"))
	if page <= 0 {
		page = 1
	}
	if rows <= 0 {
		rows = 10
	}

	db := connect("holder")
	var results []holder

	err := db.Find(nil).Sort("$quantity").Limit(rows).Skip((page - 1) * rows).All(&results)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: err}
	}
	return c.JSON(http.StatusOK, results)
}

func main() {
	e := echo.New()
	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	// e.POST("/holders", createHolder)
	e.GET("/holders/:address", getHolder)
	e.GET("/holders", getHolders)
	// e.PUT("/holders/:id", updateHolder)
	// e.DELETE("/holders/:id", deleteHolder)
	e.Logger.Fatal(e.Start(":1323"))

}
