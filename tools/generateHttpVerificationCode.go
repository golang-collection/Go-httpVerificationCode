package main

import (
	"fmt"
	"github.com/dchest/captcha"
	"log"
	"net/http"
)

/**
* @Author: super
* @Date: 2020-09-10 08:21
* @Description:
**/

//生成ID
func GenerateCodeID(w http.ResponseWriter, r *http.Request){
	id := captcha.NewLen(6)
	if _, err := fmt.Fprint(w, id); err != nil {
		log.Println("generate captcha error", err)
	}
}

//根据id生成对应的验证码
func GenerateCode(w http.ResponseWriter, r *http.Request){
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "need request id", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	if err := captcha.WriteImage(w, id, 120, 80); err != nil {
		log.Println("show captcha error", err)
	}
}

//根据id和传递过来的值对比内存中存储的id与验证码值是否一致
func VerifyCode(w http.ResponseWriter, r *http.Request){
	if err := r.ParseForm(); err != nil {
		log.Println("parseForm error", err)
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}
	// receive captcha-id and captcha-value
	id := r.FormValue("id")
	value := r.FormValue("value")
	// checking whether captcha-id and captcha-value is matched
	if captcha.VerifyString(id, value) {
		fmt.Fprint(w, "ok")
	} else {
		fmt.Fprint(w, "mismatch")
	}
}

func main() {
	// generate a new captcha, return captcha-id
	http.HandleFunc("/captcha/generate", GenerateCodeID)
	// show captcha image by captcha-id
	http.HandleFunc("/captcha/image", GenerateCode)
	// business logic
	http.HandleFunc("/login", VerifyCode)
	log.Fatal(http.ListenAndServe(":8080", nil))
}