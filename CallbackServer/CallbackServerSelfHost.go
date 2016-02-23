package main

import (
	"encoding/json"
	"fmt"
	"github.com/DuoSoftware/gorest"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

type CallbackServerSelfHost struct {
	gorest.RestService `root:"/CallbackServerSelfHost/" consumes:"application/json" produces:"application/json"`
	addCallback        gorest.EndPoint `method:"POST" path:"/Callback/AddCallback/" postdata:"CampaignCallback"`
}

func (callbackServerSelfHost CallbackServerSelfHost) AddCallback(callbackInfo CampaignCallback) {
	user := context.Get(callbackServerSelfHost.Context.Request(), "user")
	if user != nil {
		iTenant := user.(*jwt.Token).Claims["tenant"]
		iCompany := user.(*jwt.Token).Claims["company"]
		if iTenant != nil && iCompany != nil {
			tenant := int(iTenant.(float64))
			company := int(iCompany.(float64))
			authHeaderStr := fmt.Sprintf("%d#%d", tenant, company)
			fmt.Println("Start AddCallback: ", callbackInfo.CallbackUrl, "#", callbackInfo.DialoutTime.String())
			fmt.Println(authHeaderStr)

			ch := make(chan error)
			go AddCallbackInfoToRedis(company, tenant, callbackInfo, ch)
			var err = <-ch
			close(ch)
			fmt.Println(err.Error())
			if err != nil {
				callbackServerSelfHost.RB().Write(ResponseGenerator(false, "Add callback info failed", "", err.Error()))
				return
			} else {
				callbackServerSelfHost.RB().Write(ResponseGenerator(true, "Add callback info success", "", ""))
				return
				if callbackInfo.Class == "DIALER" && callbackInfo.Type == "CALLBACK" && callbackInfo.Category == "INTERNAL" {
					go UploadCampaignMgrCallbackInfo(company, tenant, callbackInfo.CampaignId, callbackInfo.CallbackObj)
				}
			}
		} else {
			callbackServerSelfHost.RB().Write(ResponseGenerator(false, "Invalid company or tenant", "", ""))
			return
		}
	} else {
		callbackServerSelfHost.RB().Write(ResponseGenerator(false, "User data not found in JWT", "", ""))
		return
	}
}

func ResponseGenerator(isSuccess bool, customMessage, result, exception string) []byte {
	res := Result{}
	res.IsSuccess = isSuccess
	res.CustomMessage = customMessage
	res.Exception = exception
	res.Result = result
	resb, _ := json.Marshal(res)
	return resb
}
