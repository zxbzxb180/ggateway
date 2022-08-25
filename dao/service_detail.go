package dao

type ServiceDetail struct {
	Info          *ServiceInfo          `json:"info" description:"服务基本信息"`
	HttpRule      *ServiceHttpRule      `json:"http" description:"Http规则"`
	TcpRule       *ServiceTcpRule       `json:"tcp" description:"tcp规则"`
	GrpcRule      *ServiceGrpcRule      `json:"grpc" description:"grpc规则"`
	LoadBalance   *ServiceLoadBalance   `json:"load_balance" description:"负载均衡"`
	AccessControl *ServiceAccessControl `json:"access_control" description:"访问控制"`
}
