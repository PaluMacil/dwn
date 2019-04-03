package setup

import (
	"github.com/PaluMacil/dwn/module/core"
	"log"
	"time"

	"github.com/PaluMacil/dwn/database"
)

type Module struct {
	Db *database.Database
}

func New(db *database.Database) *Module {
	return &Module{
		Db: db,
	}
}

const setupUser = "(SETUP)"

func (mod Module) Ensure() error {
	complete, err := mod.Db.SetupInfo.Completed()
	if err != nil {
		return err
	}
	if complete {
		log.Println("app already initialized")
		return nil
	}

	builtinGroups := []core.Group{
		{
			Name:         core.BuiltInGroupAdmin,
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{
			Name: core.BuiltInGroupSpouse,
			Permissions: []string{
				core.PermissionEditUserInfo,
				core.PermissionUnlockUser,
				core.PermissionViewUsers,
				core.PermissionEditGroups,
				core.PermissionViewGroups,
				core.PermissionEditAllVisitPages,
			},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{
			Name: core.BuiltInGroupResident,
			Permissions: []string{
				core.PermissionManageIOTDevices,
			},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{

			Name:         core.BuiltInGroupFriend,
			Permissions:  []string{},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{

			Name:         core.BuiltInGroupLandlord,
			Permissions:  []string{},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{

			Name:         core.BuiltInGroupTenant,
			Permissions:  []string{},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{

			Name:        core.BuiltInGroupPowerUser,
			Permissions: []string{core.PermissionListProjects},
			ModifiedBy:  setupUser,
		},
		{

			Name:         core.BuiltInGroupUser,
			Permissions:  []string{core.PermissionPostComments},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
	}
	for _, group := range builtinGroups {
		err = mod.Db.Groups.Set(group)
		if err != nil {
			return err
		}
	}

	err = mod.Db.SetupInfo.Set(core.SetupInfo{time.Now()})
	if err != nil {
		return err
	}
	log.Println("new app setup complete")
	return nil
}
