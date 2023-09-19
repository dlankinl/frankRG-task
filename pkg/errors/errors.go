package errors

import "errors"

var NoDirsWereFoundErr = errors.New("no dirs were found")
var NoFilesWereFoundErr = errors.New("no files were fould")
var TypeNotFileErr = errors.New("that's not a file but directory")
var NothingFoundToRenameError = errors.New("there's nothing to rename")
