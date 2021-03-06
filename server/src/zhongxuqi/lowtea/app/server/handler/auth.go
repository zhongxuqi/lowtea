package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"zhongxuqi/lowtea/errors"
	"zhongxuqi/lowtea/model"
	"zhongxuqi/lowtea/utils"
)

// CheckAdmin
func (p *MainHandler) CheckAdmin(r *http.Request) (err error) {
	var accountCookie *http.Cookie
	accountCookie, err = r.Cookie("account")
	if err != nil {
		return
	}

	// check root
	if accountCookie.Value == model.ROOT {
		return
	}

	// check admin
	var user model.User
	user, err = p.UserModel.FindByAccount(accountCookie.Value)
	if err != nil {
		return
	}
	if user.Role != model.ADMIN {
		err = errors.ERROR_PERMISSION_DENIED
	}
	return
}

// CheckRoot
func (p *MainHandler) CheckRoot(r *http.Request) (err error) {
	var accountCookie *http.Cookie
	accountCookie, err = r.Cookie("account")
	if err != nil {
		return
	}

	// check root
	if accountCookie.Value != model.ROOT {
		err = errors.ERROR_PERMISSION_DENIED
	}
	return
}

// Login do login
func (p *MainHandler) Login(w http.ResponseWriter, r *http.Request) {
	var dataStruct struct {
		Account    string `json:"account"`
		ExpireTime int64  `json:"expireTime"`
		Sign       string `json:"sign"`
	}
	err := utils.ReadReq2Struct(r, &dataStruct)
	if err != nil {
		http.Error(w, "[Login] data read and unmarshal fail: "+err.Error(), 400)
		return
	}

	if dataStruct.ExpireTime < time.Now().Unix() {
		http.Error(w, "[Login] "+ERROR_TIME_EXPIRED.Error(), 400)
		return
	}

	err = p.checkSign(dataStruct.Account, dataStruct.ExpireTime, dataStruct.Sign)
	if err != nil {
		http.Error(w, "check sign fail: "+err.Error(), 400)
		return
	}

	// if pass check sign, update session
	p.UpdateSession(w, dataStruct.Account)

	var user model.User
	user, err = p.UserModel.FindByAccount(dataStruct.Account)
	if err != nil {
		http.Error(w, "find User error: "+err.Error(), 400)
		return
	}

	var ret struct {
		model.RespBase
		User *model.User `json:"user"`
	}
	ret.Status = 200
	ret.Message = "success"
	ret.User = &user

	retStr, _ := json.Marshal(ret)
	w.Write(retStr)
}

// Register do register
func (p *MainHandler) Register(w http.ResponseWriter, r *http.Request) {
	var register model.Register
	err := utils.ReadReq2Struct(r, &register)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	// check the account
	if register.Account == model.ROOT {
		http.Error(w, errors.ERROR_USERNAME_EXISTS.Error(), 400)
		return
	}
	n, err := p.UserModel.CountByAccount(register.Account)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	} else if n > 0 {
		http.Error(w, errors.ERROR_USERNAME_EXISTS.Error(), 400)
		return
	}
	n, err = p.RegisterModel.CountByAccount(register.Account)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	} else if n > 0 {
		http.Error(w, errors.ERROR_USERNAME_EXISTS.Error(), 400)
		return
	}

	err = p.RegisterModel.Insert(register)
	if err != nil {
		http.Error(w, "[Register] insert error: "+err.Error(), 400)
		return
	}

	retStr, _ := json.Marshal(&model.RespBase{
		Status:  200,
		Message: "success",
	})
	w.Write(retStr)
}

// Logout do logout
func (p *MainHandler) Logout(w http.ResponseWriter, r *http.Request) {
	p.ClearSession(w)
	retStr, _ := json.Marshal(&model.RespBase{
		Status:  200,
		Message: "success",
	})
	w.Write(retStr)
}
