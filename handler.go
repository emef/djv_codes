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
	codeManager *CodeManager
}

func NewGetCodeHandler(codesDir string, usedCodesPath string) (*GetCodeHandler, error) {
	manager, err := NewCodeManager(codesDir, usedCodesPath)
	if err != nil {
		return nil, err
	}

	return &GetCodeHandler{manager}, nil
}

func (handler *GetCodeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var code string
	codeCookie, err := r.Cookie(CODE_COOKIE)
	if err == nil {
		code = codeCookie.Value
	} else if err == http.ErrNoCookie {
		code, err = handler.codeManager.NextCode()
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
	} else {
		http.Error(w, err.Error(), 500)
	}

	log.Printf("code:%v ip:%v useragent:%v", code, r.RemoteAddr, r.UserAgent())

	cookie := &http.Cookie{
		Name: CODE_COOKIE,
		Value: code,
		MaxAge: MAX_AGE}
  http.SetCookie(w, cookie)

	callback := r.FormValue(JSONP_CALLBACK)
	if len(callback) > 0 {
		fmt.Fprintf(w, "%v(%v);", callback, code)
	} else {
		fmt.Fprintf(w, code)
	}
}
