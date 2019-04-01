package setup

import (
	"log"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/dwn"
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

	builtinGroups := []dwn.Group{
		{
			Name:         dwn.BuiltInGroupAdmin,
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{
			Name: dwn.BuiltInGroupSpouse,
			Permissions: []string{
				dwn.PermissionEditUserInfo,
				dwn.PermissionUnlockUser,
				dwn.PermissionViewUsers,
				dwn.PermissionEditGroups,
				dwn.PermissionViewGroups,
				dwn.PermissionEditAllVisitPages,
			},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{
			Name: dwn.BuiltInGroupResident,
			Permissions: []string{
				dwn.PermissionManageIOTDevices,
			},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{

			Name:         dwn.BuiltInGroupFriend,
			Permissions:  []string{},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{

			Name:         dwn.BuiltInGroupLandlord,
			Permissions:  []string{},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{

			Name:         dwn.BuiltInGroupTenant,
			Permissions:  []string{},
			ModifiedBy:   setupUser,
			ModifiedDate: time.Now(),
		},
		{

			Name:        dwn.BuiltInGroupPowerUser,
			Permissions: []string{dwn.PermissionListProjects},
			ModifiedBy:  setupUser,
		},
		{

			Name:         dwn.BuiltInGroupUser,
			Permissions:  []string{dwn.PermissionPostComments},
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

	err = mod.Db.SetupInfo.Set(dwn.SetupInfo{time.Now()})
	if err != nil {
		return err
	}
	log.Println("new app setup complete")
	return nil
}
