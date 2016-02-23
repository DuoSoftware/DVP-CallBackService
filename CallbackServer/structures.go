package main

import (
	"time"
)

type Configuration struct {
	SecurityIp                       string
	SecurityPort                     string
	RedisIp                          string
	RedisPort                        string
	RedisDb                          int
	CallbackServerId                 string
	HostIpAddress                    string
	Port                             string
	ExternalCallbackRequestFrequency time.Duration
	CampaignServiceHost              string
	CampaignServicePort              string
}

type EnvConfiguration struct {
	SecurityIp                       string
	SecurityPort                     string
	RedisIp                          string
	RedisPort                        string
	RedisDb                          string
	CallbackServerId                 string
	HostIpAddress                    string
	Port                             string
	ExternalCallbackRequestFrequency string
	CampaignServiceHost              string
	CampaignServicePort              string
}

type CampaignCallback struct {
	Company     int
	Tenant      int
	Class       string
	Type        string
	Category    string
	DialoutTime time.Time
	CallbackUrl string
	CallbackObj string
	CampaignId  string
}

type Result struct {
	Exception     string
	CustomMessage string
	IsSuccess     bool
	Result        string
}
