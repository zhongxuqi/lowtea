package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"zhongxuqi/lowtea/errors"
	"zhongxuqi/lowtea/model"

	"gopkg.in/mgo.v2/bson"
)

func (p *MainHandler) ActionPublicDocuments(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		var err error
		err = r.ParseForm()
		if err != nil {
			http.Error(w, "parse url form error: "+err.Error(), 400)
			return
		}

		var params struct {
			PageSize  int    `json:"pageSize"`
			PageIndex int    `json:"pageIndex"`
			Keyword   string `json:"keyword"`
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
		params.Keyword = r.Form.Get("keyword")

		var respBody struct {
			model.RespBase
			Documents []model.Document `json:"documents"`
			PageTotal int              `json:"pageTotal"`
			DocTotal  int              `json:"docTotal"`
		}
		filter := bson.M{
			"status": model.STATUS_PUBLISH_PUBLIC,
		}
		if params.Keyword != "" {
			filter["title"] = &bson.M{
				"$regex": params.Keyword,
			}
		}

		var n int
		n, err = p.DocumentModel.CountByFilter(filter)
		if err != nil {
			http.Error(w, "find document count error: "+err.Error(), 500)
			return
		}
		respBody.DocTotal = n
		if n > 0 {
			respBody.PageTotal = (n-1)/params.PageSize + 1
		} else {
			respBody.PageTotal = 0
		}

		respBody.Documents, err = p.DocumentModel.SortFindByFilterWithPage(filter, "-modifyTime", params.PageSize*params.PageIndex, params.PageSize)
		if err != nil {
			http.Error(w, "find documents error: "+err.Error(), 500)
			return
		}

		for i, _ := range respBody.Documents {
			n, _ = p.StarModel.CountByDocumentId(respBody.Documents[i].Id.Hex())
			respBody.Documents[i].StarNum = n
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

func (p *MainHandler) ActionPublicDocument(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cmds := strings.Split(r.URL.Path, "/")
		if len(cmds) < 3 {
			http.Error(w, errors.ERROR_EMPTY_ID.Error(), 400)
			return
		}
		documentId := cmds[3]
		if !bson.IsObjectIdHex(documentId) {
			http.Error(w, errors.ERROR_INVAIL_ID.Error(), 400)
			return
		}
		var respBody struct {
			model.RespBase
			Document model.Document `json:"document"`
		}
		var err error
		respBody.Document, err = p.DocumentModel.FindDocument(bson.ObjectIdHex(documentId))
		if err != nil {
			http.Error(w, "find document error: "+err.Error(), 500)
			return
		}

		if respBody.Document.Status != model.STATUS_PUBLISH_PUBLIC {
			http.Error(w, errors.ERROR_PERMISSION_DENIED.Error(), 400)
			return
		}

		var n int
		n, _ = p.StarModel.CountByDocumentId(documentId)
		respBody.Document.StarNum = n

		respBody.Status = 200
		respBody.Message = "success"
		respByte, _ := json.Marshal(&respBody)
		w.Write(respByte)
		return
	}
	http.Error(w, "Not Found", 404)
	return
}
