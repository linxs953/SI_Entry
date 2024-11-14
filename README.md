# Gateway Service

Gateway服务是集群的统一入口点，主要负责HTTP请求的路由和转发，以及与scheduler组件的交互。

## 功能特点

- 提供统一的HTTP API入口
- 请求转发至scheduler服务
- 从scheduler获取和管理CRD（自定义资源定义）资源
- 支持请求的负载均衡
- 提供请求的认证和授权功能
- 实现请求限流和熔断机制
- 支持API版本控制
- 提供详细的监控指标和日志

## 系统架构

</file>