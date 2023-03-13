package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"os"
	"systest/log"
	"systest/payloads"
)

var eurekaLogger *log.Logger
var baseUrl string

func addCorsHeaders(ctx *fiber.Ctx) {
	ctx.Vary(fiber.HeaderOrigin)
	ctx.Vary(fiber.HeaderAccessControlRequestMethod)
	ctx.Vary(fiber.HeaderAccessControlRequestHeaders)
	ctx.Set(fiber.HeaderAccessControlAllowOrigin, "*")
	ctx.Set(fiber.HeaderAccessControlAllowMethods, "GET,POST,HEAD,PUT,DELETE,PATCH")
	ctx.Set(fiber.HeaderAccessControlAllowCredentials, "true")
	ctx.Set(fiber.HeaderAccessControlAllowHeaders, "*")
}

func forwardContext(ctx *fiber.Ctx) error {
	if err := proxy.Do(ctx, baseUrl+ctx.OriginalURL()); err != nil {
		return err
	}
	addCorsHeaders(ctx)
	return nil
}

func forwardRequest(agent *fiber.Agent, ctx *fiber.Ctx) error {
	agent.Set("Content-Type", ctx.Get("Content-Type"))
	agent.Set("Accept", ctx.Get("Accept"))

	code, body, errors := agent.Bytes()
	if errors != nil {
		eurekaLogger.Error("Failed to forward request:")
		for _, err := range errors {
			eurekaLogger.Error(err.Error())
		}
	}

	return ctx.Status(code).Send(body)
}

func forwardPostRequest(ctx *fiber.Ctx, payload interface{}) error {
	agent := fiber.Post(baseUrl + ctx.OriginalURL())

	if payload == nil {
		agent.Body(ctx.Body())
	} else {
		switch ctx.Get("Content-Type") {
		case "application/json":
			agent.JSON(payload)
			break
		case "application/xml":
			agent.XML(payload)
			break
		default:
			eurekaLogger.Warn("Unsupported content type: %s", ctx.Get("Content-Type"))
			agent.Body(ctx.Body())
		}
	}

	return forwardRequest(agent, ctx)
}

func RegisterEurekaRoutes(app *fiber.App) {
	eurekaLogger = log.New("Eureka Proxy")
	baseUrl = os.Getenv("EUREKA_BASE_URL")

	eurekaLogger.Info("Registering Eureka proxy routes...")
	api := app.Group("/eureka")

	api.Get("/apps", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Querying for all application instances")
		return forwardContext(ctx)
	})

	api.Get("/apps/delta", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Querying for application instances with deltas")
		return forwardContext(ctx)
	})

	api.Get("/apps/:appId", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Querying for all application instances of: %s", ctx.Params("appId"))
		return forwardContext(ctx)
	})

	api.Post("/apps/:appId", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Registering a new instance of %s", ctx.Params("appId"))

		payload := &payloads.InstanceRegistrationRequest{}
		if err := ctx.BodyParser(payload); err != nil {
			eurekaLogger.Warn("The request could not be parsed, however it shall be forwarded: %s", err.Error())
		}

		if port, err := CreateInstanceProxy(payload.Instance); err == nil {
			eurekaLogger.Info("Instance proxy created on port: %d", port)
		} else {
			eurekaLogger.Warn("Failed to create instance proxy: %s", err.Error())
		}

		return forwardPostRequest(ctx, payload)
	})

	api.Get("/apps/:appId/:instanceId", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Querying for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return forwardContext(ctx)
	})

	api.Delete("/apps/:appId/:instanceId", func(ctx *fiber.Ctx) error {
		appId := ctx.Params("appId")
		instanceId := ctx.Params("instanceId")

		eurekaLogger.Info("De-registering %s instance: %s", appId, instanceId)

		var port int
		if port = GetPortByInstanceId(instanceId); port == 0 {
			port = GetPortByAppAndHostname(appId, instanceId)
		}

		if port == 0 {
			eurekaLogger.Warn("Cannot destroy proxy, none exists for: %s:%s", appId, instanceId)
		} else {
			DestroyInstanceProxy(port)
		}

		return forwardContext(ctx)
	})

	api.Put("/apps/:appId/:instanceId", func(ctx *fiber.Ctx) error {
		eurekaLogger.Debug("Heartbeat received for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return forwardContext(ctx)
	})

	api.Put("/apps/:appId/:instanceId/status", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Updating status override for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return forwardContext(ctx)
	})

	api.Delete("/apps/:appId/:instanceId/status", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Removing status override for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return forwardContext(ctx)
	})

	api.Put("/apps/:appId/:instanceId/metadata", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Updating metadata for %s instance: %s", ctx.Params("appId"), ctx.Params("instanceId"))
		return forwardContext(ctx)
	})

	api.Get("/instances/:instanceId", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Querying for instance: %s", ctx.Params("instanceId"))
		return forwardContext(ctx)
	})

	api.Get("/vips/:vipAddress", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Querying instances under vip address: %s", ctx.Params("vipAddress"))
		return forwardContext(ctx)
	})

	api.Get("/svips/:secureVipAddress", func(ctx *fiber.Ctx) error {
		eurekaLogger.Info("Querying for instances under secure vip address: %s", ctx.Params("secureVipAddress"))
		return forwardContext(ctx)
	})
}
