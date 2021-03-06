package djv_codes

import (
	"fmt"
	"log"
	"net/http"
)

const CODE_COOKIE = "djv_code"
const MAX_AGE = 60 * 60 * 24 * 365
const JSONP_CALLBACK = "callback"

type GetCodeHandler struct {
	CodeManager *CodeManager
}

type ListCodeHandler struct {
	CodeManager *CodeManager
}

func (handler *GetCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var code string
	codeCookie, err := r.Cookie(CODE_COOKIE)
	if err == nil {
		code = codeCookie.Value
	} else if err == http.ErrNoCookie {
		code, err = handler.CodeManager.NextCode()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		http.Error(w, err.Error(), 500)
	}

	log.Printf("code:%v ip:%v useragent:%v", code, r.RemoteAddr, r.UserAgent())

	cookie := &http.Cookie{
		Name:   CODE_COOKIE,
		Value:  code,
		MaxAge: MAX_AGE}
	http.SetCookie(w, cookie)

	callback := r.FormValue(JSONP_CALLBACK)
	if len(callback) > 0 {
		fmt.Fprintf(w, "%v(\"%v\");", callback, code)
	} else {
		fmt.Fprintf(w, code)
	}
}

func (handler *ListCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	codes, err := handler.CodeManager.ListCodes()
	if err != nil || len(codes) == 0 {
		http.Error(w, err.Error(), 500)
		return
	}

	chr := codes[0][0]
	chrCount := 0

	// assume codes are 5 characters long
	for _, code := range codes {
		if code[0] != chr {
			fmt.Fprint(w, "\n\n")
			chr = code[0]
			chrCount = 0
		}

		sep := "     "
		if (chrCount + 1) % 8 == 0 {
			sep = "\n"
		}

		fmt.Fprint(w, code + sep)
		chrCount++
	}
}
