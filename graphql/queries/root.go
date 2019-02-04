package queries

import (
	"github.com/jal88/elrincondalba-ms/util"
)

var Root = util.CombineFields(QueryArticle, QueryOrder, QueryCategories)
