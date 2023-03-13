package main

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"os"
	"strconv"
	"strings"
	"systest/log"
	"systest/payloads"
)

type InstanceProxy struct {
	Instance *payloads.EurekaInstance
	App      *fiber.App
	logger   *log.Logger
}

var InstanceProxies = make(map[int]*InstanceProxy)

func getPorts() []int {
	var ports []int
	portStrings := strings.Split(os.Getenv("PROXY_PORTS"), ",")

	if len(portStrings) == 0 {
		portStrings = []string{"10000", "10100"}
	}

	if len(portStrings) == 2 {
		start, err := strconv.Atoi(portStrings[0])
		if err != nil {
			start = 10000
		}

		end, err := strconv.Atoi(portStrings[1])
		if err != nil {
			end = 10100
		}

		length := end - start
		ports = make([]int, length)

		for i := 0; i < length; i++ {
			ports[i] = start + i
		}
	} else {
		ports = make([]int, len(portStrings))

		for i, portString := range portStrings {
			if port, err := strconv.Atoi(portString); err == nil {
				ports[i] = port
			}
		}
	}

	return ports
}

func nextAvailablePort() int {
	for _, port := range getPorts() {
		if _, ok := InstanceProxies[port]; !ok {
			return port
		}
	}
	return -1
}

func CreateInstanceProxy(instance *payloads.EurekaInstance) (int, error) {
	name := instance.App + ":" + instance.HostName

	logger := log.New(name + " Proxy")
	logger.Info("Creating instance proxy...")

	baseUrl := "http://" + instance.IpAddress + ":" + strconv.Itoa(instance.Port.Port)

	ipAddress := os.Getenv("SERVER_IP_ADDRESS")
	port := nextAvailablePort()

	if port < 0 {
		return port, errors.New("no remaining instance proxy ports")
	}

	instance.IpAddress = ipAddress
	instance.Port.Port = port

	app := fiber.New(fiber.Config{
		AppName:               name,
		DisableStartupMessage: true,
	})

	InstanceProxies[port] = &InstanceProxy{
		Instance: instance,
		App:      app,
		logger:   logger,
	}

	app.Use(recover.New())

	app.All("/*", func(ctx *fiber.Ctx) error {
		CollectRequest(instance, ctx)

		if err := proxy.Do(ctx, baseUrl+ctx.OriginalURL()); err != nil {
			logger.Error("Failed to proxy request: %s", err.Error())
			return err
		}

		CollectResponse(ctx)
		return nil
	})

	logger.Info("Starting listen server on: %s:%d", ipAddress, port)
	go func() {
		if err := app.Listen(ipAddress + ":" + strconv.Itoa(port)); err != nil {
			logger.Error("Failed to start listen server: %s", err.Error())
		}
	}()

	return port, nil
}

func GetPortByInstanceId(instanceId string) int {
	for port, instanceProxy := range InstanceProxies {
		if instanceProxy.Instance.InstanceId == instanceId {
			return port
		}
	}
	return -1
}

func GetPortByAppAndHostname(app string, hostname string) int {
	for port, instanceProxy := range InstanceProxies {
		if instanceProxy.Instance.App == app && instanceProxy.Instance.HostName == hostname {
			return port
		}
	}
	return -1
}

func DestroyInstanceProxy(port int) {
	if _, ok := InstanceProxies[port]; ok {
		instanceProxy := InstanceProxies[port]
		instanceProxy.logger.Info("Destroying instance proxy...")

		if err := instanceProxy.App.Shutdown(); err != nil {
			instanceProxy.logger.Error("Failed to shutdown instance proxy: %s", err.Error())
		}

		delete(InstanceProxies, port)
	}
}
