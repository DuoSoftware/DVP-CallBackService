package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//----------Callbak Info-----------------------
func AddCallbackInfoToRedis(company, tenant int, callback CampaignCallback, aci chan error) {
	callback.Company = company
	callback.Tenant = tenant

	callbackKey := fmt.Sprintf("CallbackInfo:%s:%d:%d", callbackServerId, company, tenant)
	score := float64(callback.DialoutTime.UTC().Unix())
	jsonData, _ := json.Marshal(callback)
	_, err := RedisZadd(callbackKey, string(jsonData), score)
	if err != nil {
		aci <- err
	} else {
		aci <- nil
	}
}

func SetLastExecuteTime(executeTime string) string {
	key := fmt.Sprintf("CallbackServerLastExecuteTime:%s", callbackServerId)
	lastExeTimeStr := RedisGet(key)
	RedisSet(key, executeTime)
	if lastExeTimeStr == "" {
		return "0"
	} else {
		return lastExeTimeStr
	}
}

func ExecuteCallback() {
	tmNow := time.Now().UTC()
	tmNowUtc := tmNow.Unix()
	tmNowUtcStr := strconv.FormatFloat(float64(tmNowUtc), 'E', -1, 64)
	lastExeTimeStr := fmt.Sprintf("(%s", SetLastExecuteTime(tmNowUtcStr))
	fmt.Println("tmNowUtcStr: ", tmNowUtcStr)
	fmt.Println("lastExeTimeStr: ", lastExeTimeStr)
	callbackListSearchKey := fmt.Sprintf("CallbackInfo:%s:*", callbackServerId)
	AllCallbackList := RedisSearchKeys(callbackListSearchKey)
	for _, callbackList := range AllCallbackList {
		fmt.Println("Execute callback list: ", callbackList)
		campaignCallbacks := RedisZRangeByScore(callbackList, lastExeTimeStr, tmNowUtcStr)
		for _, cmpCallbackStr := range campaignCallbacks {
			fmt.Println("cmpCallbackStr: ", cmpCallbackStr)

			RedisZRemove(callbackList, cmpCallbackStr)
			var campCallback CampaignCallback
			json.Unmarshal([]byte(cmpCallbackStr), &campCallback)
			go SendCallback(campCallback, tmNow.String())

		}
	}
}

func SendCallback(campCallback CampaignCallback, tmNow string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UploadCallbackInfo", r)
		}
	}()

	company := campCallback.Company
	tenant := campCallback.Tenant
	callbackUrl := campCallback.CallbackUrl
	callbackObj := campCallback.CallbackObj

	fmt.Println("request:", callbackUrl)
	authToken := fmt.Sprintf("bearer %s", accessToken)
	internalAccess := fmt.Sprintf("%d:%d", tenant, company)
	req, err := http.NewRequest("POST", callbackUrl, bytes.NewBufferString(callbackObj))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	req.Header.Set("companyinfo", internalAccess)
	fmt.Println("request:", callbackObj)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	if resp.StatusCode == 200 {
		go UploadDispatchedTime(campCallback.Tenant, campCallback.Company, campCallback.SessionId, tmNow)
	}
}
