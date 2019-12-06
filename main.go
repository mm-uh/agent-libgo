package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	lib "github.com/mm-uh/agent-libgo/src"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"strings"
)

func launchServer(ip string, port string, function func(string2 string) (string, error)) {

	fmt.Println("Launching server...")

	// listen on all interfaces
	ln, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		logrus.Warn("Error listen server ", ip+":"+port)
	}
	// accept connection on port
	for {
		conn, err := ln.Accept()
		if err != nil {
			logrus.Warn("Error establishing connection")
		}

		// run loop forever (or until ctrl-c)
		// will listen for message to process ending in newline (\n)
		message, _ := bufio.NewReader(conn).ReadString('\n')
		// sample process for string received
		responseMessage, err := function(message)
		if err != nil {
			logrus.Warn("Error reading response")
			continue
		}
		// send new string back to client
		_, err = conn.Write([]byte(responseMessage))
		if err != nil {
			logrus.Warn("Error sending message")
		}
	}
}

func main() {
	name := "AdderGo"
	function := "Add"
	host := os.Args[1]
	port := os.Args[2]
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
	go launchServer("localhost", "38090", func(message string) (string, error) {
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
	go launchServer("localhost", "38091", func(message string) (string, error) {
		if message != "IsAlive?\n" {
			return "", errors.New("unknown param, not \"IsAlive?\"")
		}
		return "Yes", nil
	})
	go launchServer("localhost", "38092", func(message string) (string, error) {
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
	config := lib.NewConfiguration()
	config.BasePath = "http://" + host + ":" + port + "/api/v1"
	api := lib.NewAPIClient(config)
	defer context.Background()

	logrus.Info("Getting peers")
	peers, _, err := api.DefaultApi.GetPeers(context.Background())
	if err != nil {
		logrus.Warn("Error getting peers ", err.Error())
	}
	for _, peer := range peers {
		fmt.Println("IP: " + peer.Ip)
		fmt.Print("PORT: ")
		fmt.Print(peer.Port)
	}

	logrus.Info("Registering Agent")
	_, err = api.DefaultApi.RegisterAgent(context.Background(), agent)
	if err != nil {
		logrus.Warn("Error Registering agent ", err.Error())
	}

	logrus.Info("Getting Agent")
	agentLocation, _, err := api.DefaultApi.GetAgent(context.Background(), agent.Name)
	if err != nil {
		logrus.Warn("Error Registering agent ", err.Error())
	}
	fmt.Println(agentLocation)

	logrus.Info("Getting Agent")
	var similar lib.Agent
	similar, _, err = api.DefaultApi.GetSimilarAgent(context.Background(), agent.Name)
	if err != nil {
		logrus.Warn("Error Registering agent ", err.Error())
	}
	_, err = fmt.Println(similar.Name)
	if err != nil {
		logrus.Warn("Couldn't get similar.Name")
	}

}
