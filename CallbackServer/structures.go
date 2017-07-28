package main

import (
	"time"
)

type Configuration struct {
	SecurityIp                       string
	SecurityPort                     string
	RedisIp                          string
	RedisPort                        string
	RedisPassword                    string
	RedisDb                          int
	CallbackServerId                 string
	HostIpAddress                    string
	Port                             string
	ExternalCallbackRequestFrequency string
	CampaignServiceHost              string
	CampaignServicePort              string
	DialerServiceHost                string
	DialerServicePort                string
	AccessToken                      string
	RedisMode                        string
	RedisClusterName                 string
	SentinelHosts                    string
	SentinelPort                     string
}

type EnvConfiguration struct {
	SecurityIp                       string
	SecurityPort                     string
	RedisIp                          string
	RedisPort                        string
	RedisPassword                    string
	RedisDb                          string
	CallbackServerId                 string
	HostIpAddress                    string
	Port                             string
	ExternalCallbackRequestFrequency string
	CampaignServiceHost              string
	CampaignServicePort              string
	DialerServiceHost                string
	DialerServicePort                string
	AccessToken                      string
	RedisMode                        string
	RedisClusterName                 string
	SentinelHosts                    string
	SentinelPort                     string
}

type ScheduledCallbackInfo struct {
	Class          string
	Type           string
	Category       string
	CompanyId      int
	TenantId       int
	SessionId      string
	ContactId      string
	CallbackData   string
	RequestedTime  string
	Duration       int
	DispatchedTime string
}

type CampaignCallback struct {
	Company          int
	Tenant           int
	Class            string
	Type             string
	Category         string
	SessionId        string
	DialoutTime      time.Time
	CallbackDuration int
	CallbackUrl      string
	CallbackObj      string
	CampaignId       string
}

type Result struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        string
}
