package threat

import "errors"

var ErrMD5Exists = errors.New("Данный md5 уже занят")
var ErrMD5NotFound = errors.New("Угроза по данному md5 не найдена")
var ErrMD5IsNotAllowed = errors.New("Недопустимый md5")
var ErrSeverityIsNotAllowed = errors.New("Недопустимый уровень угрозы")