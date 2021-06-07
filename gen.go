package buns

//go:generate mockgen -destination mocks/tntcluster/pool/mock_connector.go github.com/city-mobil/gobuns/tntcluster/pool ConnectorPool
//go:generate mockgen -destination mocks/external/mock_external.go github.com/city-mobil/gobuns/external Client
//go:generate mockgen -destination mocks/health/mock_checkable.go github.com/city-mobil/gobuns/health Checkable
//go:generate mockgen -destination mocks/registry/mock_registry.go github.com/city-mobil/gobuns/registry Client
//go:generate mockgen -destination mocks/mysql/mock_mysql.go github.com/city-mobil/gobuns/mysql Adapter
//go:generate mockgen -destination mocks/redis/mock_redis.go github.com/city-mobil/gobuns/redis Redis
