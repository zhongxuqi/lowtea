package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"zhongxuqi/lowtea/model"

	"gopkg.in/mgo.v2/bson"
)

func (p *MainHandler) ActionDrafts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var accountCookie *http.Cookie
		var err error
		accountCookie, err = r.Cookie("account")
		if err != nil {
			http.Error(w, "cookie find error: "+err.Error(), 400)
			return
		}

		err = r.ParseForm()
		if err != nil {
			http.Error(w, "parse url form error: "+err.Error(), 400)
			return
		}

		var params struct {
			PageSize  int `json:"pageSize"`
			PageIndex int `json:"pageIndex"`
		}
		params.PageSize, err = strconv.Atoi(r.Form.Get("pageSize"))
		if err != nil {
			http.Error(w, "read param pageSize error: "+err.Error(), 400)
			return
		}
		params.PageIndex, err = strconv.Atoi(r.Form.Get("pageIndex"))
		if err != nil {
			http.Error(w, "read param pageIndex error: "+err.Error(), 400)
			return
		}

		var respBody struct {
			model.RespBase
			Documents []model.Document `json:"documents"`
			PageTotal int              `json:"pageTotal"`
		}
		filter := bson.M{
			"account": accountCookie.Value,
			"status":  model.STATUS_DRAFT,
		}

		var n int
		n, err = p.DocumentColl.Find(filter).Count()
		if err != nil {
			http.Error(w, "find document count error: "+err.Error(), 500)
			return
		}
		if n > 0 {
			respBody.PageTotal = (n-1)/params.PageSize + 1
		} else {
			respBody.PageTotal = 0
		}

		err = p.DocumentColl.Find(&filter).Sort("-createTime").Skip(params.PageSize * params.PageIndex).Limit(params.PageSize).All(&respBody.Documents)
		if err != nil {
			http.Error(w, "find documents error: "+err.Error(), 500)
			return
		}

		respBody.Status = 200
		respBody.Message = "success"
		respByte, _ := json.Marshal(respBody)
		w.Write(respByte)
		return
	}
	http.Error(w, "Not Found", 404)
	return

}