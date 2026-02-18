package model

import "path/filepath"

type ImmutableBase = int
type MutableBase = rune

var SystemWord ImmutableBase = 0
var Separators ImmutableBase = 1
var Operations ImmutableBase = 2

var SymConsts MutableBase = '3'
var BoolConsts MutableBase = '4'
var Identificators MutableBase = '5'
var NumConsts MutableBase = '6'

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
