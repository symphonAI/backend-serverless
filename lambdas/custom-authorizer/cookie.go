package main

import (
	"fmt"
	"strings"
)

func extractJwtFromCookie(cookieStr string) (jwt *string){
		// Parse the cookie value to obtain the token
		cookie := strings.Split(cookieStr, ";")
		var tokenString string
		for _, c := range cookie {
			fmt.Println("Cookie:::", cookie)
			if strings.Contains(c, "jwt=") {
				tokenString = strings.TrimPrefix(c, "jwt=")
				break
			}
		}
		return &tokenString
}