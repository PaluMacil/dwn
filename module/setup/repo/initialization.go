package repo

import (
	"fmt"
	"log"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/module/setup"
)

type InitializationRepo struct {
	store database.Storer
	db    *database.Database
}

func NewInitializationRepo(store database.Storer, db *database.Database) *InitializationRepo {
	return &InitializationRepo{store, db}
}

func (p InitializationRepo) Get() (setup.Initialization, error) {
	var initialization = setup.Initialization{}
	item, err := p.store.Get(&initialization)
	if err != nil {
		return initialization, err
	}
	initialization, ok := item.(setup.Initialization)
	if !ok {
		return initialization, fmt.Errorf("got data of type %T but wanted setup.Initialization", initialization)
	}
	return initialization, err
}

func (p InitializationRepo) dbInitCompleted() (bool, error) {
	_, err := p.Get()
	if p.db.IsKeyNotFoundErr(err) {
		return false, nil
	}
	return true, err
}

func (p InitializationRepo) Set(initialization setup.Initialization) error {
	return p.store.Set(&initialization)
}

// EnsureDatabase checks that the database has been set up. If it hasn't been, it performs initialization
func (p InitializationRepo) EnsureDatabase() error {
	complete, err := p.dbInitCompleted()
	if err != nil {
		return err
	}
	if complete {
		log.Println("database already initialized")
		return nil
	}

	// skip ID 0 for representation of the unauthenticated user
	skipID, err := p.db.NextID()
	if err != nil {
		return err
	} else if skipID != 0 {
		return fmt.Errorf("first number in sequence was not 0; it was %v", skipID)
	}

	setupUserID, err := p.db.NextID()
	if err != nil {
		return err
	}
	setupUser := core.User{
		ID:          setupUserID,
		DisplayName: core.DisplayName("(SETUP)"),
		Emails: []core.Email{
			{
				Email:    "(SETUP)",
				Verified: true,
			},
		},
		Locked: true,
	}
	if err := p.db.Users.Set(setupUser); err != nil {
		return err
	}

	envUserID, err := p.db.NextID()
	if err != nil {
		return err
	}
	envUser := core.User{
		ID:          envUserID,
		DisplayName: core.DisplayName("(ENV)"),
		Emails: []core.Email{
			{
				Email:    "(ENV)",
				Verified: true,
			},
		},
		Locked: true,
	}
	if err := p.db.Users.Set(envUser); err != nil {
		return err
	}

	builtinGroups := []core.Group{
		{
			Name:         core.BuiltInGroupAdmin,
			ModifiedBy:   setupUser.ID,
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
			},
			ModifiedBy:   setupUser.ID,
			ModifiedDate: time.Now(),
		},
		{
			Name: core.BuiltInGroupResident,
			Permissions: []string{
				core.PermissionManageIOTDevices,
			},
			ModifiedBy:   setupUser.ID,
			ModifiedDate: time.Now(),
		},
		{
			Name:         core.BuiltInGroupFriend,
			Permissions:  []string{},
			ModifiedBy:   setupUser.ID,
			ModifiedDate: time.Now(),
		},
		{
			Name:         core.BuiltInGroupPowerUser,
			Permissions:  []string{core.PermissionListProjects},
			ModifiedBy:   setupUser.ID,
			ModifiedDate: time.Now(),
		},
		{
			Name:         core.BuiltInGroupUser,
			Permissions:  []string{core.PermissionPostComments},
			ModifiedBy:   setupUser.ID,
			ModifiedDate: time.Now(),
		},
	}
	for _, group := range builtinGroups {
		err = p.db.Groups.Set(group)
		if err != nil {
			return err
		}
	}

	err = p.Set(setup.Initialization{
		DatabaseInitDate: time.Now(),
		WizardComplete:   false,
		SetupUserID:      setupUserID,
		EnvUserID:        envUserID,
	})
	if err != nil {
		return err
	}
	log.Println("new app setup complete")
	return nil
}
