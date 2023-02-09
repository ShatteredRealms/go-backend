package main

import (
    "context"
    "fmt"
    "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/model"
    "github.com/WilSimpson/ShatteredRealms/go-backend/pkg/service"
    log "github.com/sirupsen/logrus"
    "google.golang.org/grpc"
)

func seedDatabaseIfNeeded(
    ctx context.Context,
    userService service.UserService,
    permissionService service.PermissionService,
    roleService service.RoleService,
    servicesInfo map[string]grpc.ServiceInfo,
) {
    var err error
    superAdminRole := roleService.FindByName(ctx, "Super Admin")
    if superAdminRole.Model.ID == 0 {
        // Create Super Admin role
        superAdminRole, err = roleService.Create(
            ctx,
            &model.Role{
                Name: "Super Admin",
            },
        )
        if err != nil {
            log.WithContext(ctx).Errorf("creating super admin: %v", err)
            return
        }
    }

    currentPermissions := createSetOfPermissions(
        permissionService.FindPermissionsForRoleID(ctx, superAdminRole.ID),
    )

    // Assign all permissions with Other set to true
    for packageService, serviceInfo := range servicesInfo {
        for _, methodInfo := range serviceInfo.Methods {
            permission := fmt.Sprintf("/%s/%s", packageService, methodInfo.Name)
            // Only add the permission if it doesn't exist already
            if _, ok := currentPermissions[permission]; !ok {
                err = permissionService.AddPermissionForRole(ctx, &model.RolePermission{
                    RoleID:     superAdminRole.ID,
                    Permission: permission,
                    Other:      true,
                })

                if err != nil {
                    log.WithContext(ctx).Errorf("creating permission %s for super admin: %v", permission, err)
                }
            }
        }
    }

    if len(userService.FindAll(ctx)) > 0 {
        return
    }

    _, err = userService.Create(
        ctx,
        &model.User{
            FirstName: "Wil",
            LastName:  "Simpson",
            Username:  "unreal",
            Email:     "wil@forever.dev",
            Password:  "password",
            Roles:     []*model.Role{superAdminRole},
        })

    if err != nil {
        log.WithContext(ctx).Errorf("creating super user: %v", err)
    }
}

func createSetOfPermissions(permissions []*model.RolePermission) map[string]struct{} {
    out := make(map[string]struct{}, len(permissions))

    for _, v := range permissions {
        out[v.Permission] = struct{}{}
    }

    return out
}
