package queries

import (
	"github.com/j13v/elrincondalba-ms/util"
)

var Root = util.CombineFields(QueryArticle, QueryOrder, QueryCategories, QueryCatalog, QueryAuth)
