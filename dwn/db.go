package dwn

import (
	"github.com/asdine/storm"
)

// Db is a global database object which I am likely to move or replace
// with another approach
var Db *storm.DB
