package main

import (
	"context"
	"errors"
	"github.com/mm-uh/agent-libgo/src"
	platform "github.com/mm-uh/go-agent-platform/src"
	"github.com/sirupsen/logrus"
	"math/rand"
	"net"
	"strconv"
)

type AgentWrapper struct {
	Name              string
	NodeAddr          src.Addr
	IsAliveAddr       src.Addr
	DocumentationAddr src.Addr
}

func NewAgentWrapper(agentName string, addresses []src.Addr) *AgentWrapper {
	if len(addresses) != 3 {
		return &AgentWrapper{}
	}
	return &AgentWrapper{
		Name:              agentName,
		NodeAddr:          addresses[0],
		IsAliveAddr:       addresses[1],
		DocumentationAddr: addresses[2],
	}
}

func (agent *AgentWrapper) SendToAgent(wrapper *PlatformWrapper, request string) (string, error) {
	if !agent.refreshAgent(wrapper, false) {
		return "", nil
	}

	response, err := platform.MakeRequest(getEndpoint(agent.NodeAddr, 0), request)
	if err != nil {
		return "", nil
	}
	return response, nil
}

func getEndpoint(addr src.Addr, plus int) string {
	port := strconv.Itoa(int(addr.Port) + plus)
	return addr.Ip + port
}

func (agent *AgentWrapper) IsAlive(wrapper *PlatformWrapper) bool {
	return platform.NodeIsAlive(getEndpoint(agent.IsAliveAddr, 0))
}

func (agent *AgentWrapper) GetDocumentation(wrapper *PlatformWrapper) (string, error) {
	if !agent.refreshAgent(wrapper, false) {
		return "", nil
	}

	response, err := platform.MakeRequest(getEndpoint(agent.DocumentationAddr, 0), "Doc")
	if err != nil {
		return "", nil
	}
	return response, nil
}

func (agent *AgentWrapper) refreshAgent(wrapper *PlatformWrapper, forced bool) bool {
	port := strconv.Itoa(int(agent.IsAliveAddr.Port))
	if isOpen(agent.IsAliveAddr.Ip, port) && !forced {
		return true
	}
	tmpAgent, err := wrapper.GetAgent(agent.Name)
	if err != nil {
		return false
	}
	*agent = *tmpAgent
	return true
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

	_ = wrapper.GetPeers()

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
func (wrapper *PlatformWrapper) GetPeers() error {

	defer context.Background()

	tmpPeers, _, err := wrapper.ApiClient.DefaultApi.GetPeers(context.Background())
	if err != nil {
		logrus.Warn("Could't get peers")
		return err
	}
	wrapper.KnowPeers = Union(wrapper.KnowPeers, tmpPeers)

	err = wrapper.UpdatePeers()
	if err != nil {
		logrus.Warn("Could't update peers")
		return err
	}

	return nil
}

func (wrapper *PlatformWrapper) UpdatePeers() error {
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

func (wrapper *PlatformWrapper) UpdateApi() error {
	_ = wrapper.GetPeers()

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

func (wrapper *PlatformWrapper) GetAgent(agentName string) (*AgentWrapper, error) {
	defer context.Background()

	err := wrapper.UpdateApi()
	if err != nil {
		return &AgentWrapper{}, err
	}
	agent, _, err := wrapper.ApiClient.DefaultApi.GetAgent(context.Background(), agentName)
	if err != nil {
		return &AgentWrapper{}, err
	}
	return NewAgentWrapper(agentName, agent), nil
}

func (wrapper *PlatformWrapper) GetSimilar(agentName string) (*AgentWrapper, error) {
	defer context.Background()

	err := wrapper.UpdateApi()
	if err != nil {
		return &AgentWrapper{}, err
	}
	agent, _, err := wrapper.ApiClient.DefaultApi.GetSimilarAgent(context.Background(), agentName)
	if err != nil {
		return &AgentWrapper{}, err
	}
	return NewAgentWrapper(agentName, agent), nil
}

func (wrapper *PlatformWrapper) RegisterAgent(agent *src.Agent) error {
	defer context.Background()

	err := wrapper.UpdateApi()
	if err != nil {
		return err
	}
	_, err = wrapper.ApiClient.DefaultApi.RegisterAgent(context.Background(), *agent)
	if err != nil {
		return err
	}
	return nil
}

func (wrapper *PlatformWrapper) GetAllAgents() ([]string, error) {
	defer context.Background()

	err := wrapper.UpdateApi()
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
