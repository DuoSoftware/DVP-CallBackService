package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
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
	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("companyinfo", authToken)
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

func UploadScheduledCallbackInfo(tenant, company int, callbackData string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UploadScheduledCallbackInfo", r)
		}
	}()
	fmt.Println("request:", callbackData)

	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/ScheduledCallback", CreateHost(campaignServiceHost, campaignServicePort))
	authToken := fmt.Sprintf("%d:%d", tenant, company)
	req, err := http.NewRequest("POST", serviceurl, bytes.NewBufferString(callbackData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("companyinfo", authToken)
	fmt.Println("request:", serviceurl)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
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

func UploadDispatchedTime(tenant, company int, sessionId, dispatchedTime string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in UploadDispatchedTime", r)
		}
	}()

	serviceurl := fmt.Sprintf("http://%s/DVP/API/1.0.0.0/CampaignManager/ScheduledCallback/%s/Dispatched/%s", CreateHost(campaignServiceHost, campaignServicePort), sessionId, dispatchedTime)
	authToken := fmt.Sprintf("%d:%d", tenant, company)
	req, err := http.NewRequest("PUT", serviceurl, bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "bearer "+accessToken)
	req.Header.Set("companyinfo", authToken)
	fmt.Println("request:", serviceurl)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer resp.Body.Close()
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
