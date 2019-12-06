package main

import (
	"context"
	"errors"
	"fmt"
	lib "github.com/mm-uh/agent-libgo/src"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func AgentMultiplier(name, function, hostP, portP string) {
	test := []lib.TestCase{
		{
			Input:  "1 2",
			Output: "3",
		},
		{
			Input:  "2 3",
			Output: "5",
		},
	}
	// Handler for agent
	go LaunchServer("localhost", "38090", func(message string) (string, error) {
		message = message[:len(message)-1]
		args := strings.Split(message, " ")
		if len(args) != 2 {
			return "", errors.New("unknown number of params")
		}
		a, err := strconv.Atoi(args[0])
		if err != nil {
			return "", err
		}
		b, err := strconv.Atoi(args[1])
		if err != nil {
			return "", err
		}
		return strconv.Itoa(a + b), nil
	})
	// Handler for IsAlive service
	go LaunchServer("localhost", "38091", func(message string) (string, error) {
		if message != "IsAlive?\n" {
			return "", errors.New("unknown param, not \"IsAlive?\"")
		}
		return "Yes", nil
	})
	go LaunchServer("localhost", "38092", func(message string) (string, error) {
		args := strings.Split(message, " ")
		if len(args) != 2 {
			return "", errors.New("unknown number of params")
		}
		a, err := strconv.Atoi(args[0])
		if err != nil {
			return "", err
		}
		b, err := strconv.Atoi(args[1])
		if err != nil {
			return "", err
		}
		return strconv.Itoa(a + b), nil
	})
	endpoints := []lib.Addr{
		{
			Ip:   "localhost",
			Port: 38090,
		},
	}
	documentation := map[string]lib.Addr{
		"localhost:38090": {
			Ip:   "localhost",
			Port: 38092,
		},
	}
	isAlive := map[string]lib.Addr{
		"localhost:38090": {
			Ip:   "localhost",
			Port: 38091,
		},
	}

	//# Agent | Agent to register
	//agent = lib_agent.Agent(name=name, function=function, endpoint_service=endpoint_service,
	//documentation=documentation, is_alive_service=is_alive_service, test_cases=test_case)
	agent := lib.Agent{
		Name:            name,
		Function:        function,
		EndpointService: endpoints,
		IsAliveService:  isAlive,
		Documentation:   documentation,
		TestCases:       test,
	}
	plastform := NewPlatformWWrapper(hostP, portP)

	logrus.Info("Getting peers")
	err := plastform.GetPeers()
	if err != nil {
		logrus.Warn("Error getting peers ", err.Error())
	}
	for _, peer := range plastform.KnowPeers {
		fmt.Println("IP: " + peer.Ip)
		fmt.Print("PORT: ")
		fmt.Print(peer.Port)
	}

	logrus.Info("Registering Agent")
	err = plastform.RegisterAgent(&agent)
	if err != nil {
		logrus.Warn("Error Registering agent ", err.Error())
	}

	logrus.Info("Getting Agent")
	agentLocation, err := plastform.GetAgent(agent.Name)
	if err != nil {
		logrus.Warn("Error Registering agent ", err.Error())
	}
	fmt.Println(agentLocation)

	logrus.Info("Getting Agent")
	var similar lib.Agent
	similar, _, err = api.DefaultApi.GetSimilarAgent(context.Background(), agentLocation.Name)
	if err != nil {
		logrus.Warn("Error Registering agent ", err.Error())
	}
	_, err = fmt.Println(similar.Name)
	if err != nil {
		logrus.Warn("Couldn't get similar.Name")
	}


}
