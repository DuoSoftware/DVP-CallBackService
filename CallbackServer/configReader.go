package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var dirPath string
var securityIp string
var securityPort string
var redisPassword string
var redisIp string
var redisPort string
var redisDb int
var callbackServerId string
var hostIpAddress string
var port string
var externalCallbackRequestFrequency time.Duration
var campaignServiceHost string
var campaignServicePort string
var dialerServiceHost string
var dialerServicePort string
var accessToken string
var redisMode string
var redisClusterName string
var sentinelHosts string
var sentinelPort string

func GetDirPath() string {
	envPath := os.Getenv("GO_CONFIG_DIR")
	if envPath == "" {
		envPath = "./"
	}
	fmt.Println(envPath)
	return envPath
}

func GetDefaultConfig() Configuration {
	confPath := filepath.Join(dirPath, "conf.json")
	fmt.Println("GetDefaultConfig config path: ", confPath)
	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	defconfiguration := Configuration{}
	deferr := json.Unmarshal(content, &defconfiguration)

	if deferr != nil {
		fmt.Println("error:", deferr)
		defconfiguration.SecurityIp = "127.0.0.1"
		defconfiguration.SecurityPort = "6379"
		defconfiguration.RedisIp = "127.0.0.1"
		defconfiguration.RedisPort = "6379"
		defconfiguration.RedisPassword = "DuoS123"
		defconfiguration.RedisDb = 6
		defconfiguration.CallbackServerId = "1"
		defconfiguration.HostIpAddress = "127.0.0.1"
		defconfiguration.Port = "2226"
		defconfiguration.ExternalCallbackRequestFrequency = 300
		defconfiguration.CampaignServiceHost = "127.0.0.1"
		defconfiguration.CampaignServicePort = "8827"
		defconfiguration.DialerServiceHost = "127.0.0.1"
		defconfiguration.DialerServicePort = "2226"
		defconfiguration.AccessToken = ""
		defconfiguration.RedisMode = "instance"
		//instance, cluster, sentinel
		defconfiguration.RedisClusterName = "redis-cluster"
		defconfiguration.SentinelHosts = "138.197.90.92,45.55.205.92,138.197.90.92"
		defconfiguration.SentinelPort = "16389"
	}

	return defconfiguration
}

func LoadDefaultConfig() {
	defconfiguration := GetDefaultConfig()
	securityIp = fmt.Sprintf("%s:%s", defconfiguration.RedisIp, defconfiguration.RedisPort)
	securityPort = defconfiguration.RedisPort
	redisPassword = defconfiguration.RedisPassword
	redisIp = fmt.Sprintf("%s:%s", defconfiguration.RedisIp, defconfiguration.RedisPort)
	redisPort = defconfiguration.RedisPort
	redisDb = defconfiguration.RedisDb
	callbackServerId = defconfiguration.CallbackServerId
	hostIpAddress = defconfiguration.HostIpAddress
	port = defconfiguration.Port
	externalCallbackRequestFrequency = defconfiguration.ExternalCallbackRequestFrequency
	campaignServiceHost = defconfiguration.CampaignServiceHost
	campaignServicePort = defconfiguration.CampaignServicePort
	dialerServiceHost = defconfiguration.DialerServiceHost
	dialerServicePort = defconfiguration.DialerServicePort
	accessToken = defconfiguration.AccessToken
	redisMode = defconfiguration.RedisMode
	redisClusterName = defconfiguration.RedisClusterName
	sentinelHosts = defconfiguration.SentinelHosts
	sentinelPort = defconfiguration.SentinelPort
}

