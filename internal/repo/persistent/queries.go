package persistent

import _ "embed"

//go:embed query/identity_context/delete_user.sql
var deleteUserQuery string

//go:embed query/identity_context/insert_user.sql
var insertUserQuery string

//go:embed query/identity_context/select_user.sql
var selectUserQuery string

//go:embed query/identity_context/update_user.sql
var updateUserQuery string
