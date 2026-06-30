package persistent

import _ "embed"

//go:embed query/registry/insert_service.sql
var insertServiceQuery string

//go:embed query/registry/insert_endpoint.sql
var insertEndpointQuery string

//go:embed query/registry/list_service_with_endpoints.sql
var selectServiceWithEndpointsQuery string

//go:embed query/registry/select_endpoint_by_service.sql
var selectEndpointsByServiceQuery string
