package semantics

import (
	"strconv"
	. "translator/lexical/bases/model"
)

func mark[T MutableBase | ImmutableBase](code int, base T) string {
	mark_code := string(base) + strconv.Itoa(code)
	return mark_code // exm I3 I1 etc
}

func SemanticProcedure5(buffer *string) {
	empty_buffer := ""
	buffer = &empty_buffer
}
