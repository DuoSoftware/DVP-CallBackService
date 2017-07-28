package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

//----------Campaign Manager Service-----------------------
func UploadCampaignMgrCallbackInfo(company, tenant int, campaignId, callback string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UploadCallbackInfo", r)
		}
	}()
	fmt.Println("request:", callback)

	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/Campaign/%s/Callback", CreateHost(campaignServiceHost, campaignServicePort), campaignId)
	authToken := fmt.Sprintf("%d:%d", tenant, company)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBufferString(callback))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	fmt.Println("request:", serviceurl)
	client := &http.Client{}
	fmt.Println("-------------------------")
	resp, err := client.Do(req)
	fmt.Println("+++++++++++++++++++++++++")
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("==========================")
	defer resp.Body.Close()
	fmt.Println("]]]]]]]]]]]]]]]]]]]]]]]]]]]")
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, errb := ioutil.ReadAll(resp.Body)
	if errb != nil {
		fmt.Println(err.Error())
	} else {
		result := string(body)
		fmt.Println("response Body:", result)
	}
}

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
	tmNowUtc := time.Now().UTC().Unix()
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
			go SendCallback(campCallback.Company, campCallback.Tenant, campCallback.CallbackUrl, campCallback.CallbackObj)
		}
	}
}

func SendCallback(company, tenant int, callbackUrl, callbackObj string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UploadCallbackInfo", r)
		}
	}()
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
}
