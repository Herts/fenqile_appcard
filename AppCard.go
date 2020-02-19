package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"github.com/spf13/viper"
	"log"
	"os"
	"time"
)

type System struct {
	Controller string `json:"controller"`
}
type Data struct {
	StateFilter string `json:"state_filter"`
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
}
type OrderInfoDetailReq struct {
	System System `json:"system"`
	Data   Data   `json:"data"`
}

func main() {
	csvWriter := CsvWriter("result.csv")
	csvWriter2 := CsvWriter("result-min.csv")
	viper.AddConfigPath(".")

	var cookieName string
	flag.StringVar(&cookieName, "config", "cookie", "cookie name in config")
	flag.Parse()

	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		log.Println("找不到cookie")
	}
	rawCookies := viper.GetString(cookieName)
	userAgent := "Mozilla/5.0 (iPhone; CPU iPhone OS 13_2_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1"
	req := gorequest.New()

	orders := getAllOrders(req, rawCookies, userAgent)
	var startedTime time.Time
	var smsCode string

	for i, order := range orders {
		productInfo := order.TemplateContent[1].OrderGoodsInfo.GoodsInfo.ProductInfo
		orderId := order.OrderInfo.OrderID
		createdTime := time.Unix(order.OrderInfo.CreateTime/1000, 0)
		fmt.Printf("%d, %s,%s,%s\n", i, productInfo, orderId, createdTime)
	}

	var startPoint, endPoint int
	fmt.Print("输入开始位置（包括）:")
	fmt.Scanln(&startPoint)
	fmt.Print("输入结束位置（包括）:")
	fmt.Scanln(&endPoint)

	orders = orders[startPoint : endPoint+1]

	for i, order := range orders {
		if i%5 == 0 {
			for time.Now().Sub(startedTime).Seconds() < 60 {
				fmt.Printf("距离上次发验证码只有 %f 秒\n", time.Now().Sub(startedTime).Seconds())
				time.Sleep(5 * time.Second)
			}
			smsCode = querySendSMS(req, userAgent, rawCookies)
			startedTime = time.Now()
		}
		if i > endPoint {
			break
		}
		productInfo := order.TemplateContent[1].OrderGoodsInfo.GoodsInfo.ProductInfo
		orderId := order.OrderInfo.OrderID
		createdTime := time.Unix(order.OrderInfo.CreateTime/1000, 0)
		orderInfoResp := getOrderFullInfo(smsCode, orderId, req, userAgent, rawCookies)
		var cardNum, cardPasswd string
		for _, info := range orderInfoResp.VirtualInfo.FuluInfo {
			cardNum, cardPasswd = info.CardNumber.Value, info.Passwd.Value
		}
		fmt.Printf("%d/%d, %s,%s,%s,%s,%s\n", i+1, len(orders), productInfo, orderId, createdTime, cardNum, cardPasswd)
		csvWriter.Write([]string{productInfo, orderId, createdTime.String(), cardNum, cardPasswd})
		csvWriter.Flush()
		csvWriter2.Write([]string{cardNum, cardPasswd})
		csvWriter2.Flush()
	}

}

func getAllOrders(req *gorequest.SuperAgent, rawCookies string, userAgent string) []orderInfo {
	offset := 0
	limit := 50

	orders := getOrders(offset, limit, req, rawCookies, userAgent)
	for {
		offset += limit
		limit += limit
		tmpOrders := getOrders(offset, limit, req, rawCookies, userAgent)
		if len(tmpOrders) == 0 {
			break
		}
		orders = append(orders, tmpOrders...)
	}
	return orders
}

func CsvWriter(filename string) *csv.Writer {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	csvWriter := csv.NewWriter(f)
	return csvWriter
}

func getOrderFullInfo(smsCode string, orderId string, req *gorequest.SuperAgent, userAgent string, rawCookies string) OrderFullInfoResp {
	type Payload struct {
		SendType int    `json:"send_type"`
		SmsCode  string `json:"sms_code"`
		OrderID  string `json:"order_id"`
		SaleType int    `json:"sale_type"`
		IsWeex   int    `json:"is_weex"`
	}
	data := Payload{
		SendType: 8,
		SmsCode:  smsCode,
		OrderID:  orderId,
		SaleType: 800,
		IsWeex:   1,
	}
	jsonBody, _ := json.Marshal(data)

	resp, _, errs := req.Post("https://trade.m.fenqile.com/order/query_verify_fulu_and_coupon_sms.json").
		Type("json").
		Send(string(jsonBody)).
		Set("User-Agent", userAgent).
		Set("Cookie", rawCookies).
		End()
	if errs != nil {
		fmt.Println(errs)
	}
	var orderInfoResp OrderFullInfoResp
	err := json.NewDecoder(resp.Body).Decode(&orderInfoResp)
	if err != nil {
		log.Println(err)
	}
	return orderInfoResp
}

func querySendSMS(req *gorequest.SuperAgent, userAgent string, rawCookies string) string {
	type QuerySendSMSReq struct {
		SendType int `json:"send_type"`
		IsWeex   int `json:"is_weex"`
	}
	queryBody := QuerySendSMSReq{
		SendType: 8,
		IsWeex:   1,
	}
	jsonBody, _ := json.Marshal(queryBody)
	_, body, errs := req.Post("https://trade.m.fenqile.com/order/query_send_sms.json").
		Type("json").
		Send(string(jsonBody)).
		Set("User-Agent", userAgent).
		Set("Cookie", rawCookies).
		End()
	if errs != nil {
		fmt.Println(errs)
	}
	fmt.Println(body)
	var smsCode string
	for {
		fmt.Print("Input sms code:")
		fmt.Scanln(&smsCode)
		if len(smsCode) == 6 {
			break
		} else {
			fmt.Print("输入验证码:")
		}
	}
	return smsCode
}

func getOrders(offset int, limit int, req *gorequest.SuperAgent, rawCookies, userAgent string) []orderInfo {
	reqBody := OrderInfoDetailReq{
		System: System{Controller: ""},
		Data:   Data{StateFilter: "", Offset: offset, Limit: limit},
	}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		log.Println(err)
	}

	resp, _, errs := req.Post("https://order.m.fenqile.com/route0001/order/getOrderInfoDetail.json").
		Type("json").
		Send(string(jsonBody)).
		Set("User-Agent", userAgent).
		Set("Cookie", rawCookies).
		End()
	if errs != nil {
		log.Println(errs)
	}
	var orderInfoDetail OrderInfoDetailResp
	err = json.NewDecoder(resp.Body).Decode(&orderInfoDetail)
	if err != nil {
		log.Println(err)
	}

	return orderInfoDetail.Data.ResultRows
}

func getOrderIds(resp OrderInfoDetailResp) []string {
	var orderIds []string
	for _, order := range resp.Data.ResultRows {
		orderIds = append(orderIds, order.OrderInfo.OrderID)
	}
	return orderIds
}
