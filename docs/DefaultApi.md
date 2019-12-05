# \DefaultApi

All URIs are relative to *http://agentplatform.grs.uh.cu/api/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**GetAgent**](DefaultApi.md#GetAgent) | **Get** /getAgent/{Name} | 
[**GetAgentsByFunction**](DefaultApi.md#GetAgentsByFunction) | **Get** /getAgentsForFunction/{Name} | 
[**GetAgentsNames**](DefaultApi.md#GetAgentsNames) | **Get** /getAllAgentsNames | 
[**GetPeers**](DefaultApi.md#GetPeers) | **Get** /getPeers | 
[**GetSimilarAgent**](DefaultApi.md#GetSimilarAgent) | **Get** /getSimilarAgents/{Name} | 
[**RegisterAgent**](DefaultApi.md#RegisterAgent) | **Post** /registerAgent | 



## GetAgent

> []Addr GetAgent(ctx, name)



Get the agent that follow a simple criteria

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string**| Name of the Agent | 

### Return type

[**[]Addr**](Addr.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAgentsByFunction

> [][]Addr GetAgentsByFunction(ctx, name)



Get the agents that match with the function name passed as params

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string**| Name of the Function | 

### Return type

[**[][]Addr**](array.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetAgentsNames

> []string GetAgentsNames(ctx, )



Get all agents names registered in the platforms 

### Required Parameters

This endpoint does not need any parameter.

### Return type

**[]string**

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetPeers

> []Addr GetPeers(ctx, )



Return all peers connected to the platform network 

### Required Parameters

This endpoint does not need any parameter.

### Return type

[**[]Addr**](Addr.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## GetSimilarAgent

> Agent GetSimilarAgent(ctx, name)



Get the agents that are similars to the agent passed as paramerter

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**name** | **string**| Name of the Agent | 

### Return type

[**Agent**](Agent.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## RegisterAgent

> RegisterAgent(ctx, body)



Register a new Agent in the platform

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**Agent**](Agent.md)| Agent to register | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: */*

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

