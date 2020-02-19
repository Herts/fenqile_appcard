package main

type OrderInfoDetailResp struct {
	Result int `json:"result"`
	System struct {
		NewVersion string `json:"new_version"`
		UID        int    `json:"uid"`
		AuthStatus int    `json:"auth_status"`
		SessionID  string `json:"session_id"`
		NeedLogin  int    `json:"need_login"`
		Email      string `json:"email"`
		TraceID    string `json:"trace_id"`
		Redirect   int    `json:"redirect"`
		Controller string `json:"controller"`
		TokenID    string `json:"token_id"`
		TimeStamp  int    `json:"time_stamp"`
	} `json:"system"`
	Data struct {
		Result     int         `json:"result"`
		ResultRows []orderInfo `json:"result_rows"`
		ResInfo    string      `json:"res_info"`
	} `json:"data"`
	ResInfo string `json:"res_info"`
}

type orderInfo struct {
	TemplateContent []struct {
		MerchInfo struct {
			MerchID   string `json:"merch_id"`
			MerchName string `json:"merch_name"`
			URL       string `json:"url"`
		} `json:"merch_info,omitempty"`
		StateInfo struct {
			StateDesc string `json:"state_desc"`
		} `json:"state_info,omitempty"`
		Type           int `json:"type"`
		OrderGoodsInfo struct {
			Title struct {
				Title     string `json:"title"`
				FontColor string `json:"font_color"`
			} `json:"title"`
			LabelList []interface{} `json:"label_list"`
			GoodsInfo struct {
				SkuPic      string `json:"sku_pic"`
				SkuID       string `json:"sku_id"`
				ProductInfo string `json:"product_info"`
			} `json:"goods_info"`
			TitleLabelList []interface{} `json:"title_label_list"`
		} `json:"order_goods_info,omitempty"`
		OrderCapitalInfo []struct {
			Content     string `json:"content"`
			TotalAmount string `json:"total_amount"`
			FontSize    string `json:"font_size"`
			Firstpay    string `json:"firstpay"`
			FontColor   string `json:"font_color"`
			HandlingFee string `json:"handling_fee"`
		} `json:"order_capital_info,omitempty"`
		ButtonList []struct {
			ButtonKey    string `json:"button_key"`
			Text         string `json:"text"`
			Sort         int    `json:"sort"`
			BgColor      string `json:"bg_color"`
			BgPressColor string `json:"bg_press_color"`
			FontColor    string `json:"font_color"`
			Key          string `json:"key"`
			URL          string `json:"url"`
		} `json:"button_list,omitempty"`
	} `json:"template_content"`
	OrderInfo struct {
		DetailURL        string `json:"detail_url"`
		CheckRequire     int    `json:"check_require"`
		OrderState       int    `json:"order_state"`
		PayWay           int    `json:"pay_way"`
		ParentPayOrderID string `json:"parent_pay_order_id"`
		ParentOrderID    string `json:"parent_order_id"`
		CreateTime       int64  `json:"create_time"`
		AppSourceType    int    `json:"app_source_type"`
		RepayType        int    `json:"repay_type"`
		OrderID          string `json:"order_id"`
		SaleType         int    `json:"sale_type"`
		OrderType        int    `json:"order_type"`
	} `json:"order_info"`
}

type OrderFullInfoResp struct {
	Retcode     int    `json:"retcode"`
	Retmsg      string `json:"retmsg"`
	VirtualInfo struct {
		FuluInfo []struct {
			CardNumber struct {
				Title string `json:"title"`
				Value string `json:"value"`
			} `json:"card_number"`
			Passwd struct {
				Title string `json:"title"`
				Value string `json:"value"`
			} `json:"passwd"`
		} `json:"fulu_info"`
		CouponInfo []interface{} `json:"coupon_info"`
	} `json:"virtual_info"`
}
