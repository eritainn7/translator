package model

import "path/filepath"

type MutableBase interface{}
type ImmutableBase interface{}

type systemWord = ImmutableBase
type separators = ImmutableBase
type operations = ImmutableBase

type symConsts = MutableBase
type boolConsts = MutableBase
type identificators = MutableBase
type numConsts = MutableBase

var SystemWord systemWord = 0
var Separators separators = 1
var Operations operations = 2

var SymConsts symConsts = 3
var BoolConsts boolConsts = 4
var Identificators identificators = 5
var NumConsts numConsts = 6

var path_to_system_word_base, _ = filepath.Abs("lexical/bases/tables/system_words.csv")
var path_to_separators_base, _ = filepath.Abs("lexical/bases/tables/separators.csv")
var path_to_operations_base, _ = filepath.Abs("lexical/bases/tables/operations.csv")
var path_to_sym_consts_base, _ = filepath.Abs("lexical/bases/tables/sym_consts.csv")
var path_to_bool_consts_base, _ = filepath.Abs("lexical/bases/tables/bool_consts.csv")
var path_to_identificators_base, _ = filepath.Abs("lexical/bases/tables/identificators.csv")
var path_to_num_consts_base, _ = filepath.Abs("lexical/bases/tables/num_consts.csv")

var endpoints map[any]string = map[any]string{
	SystemWord:     path_to_system_word_base,
	Separators:     path_to_separators_base,
	Operations:     path_to_operations_base,
	SymConsts:      path_to_sym_consts_base,
	BoolConsts:     path_to_bool_consts_base,
	Identificators: path_to_identificators_base,
	NumConsts:      path_to_num_consts_base,
}
