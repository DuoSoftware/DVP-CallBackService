package main

import (
	"encoding/json"
	"fmt"
	"github.com/DuoSoftware/gorest"
	"time"
)

type CallbackServerSelfHost struct {
	gorest.RestService    `root:"/DVP/API/" consumes:"application/json" produces:"application/json"`
	addCallback           gorest.EndPoint `method:"POST" path:"/{Version:string}/Callback/AddCallback/" postdata:"CampaignCallback"`
	addCallbackByDuration gorest.EndPoint `method:"POST" path:"/{Version:string}/Callback/AddCallbackByDuration/" postdata:"CampaignCallback"`
}

func (callbackServerSelfHost CallbackServerSelfHost) AddCallback(callbackInfo CampaignCallback, Version string) {
	company, tenant, _, msg := decodeJwtDialer(callbackServerSelfHost, "dialer", "read")
	if company != 0 && tenant != 0 {
		authHeaderStr := fmt.Sprintf("%d#%d", tenant, company)
		fmt.Println("Start AddCallback: ", callbackInfo.CallbackUrl, "#", callbackInfo.DialoutTime.String())
		fmt.Println(authHeaderStr)

		ch := make(chan error)
		go AddCallbackInfoToRedis(company, tenant, callbackInfo, ch)
		if <-ch != nil {
			var err = <-ch
			fmt.Println(err.Error())
			close(ch)

			errStr, _ := json.Marshal(ResponseGenerator(false, "Add callback info failed", "", err.Error()))

			callbackServerSelfHost.RB().Write(errStr)
		} else {
			close(ch)

			resStr, _ := json.Marshal(ResponseGenerator(true, "Add callback info success", "", ""))
			callbackServerSelfHost.RB().Write(resStr)
			if callbackInfo.Class == "DIALER" && callbackInfo.Type == "CALLBACK" && callbackInfo.Category == "INTERNAL" {
				go UploadCampaignMgrCallbackInfo(company, tenant, callbackInfo.CampaignId, callbackInfo.CallbackObj)
			}
		}

		return
	} else {
		defStr, _ := json.Marshal(msg)
		callbackServerSelfHost.RB().Write(defStr)
		return
	}
}

func (callbackServerSelfHost CallbackServerSelfHost) AddCallbackByDuration(callbackInfo CampaignCallback, Version string) {
	company, tenant, _, msg := decodeJwtDialer(callbackServerSelfHost, "dialer", "read")
	if company != 0 && tenant != 0 {
		authHeaderStr := fmt.Sprintf("%d#%d", tenant, company)

		tmNow := time.Now().UTC()
		secCount := tmNow.Second() + callbackInfo.CallbackDuration
		callbackTime := time.Date(tmNow.Year(), tmNow.Month(), tmNow.Day(), tmNow.Hour(), tmNow.Minute(), secCount, 0, time.UTC)
		fmt.Println("callbackTime:: ", callbackTime)

		callbackInfo.Company = company
		callbackInfo.Tenant = tenant
		callbackInfo.DialoutTime = callbackTime
		callbackInfo.CallbackUrl = fmt.Sprintf("http://%s/DVP/DialerAPI/ResumeCallback", CreateHost(dialerServiceHost, dialerServicePort))

		fmt.Println("Start AddCallback: ", callbackInfo.CallbackUrl, "#", callbackInfo.DialoutTime.String())
		fmt.Println(authHeaderStr)

		ch := make(chan error)
		go AddCallbackInfoToRedis(company, tenant, callbackInfo, ch)
		if <-ch != nil {
			var err = <-ch
			fmt.Println(err.Error())
			close(ch)

			errStr, _ := json.Marshal(ResponseGenerator(false, "Add callback info failed", "", err.Error()))

			callbackServerSelfHost.RB().Write(errStr)
		} else {
			close(ch)

			resStr, _ := json.Marshal(ResponseGenerator(true, "Add callback info success", "", ""))
			callbackServerSelfHost.RB().Write(resStr)
			if callbackInfo.Class == "DIALER" && callbackInfo.Type == "CALLBACK" && callbackInfo.Category == "INTERNAL" {
				go UploadCampaignMgrCallbackInfo(company, tenant, callbackInfo.CampaignId, callbackInfo.CallbackObj)
			}
		}

		return
	} else {
		defStr, _ := json.Marshal(msg)
		callbackServerSelfHost.RB().Write(defStr)
		return
	}
}
