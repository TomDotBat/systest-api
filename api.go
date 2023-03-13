package main

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"os"
	"strconv"
	"strings"
	"systest/log"
	"systest/payloads"
)

func getAgent(method string, url string) (*fiber.Agent, error) {
	switch strings.ToUpper(method) {
	case "GET":
		return fiber.Get(url), nil
	case "POST":
		return fiber.Post(url), nil
	case "PUT":
		return fiber.Put(url), nil
	case "PATCH":
		return fiber.Patch(url), nil
	case "DELETE":
		return fiber.Delete(url), nil
	case "HEAD":
		return fiber.Head(url), nil
	}
	return nil, errors.New("invalid method: " + method)
}

func RegisterApiRoutes(app *fiber.App) {
	logger := log.New("API")
	validate := validator.New()

	ipAddress := os.Getenv("SERVER_IP_ADDRESS")

	app.Post("/request", func(ctx *fiber.Ctx) error {
		logger.Info("Received SysTest request...")

		if IsCollectingRequests() {
			logger.Warn("The request collector is already active.")
			return ctx.Status(fiber.StatusConflict).JSON(payloads.ErrorResponse{
				Status:  fiber.StatusConflict,
				Message: "The request collector is already active",
			})
		}

		payload := new(payloads.SysTestRequest)
		if err := ctx.BodyParser(payload); err != nil {
			logger.Warn("Failed to parse request body: %s", err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(payloads.ErrorResponse{
				Status:  fiber.StatusBadRequest,
				Message: "The request body could not be parsed",
			})
		}

		if err := validate.Struct(payload); err != nil {
			logger.Warn("Failed to validate request body:\n%s", err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(payloads.ErrorResponse{
				Status:  fiber.StatusBadRequest,
				Message: "The request body is invalid",
			})
		}

		logger.Info("Determining where the request should be forwarded...")

		var port int
		if port = GetPortByInstanceId(payload.InstanceId); port == -1 {
			port = GetPortByAppAndHostname(payload.App, payload.InstanceId)
		}

		if port == -1 {
			logger.Warn("Could not locate a proxy for instance: %s", payload.InstanceId)
			return ctx.Status(fiber.StatusNotFound).JSON(payloads.ErrorResponse{
				Status:  fiber.StatusNotFound,
				Message: "A proxy does not exist for the specified instance",
			})
		}

		logger.Info("Forwarding request to %s:%d...", ipAddress, port)

		agent, err := getAgent(payload.Method, "http://"+ipAddress+":"+strconv.Itoa(port)+payload.Path)
		if err != nil {
			logger.Warn("Failed to create agent: %s", err.Error())
			return ctx.Status(fiber.StatusBadRequest).JSON(payloads.ErrorResponse{
				Status:  fiber.StatusBadRequest,
				Message: "Unsupported HTTP method",
			})
		}

		logger.Info("HTTP agent created for %s request to: http://%s:%d%s", payload.Method, ipAddress, port, payload.Path)

		agent.Body([]byte(payload.Body))
		for key, value := range payload.Headers {
			agent.Set(key, value)
		}

		StartCollectingRequests()

		logger.Info("Sending SysTest request...")
		if _, _, errs := agent.Bytes(); errs != nil {
			logger.Error("Failed to send request:")
			for _, err := range errs {
				logger.Error(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).JSON(payloads.ErrorResponse{
				Status:  fiber.StatusInternalServerError,
				Message: "Failed to send the request, read the logs for more information",
			})
		}

		logger.Info("The SysTest request was sent successfully.")
		return ctx.Status(fiber.StatusOK).JSON(StopCollectingRequests())
	})
}
