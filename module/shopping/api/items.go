package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/module/shopping"
	"github.com/PaluMacil/dwn/webserver/errs"
)

// POST api/shopping/items
func addHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if spouse, err := cur.Is(core.BuiltInGroupSpouse); err != nil {
		return err
	} else if !spouse {
		return errs.StatusForbidden
	}

	var item shopping.Item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid shopping item")}
	}
	item.Added = time.Now()
	item.AddedBy = cur.User.DisplayName
	err = db.Shopping.Items.Set(item)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(item)
}

// DELETE api/shopping/items
func removeHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if spouse, err := cur.Is(core.BuiltInGroupSpouse); err != nil {
		return err
	} else if !spouse {
		return errs.StatusForbidden
	}
	qry, err := url.QueryUnescape(vars["name"])
	if len(qry) == 0 || err != nil {
		return errs.StatusError{http.StatusBadRequest, errors.New("invalid or missing shopping item name")}
	}
	w.WriteHeader(http.StatusNoContent)
	return db.Shopping.Items.Delete(qry)
}

// GET api/shopping/items
func listHandler(
	db *database.Database,
	cur core.Current,
	vars map[string]string,
	w http.ResponseWriter,
	r *http.Request,
) error {
	if spouse, err := cur.Is(core.BuiltInGroupSpouse); err != nil {
		return err
	} else if !spouse {
		return errs.StatusForbidden
	}
	items, err := db.Shopping.Items.All()
	if err != nil {
		return fmt.Errorf("getting list of all shopping items: %w", err)
	}

	return json.NewEncoder(w).Encode(items)
}
