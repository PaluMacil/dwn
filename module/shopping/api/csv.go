package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/PaluMacil/dwn/database"
	"github.com/PaluMacil/dwn/module/core"
	"github.com/PaluMacil/dwn/module/shopping"
	"github.com/PaluMacil/dwn/webserver/errs"
	"io"
	"net/http"
	"strconv"
	"time"
)

// POST api/shopping/items/csv
func importCSVHandler(
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
	file, _, err := r.FormFile("file")
	if err != nil {
		return err
	}
	defer file.Close()
	rdr := csv.NewReader(file)
	for i := 0; true; i++ {
		record, err := rdr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if i == 0 {
			continue
		}
		if len(record) != 3 {
			return errs.StatusError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("parsing csv upload: expected receord length 3, got %d", len(record)),
			}
		}
		quantity, err := strconv.Atoi(record[2])
		if err != nil {
			return errs.StatusError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("parsing csv upload: expected numeric quantity field, got %s", record[2]),
			}
		}
		item := shopping.Item{
			Name:     record[0],
			Quantity: quantity,
			Note:     record[1],
			AddedBy:  cur.User.DisplayName,
			Added:    time.Now(),
		}
		err = db.Shopping.Items.Set(item)
		if err != nil {
			return err
		}
	}
	items, err := db.Shopping.Items.All()
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(items)
}
