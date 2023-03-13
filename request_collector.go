package main

import (
	"github.com/gofiber/fiber/v2"
	"systest/log"
	"systest/payloads"
	"time"
)

var requestLogger = log.New("Request Collector")
var startTime int64
var currentRequest *payloads.HttpRequest

func getUnixTime() int64 {
	return time.Now().UnixMilli()
}

func StartCollectingRequests() {
	startTime = getUnixTime()
	requestLogger.Info("Request collection started.")
}

func StopCollectingRequests() *payloads.HttpRequest {
	requestLogger.Info("Request collection completed after %dms.", getUnixTime()-startTime)
	startTime = 0

	response := currentRequest
	currentRequest = nil
	return response
}

func IsCollectingRequests() bool {
	return startTime != 0
}

func CollectRequest(instance *payloads.EurekaInstance, ctx *fiber.Ctx) {
	if startTime == 0 {
		return
	}

	instanceId := instance.InstanceId
	if instanceId == "" {
		instanceId = instance.HostName
	}

	requestLogger.Info("Collecting request to %s...", instanceId)

	request := &payloads.HttpRequest{
		App:        instance.App,
		InstanceId: instanceId,
		Method:     ctx.Method(),
		Path:       ctx.Path(),
		Headers:    ctx.GetReqHeaders(),
		Body:       string(ctx.Body()),
		Time:       getUnixTime(),
		Parent:     currentRequest,
		Children:   make([]*payloads.HttpRequest, 0),
	}

	if currentRequest == nil {
		currentRequest = request
	} else {
		currentRequest.Children = append(currentRequest.Children, request)
		currentRequest = request
	}
}

func CollectResponse(ctx *fiber.Ctx) {
	if startTime == 0 {
		return
	}

	requestLogger.Info("Collecting response...")

	if currentRequest == nil {
		requestLogger.Warn("Attempted to collect a response before a request was collected.")
	} else if currentRequest.Response != nil {
		requestLogger.Warn("A response was already collected for the current request, this will be overridden.")
	}

	currentRequest.Response = &payloads.HttpResponse{
		Status:  ctx.Response().StatusCode(),
		Headers: ctx.GetRespHeaders(),
		Body:    string(ctx.Response().Body()),
		Time:    getUnixTime(),
	}

	if currentRequest.Parent != nil {
		currentRequest = currentRequest.Parent
	}
}
