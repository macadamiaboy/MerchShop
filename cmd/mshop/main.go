package main

import (
	"fmt"

	"github.com/macadamiaboy/AvitoMerchShop/internal/helpers/auth"
)

func main() {
	login := "phil16"

	token, err := auth.GenToken(login)
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println(token)

	if err = auth.Verify(token, login); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Everything's working correctly")
	}
}