func LoadConfiguration() {
	dirPath = GetDirPath()
	confPath := filepath.Join(dirPath, "custom-environment-variables.json")
	fmt.Println("InitiateRedis config path: ", confPath)

	content, operr := ioutil.ReadFile(confPath)
	if operr != nil {
		fmt.Println(operr)
	}

	envconfiguration := EnvConfiguration{}
	enverr := json.Unmarshal(content, &envconfiguration)
	if enverr != nil {
		fmt.Println("error:", enverr)
		LoadDefaultConfig()
	} else {
		var converr error
		defConfig := GetDefaultConfig()
		securityIp = os.Getenv(envconfiguration.SecurityIp)
		securityPort = os.Getenv(envconfiguration.SecurityPort)
		redisPassword = os.Getenv(envconfiguration.RedisPassword)
		redisIp = os.Getenv(envconfiguration.RedisIp)
		redisPort = os.Getenv(envconfiguration.RedisPort)
		redisDb, converr = strconv.Atoi(os.Getenv(envconfiguration.RedisDb))
		callbackServerId = os.Getenv(envconfiguration.CallbackServerId)
		hostIpAddress = os.Getenv(envconfiguration.HostIpAddress)
		port = os.Getenv(envconfiguration.Port)
		externalCallbackRequestFrequencyTemp := os.Getenv(envconfiguration.ExternalCallbackRequestFrequency)
		campaignServiceHost = os.Getenv(envconfiguration.CampaignServiceHost)
		campaignServicePort = os.Getenv(envconfiguration.CampaignServicePort)
		dialerServiceHost = os.Getenv(envconfiguration.DialerServiceHost)
		dialerServicePort = os.Getenv(envconfiguration.DialerServicePort)
		accessToken = os.Getenv(envconfiguration.AccessToken)
		redisMode = os.Getenv(envconfiguration.RedisMode)
		redisClusterName = os.Getenv(envconfiguration.RedisClusterName)
		sentinelHosts = os.Getenv(envconfiguration.SentinelHosts)
		sentinelPort = os.Getenv(envconfiguration.SentinelPort)

		if securityIp == "" {
			securityIp = defConfig.SecurityIp
		}
		if securityPort == "" {
			securityPort = defConfig.SecurityPort
		}
		if redisIp == "" {
			redisIp = defConfig.RedisIp
		}
		if redisPort == "" {
			redisPort = defConfig.RedisPort
		}
		if redisPassword == "" {
			redisPassword = defConfig.RedisPassword
		}
		if redisDb == 0 || converr != nil {
			redisDb = defConfig.RedisDb
		}
		if callbackServerId == "" {
			callbackServerId = defConfig.CallbackServerId
		}
		if hostIpAddress == "" {
			hostIpAddress = defConfig.HostIpAddress
		}
		if port == "" {
			port = defConfig.Port
		}
		if externalCallbackRequestFrequencyTemp == "" {
			externalCallbackRequestFrequency = defConfig.ExternalCallbackRequestFrequency
		} else {
			externalCallbackRequestFrequency, _ = time.ParseDuration(externalCallbackRequestFrequencyTemp + "s")
		}
		if campaignServiceHost == "" {
			campaignServiceHost = defConfig.CampaignServiceHost
		}
		if campaignServicePort == "" {
			campaignServicePort = defConfig.CampaignServicePort
		}
		if dialerServiceHost == "" {
			dialerServiceHost = defConfig.DialerServiceHost
		}
		if dialerServicePort == "" {
			dialerServicePort = defConfig.DialerServicePort
		}
		if accessToken == "" {
			accessToken = defConfig.AccessToken
		}
		if redisMode == "" {
			redisMode = defConfig.RedisMode
		}
		if redisClusterName == "" {
			redisClusterName = defConfig.RedisClusterName
		}
		if sentinelHosts == "" {
			sentinelHosts = defConfig.SentinelHosts
		}
		if sentinelPort == "" {
			sentinelPort = defConfig.SentinelPort
		}

		redisIp = fmt.Sprintf("%s:%s", redisIp, redisPort)
		securityIp = fmt.Sprintf("%s:%s", securityIp, securityPort)
	}

	fmt.Println("RedisMode:", redisMode)
	fmt.Println("redisIp:", redisIp)
	fmt.Println("redisDb:", redisDb)
	fmt.Println("SentinelHosts:", sentinelHosts)
	fmt.Println("SentinelPort:", sentinelPort)
}
