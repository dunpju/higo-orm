package DaoException

import "github.com/dunpju/higo-throw/exception"

func Throw(message string, code int) {
	exception.Throw(exception.Code(code), exception.Message(message), exception.Data(""))
}
