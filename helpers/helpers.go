package helpers

import (
	"github.com/graphql-go/graphql"
)


func CombineFields( fields ...graphql.Fields) graphql.Fields {
   resultFields := graphql.Fields{}

   for _, field := range fields {
      for nameProp, fieldDefinition := range field {
          resultFields[nameProp] = fieldDefinition
      }
   }

   return resultFields;
}
