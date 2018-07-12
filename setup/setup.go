package setup

import (
	"log"
	"time"

	"github.com/PaluMacil/dwn/app"
	"github.com/PaluMacil/dwn/db"
)

type Module struct {
	*app.App
}

func New(app *app.App) *Module {
	return &Module{
		App: app,
	}
}

func (mod Module) Ensure() error {
	complete, err := mod.Db.SetupInfo.Completed()
	if err != nil {
		return err
	}
	if complete {
		log.Println("app already initialized")
		return nil
	}

	builtinGroups := []db.Group{
		db.Group{
			Name:         db.BuiltInGroupAdmin,
			ModifiedBy:   "(SETUP)",
			ModifiedDate: time.Now(),
		},
		db.Group{
			Name: db.BuiltInGroupSpouse,
			Permissions: []string{
				db.PermissionEditUserInfo,
				db.PermissionUnlockUser,
				db.PermissionViewUsers,
				db.PermissionEditGroups,
				db.PermissionViewGroups,
				db.PermissionEditAllVisitPages,
			},
			ModifiedBy:   "(SETUP)",
			ModifiedDate: time.Now(),
		},
		db.Group{
			Name: db.BuiltInGroupResident,
			Permissions: []string{
				db.PermissionManageIOTDevices,
			},
			ModifiedBy:   "(SETUP)",
			ModifiedDate: time.Now(),
		},
		db.Group{

			Name:         db.BuiltInGroupFriend,
			Permissions:  []string{},
			ModifiedBy:   "(SETUP)",
			ModifiedDate: time.Now(),
		},
		db.Group{

			Name:         db.BuiltInGroupLandlord,
			Permissions:  []string{},
			ModifiedBy:   "(SETUP)",
			ModifiedDate: time.Now(),
		},
		db.Group{

			Name:         db.BuiltInGroupTenant,
			Permissions:  []string{},
			ModifiedBy:   "(SETUP)",
			ModifiedDate: time.Now(),
		},
		db.Group{

			Name:         db.BuiltInGroupUser,
			Permissions:  []string{db.PermissionPostComments},
			ModifiedBy:   "(SETUP)",
			ModifiedDate: time.Now(),
		},
	}
	for _, group := range builtinGroups {
		err = mod.Db.Groups.Set(group)
		if err != nil {
			return err
		}
	}

	err = mod.Db.SetupInfo.Set(db.SetupInfo{time.Now()})
	if err != nil {
		return err
	}
	log.Println("new app setup complete")
	return nil
}
