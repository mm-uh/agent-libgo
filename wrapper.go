package main

import (
	"context"
	"errors"
	"github.com/mm-uh/agent-libgo/src"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"strconv"
)

type AgentWrapper struct {
	Node          src.Addr
	IsAlive       src.Addr
	Documentation src.Addr
}

func NewAgentWrapper(addresses []src.Addr) *AgentWrapper {
	if len(addresses) != 3 {
		return &AgentWrapper{}
	}
	return &AgentWrapper{
		Node:          addresses[0],
		IsAlive:       addresses[1],
		Documentation: addresses[2],
	}
}

type PlatformWrapper struct {
	Host      string
	Port      string
	ApiClient *src.APIClient
	KnowPeers []src.Addr
}

func NewPlatformWWrapper(host, port string) *PlatformWrapper {
	wrapper := PlatformWrapper{
		Host:      host,
		Port:      port,
		KnowPeers: make([]src.Addr, 0),
	}

	config := src.NewConfiguration()
	config.BasePath = "http://" + host + ":" + port + "/api/v1"
	wrapper.ApiClient = src.NewAPIClient(config)

	_ = wrapper.getPeers()

	return &wrapper
}
func Union(a, b []src.Addr) []src.Addr {
	m := make(map[src.Addr]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; !ok {
			a = append(a, item)
		}
	}
	return a
}
func (wrapper *PlatformWrapper) getPeers() error {

	defer context.Background()

	tmpPeers, _, err := wrapper.ApiClient.DefaultApi.GetPeers(context.Background())
	if err != nil {
		logrus.Warn("Could't get peers")
		return err
	}
	wrapper.KnowPeers = Union(wrapper.KnowPeers, tmpPeers)

	err = wrapper.updatePeers()
	if err != nil {
		logrus.Warn("Could't update peers")
		return err
	}

	return nil
}

func (wrapper *PlatformWrapper) updatePeers() error {
	defer context.Background()

	if len(wrapper.KnowPeers) < 1 {
		return nil
	}
	for i := 0; i < 10; i++ {
		rnd := rand.Intn(len(wrapper.KnowPeers))
		port := strconv.Itoa(int(wrapper.KnowPeers[rnd].Port + 1000))
		if !isOpen(wrapper.KnowPeers[rnd].Ip, port) {
			continue
		}
		config := src.NewConfiguration()
		config.Host = "http://" + wrapper.KnowPeers[rnd].Ip + ":" + port + "/api/v1"
		apiClient := src.NewAPIClient(config)
		tmpPeers, _, err := apiClient.DefaultApi.GetPeers(context.Background())
		if err != nil {
			logrus.Warn("Could't get peers")
			continue
		}
		wrapper.KnowPeers = Union(wrapper.KnowPeers, tmpPeers)
	}

	return nil
}

func (wrapper *PlatformWrapper) updateApi() error {
	_ = wrapper.getPeers()

	if !isOpen(wrapper.Host, wrapper.Port) {
		for _, node := range wrapper.KnowPeers {
			port := strconv.Itoa(int(node.Port + 1000))
			if !isOpen(node.Ip, port) {
				wrapper.Port = port
				wrapper.Host = node.Ip

				config := src.NewConfiguration()
				config.BasePath = "http://" + wrapper.Host + ":" + port + "/api/v1"
				wrapper.ApiClient = src.NewAPIClient(config)
				return nil
			}
		}
		return errors.New("couldn't get any platform available")
	}

	return nil
}

func (wrapper *PlatformWrapper) getAgent(agentName string) (*AgentWrapper, error) {
	defer context.Background()

	err := wrapper.updateApi()
	if err != nil {
		return &AgentWrapper{}, err
	}
	agent, _, err := wrapper.ApiClient.DefaultApi.GetAgent(context.Background(), agentName)
	if err != nil {
		return &AgentWrapper{}, err
	}
	return NewAgentWrapper(agent), nil
}

func (wrapper *PlatformWrapper) registerAgent(agent *src.Agent) error {
	defer context.Background()

	err := wrapper.updateApi()
	if err != nil {
		return err
	}
	_, err = wrapper.ApiClient.DefaultApi.RegisterAgent(context.Background(), *agent)
	if err != nil {
		return err
	}
	return nil
}

func (wrapper *PlatformWrapper) getAllAgents() ([]string, error) {
	defer context.Background()

	err := wrapper.updateApi()
	if err != nil {
		return nil, err
	}
	agents, _, err := wrapper.ApiClient.DefaultApi.GetAgentsNames(context.Background())
	if err != nil {
		return nil, err
	}
	return agents, nil
}

func isOpen(host, port string) bool {
	ln, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}
