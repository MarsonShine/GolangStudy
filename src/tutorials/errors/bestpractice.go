package errors

import "net/http"

type appHandler func(http.ResponseWriter, *http.Request) error

func viewRecord2(w http.ResponseWriter, r *http.Request) error {
	key := dbGetKey()
	if err := dbGet(key); err != nil {
		return err
	}
	return otherExecute(w, key)
}

func otherExecute(w http.ResponseWriter, key string) error {
	return nil
}

// 实现 http.Handler 接口
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
