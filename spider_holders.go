package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

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

// func autoSpider(e *echo.Echo) {
func autoSpider() {

	// 获取数据
	page := 1
	rows := 1000

	for {
		url := "https://explorer.binance.org/api/v1/asset-holders?page=" + strconv.Itoa(page) + "&rows=" + strconv.Itoa(rows) + "&asset=COS-2E4"
		fmt.Println(fmt.Sprintf("page:%v method:%v,all:%t", page, url, page*rows))

		resp, _ := http.Get(url)
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		asset_holders := asset_holders{}

		err := json.Unmarshal([]byte(body), &asset_holders)
		if err != nil {
			fmt.Printf("Unmarshal err, %v\n", err)
			continue
		}

		c := connect("holder")
		len := len(asset_holders.AddressHolders)
		// 插入数据库
		for i := 0; i < len; i++ {
			holderData := bson.M{
				"address":    asset_holders.AddressHolders[i].Address,
				"quantity":   asset_holders.AddressHolders[i].Quantity,
				"percentage": asset_holders.AddressHolders[i].Percentage,
				"tag":        asset_holders.AddressHolders[i].Tag,
			}
			find := bson.M{"address": asset_holders.AddressHolders[i].Address}
			update := bson.M{"$set": holderData}
			_, err := c.Upsert(find, update)
			if err != nil {
				fmt.Println("update err, %t", err)
			}
		}
		if page*rows >= asset_holders.TotalNum {
			break
		}
		// 页数+1
		page++
	}
	// time.AfterFunc(30*time.Minute, autoSpider)
}

func main() {
	i := 0
	for {
		i++
		fmt.Println("times:%t", i)
		autoSpider()
		time.Sleep(30 * time.Minute)
	}
}
