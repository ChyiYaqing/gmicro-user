package role

func AccessibleRoles() map[string][]string {
	const userV1 = "/user.v1.User/"

	return map[string][]string{
		// userV1 + "Login":  {"admin", "user"},
		userV1 + "Create": {"admin"},
		userV1 + "Get":    {"admin", "user"},
	}
}
