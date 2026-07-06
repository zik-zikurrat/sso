package persistent

import _ "embed"

//go:embed query/registry/insert_service.sql
var insertServiceQuery string

//go:embed query/registry/insert_endpoint.sql
var insertEndpointQuery string

//go:embed query/registry/list_service_with_endpoints.sql
var selectServiceWithEndpointsQuery string

//go:embed query/registry/select_endpoints_by_service_id.sql
var selectEndpointsByServiceIDQuery string

//go:embed query/registry/select_service_by_id.sql
var selectServiceByIDQuery string

//go:embed query/registry/delete_service.sql
var deleteServiceQuery string

//go:embed query/registry/update_service.sql
var updateServiceQuery string

//go:embed query/identity_context/insert_user.sql
var insertUserQuery string

//go:embed query/identity_context/update_user.sql
var updateUserQuery string

//go:embed query/identity_context/delete_user.sql
var deleteUserQuery string

//go:embed query/identity_context/select_user.sql
var selectUserQuery string

//go:embed query/identity_context/insert_session.sql
var insertSessionQuery string
