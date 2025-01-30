package helpers

import (
	"html/template"
	"net/http"
)

type ErrData struct {
	Code int
	Msg  string
}

func WriteError(w http.ResponseWriter, errMsg string, code int) ErrData {
	var Err ErrData
	Err.Code = code
	Err.Msg = errMsg
	return Err
}

func ExecuteTmpl(w http.ResponseWriter, temp string, code int, errMsg string, data any) {
	w.WriteHeader(code)
	tmpl := template.Must(template.ParseGlob("./assets/templates/*.html"))

	if temp == "error.html" {
		data = WriteError(w, errMsg, code)
	}
	err := tmpl.ExecuteTemplate(w, temp, data)

	if err != nil {
		err2 := tmpl.ExecuteTemplate(w, "error_page.html", WriteError(w, "Internal Server Error!", http.StatusInternalServerError))
		if err2 != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}
}
