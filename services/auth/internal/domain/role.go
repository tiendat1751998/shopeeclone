package domain

type Role string

const (
	RoleAdmin     Role = "admin"
	RoleSeller    Role = "seller"
	RoleBuyer     Role = "buyer"
	RoleService   Role = "service"
	RoleSupport   Role = "support"
	RoleModerator Role = "moderator"
)

type Permission string

const (
	PermUserRead       Permission = "user:read"
	PermUserWrite      Permission = "user:write"
	PermUserDelete     Permission = "user:delete"
	PermProductRead    Permission = "product:read"
	PermProductWrite   Permission = "product:write"
	PermProductDelete  Permission = "product:delete"
	PermOrderRead      Permission = "order:read"
	PermOrderWrite     Permission = "order:write"
	PermInventoryRead  Permission = "inventory:read"
	PermInventoryWrite Permission = "inventory:write"
	PermPaymentRead    Permission = "payment:read"
	PermPaymentWrite   Permission = "payment:write"
	PermAdmin          Permission = "admin:*"
	PermAuditLog       Permission = "audit:read"
	PermSessionManage  Permission = "session:manage"
	PermRoleManage     Permission = "role:manage"
)

type RoleDefinition struct {
	Role        Role         `json:"role"`
	Permissions []Permission `json:"permissions"`
	Parent      *Role        `json:"parent,omitempty"`
	Description string       `json:"description"`
}

var RoleHierarchy = map[Role]*RoleDefinition{
	RoleAdmin: {
		Role:        RoleAdmin,
		Permissions: []Permission{PermAdmin, PermUserRead, PermUserWrite, PermUserDelete, PermRoleManage, PermAuditLog, PermSessionManage},
		Description: "Full system access",
	},
	RoleSeller: {
		Role:        RoleSeller,
		Permissions: []Permission{PermProductRead, PermProductWrite, PermProductDelete, PermOrderRead, PermInventoryRead, PermInventoryWrite},
		Description: "Seller account with product and inventory management",
	},
	RoleBuyer: {
		Role:        RoleBuyer,
		Permissions: []Permission{PermProductRead, PermOrderRead, PermOrderWrite, PermPaymentRead},
		Description: "Buyer account with order management",
	},
	RoleService: {
		Role:        RoleService,
		Permissions: []Permission{PermUserRead, PermOrderRead, PermInventoryRead, PermPaymentRead, PermProductRead},
		Description: "Internal service account",
	},
	RoleSupport: {
		Role:        RoleSupport,
		Permissions: []Permission{PermUserRead, PermOrderRead, PermOrderWrite},
		Description: "Customer support agent",
	},
	RoleModerator: {
		Role:        RoleModerator,
		Permissions: []Permission{PermUserRead, PermProductRead, PermProductWrite},
		Description: "Content moderator",
	},
}

func GetPermissionsForRole(role Role) []Permission {
	def, exists := RoleHierarchy[role]
	if !exists {
		return nil
	}
	return def.Permissions
}

func HasPermission(userRoles []Role, required Permission) bool {
	for _, role := range userRoles {
		perms := GetPermissionsForRole(role)
		for _, p := range perms {
			if p == PermAdmin || p == required {
				return true
			}
		}
	}
	return false
}

func IsAdmin(roles []Role) bool {
	for _, r := range roles {
		if r == RoleAdmin {
			return true
		}
	}
	return false
}
