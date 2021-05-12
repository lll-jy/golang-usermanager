package main

import (
	"fmt"
	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
	"git.garena.com/jiayu.li/entry-task/test/server_helpers"
	"net/http"
	"testing"
)

func Test_temp(t *testing.T) {
	i := 15
	client := &http.Client{}
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     fmt.Sprintf("testuser%d", i),
			Password: fmt.Sprintf("testpass%d%d", i*2, i*2),
		},
		&protocol.User{
			Name:     fmt.Sprintf("testuser%d", i),
			Password: fmt.Sprintf("testpass%d%d", i*2, i*2),
		},
		handlers.InfoErr{},
		paths.PlaceholderPath,
	)
	request := server_helpers.MakeRequest(http.MethodPost, "http://localhost:8080/delete", t)
	t.Logf("header is %v", request.Header)
	//request.Header.Set("Cookie", cookieString)
	// MTYyMDc4MzY0OHxXMHFtNThqS0VDTnFjbGlyVG1ETkcwMnVvWlk4RktIaS1NSnRmWWx2RURhRWU0SHdpVmRRazRfbWVBcHk1VHZ0bEEyUDc4WXFUbk90TXY3UHV6ZG83aUM5SS0tRVhWUG41LW9JM25pZE9HaDVmS25mUEFpOUJxWElWVWs5TWdwQnhueERtS0ptMEI1bWc1ZGs0Q1Z1RUhqSUJ0SzItQzZhVW5hOUFPLW9MRU5UbkNXWVRrNkxBV280ZlJrTW5VZVZFNXA0ckpmU3BiUE1EWXptTUlwTnBnSDBJcDJ1Q3dpWWUyaURfLUF5dURtd2VWMHJTdTFNVUVjQVZJbzRSQkl1WHlFUlY0SmZLNGV1XzA5bkpvV2VXU0NFLU96TlNsSVdsdHdlb1BqVXE2UUFwSXNDX0RXVnRkXzJTOHVFVGxhZUdJaUNuUm5KdVhmV09IZlpyUXQ4cDh4czdFMVpVUT09fNyDQWhBoNxUZRs4nJO1RMkS7Jfx4FRFkACmh8BhqV5z
	// MTYyMDc4MzU3M3xPNnpWSTMtUENtSzZDZ3kxVGNObjRJQzVZNGdpTHVGdnRoNVVHUzNfOEtUdUhXVFpuRkxsNUZvTTloUTJWMDBQNW4zVzM4OHlOXzJMQV9UUTIwQ2ZRVFZhdVBUY0VwazVoa3ZJUGcxVF8wYUZqam53aEZzZXFaVldMdHFUZ19QcHNZQWlvWmpOS0hDbTVEcTB6Y05rY3RGSm9kOVhBd2k4Yzl3em9GUFZrMk1FRTNxc3NUOVVJc2tGWFVzPXy7bM4i-K9htUjg7jelZCnZzRQ885tsefSVMIqVE5KdeQ==
	//  map[Accept-Encoding:[gzip] Content-Length:[0] Cookie:[MTYyMDc4NDQ0OXxsbFhlcFNkeDJobHBpUzZsRklhZWE4TXJoVXh1aFVBaVpHbHBsUmVjbi1oSUVuc0M4d2pWVlFaZnBMR2NRemZpNloyei1Wbm5ySG9pd1ZzZGhmcGJ4WWdNWUlhT0Z0VDlzell1NzlmSW40VlI1ODJoSUYyempScVdBS0wtWWFSbjRRWVFqbjlna3k5cXE1ODBIT0tJSFhuTW9QblFiaWdEOGRUa1doTU02ekpWRS1EdXNmWXZVTnJ5dDhZZGVVMWtwejg9fBT0P1dri4Yd7BzlLUdJxZnHdA-GTGWCClb6C2uV_FJY] User-Agent:[Go-http-client/1.1]]
	// map[

	//Accept:[text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9]
	//Accept-Encoding:[gzip, deflate, br]
	//Accept-Language:[en-US,en;q=0.9,zh-CN;q=0.8,zh;q=0.7]
	//Cache-Control:[max-age=0]
	//Connection:[keep-alive]
	//Content-Length:[0]
	//Content-Type:[application/x-www-form-urlencoded]
	//Cookie:[session=MTYyMDc4NDU5OXxUNG9NTWdXbG5ZN2hMYlhmV3ItMUpKamJvd1BaUlNBeVdwTWFFY1hiVlBPNXBuLWU0bEFKRDY1ZWpsQTJJYzlxT05idFhoWVF2WmZ2TFVMSllkVEdQbXlVVWtQcVN0NS0yc2VWbEo3YTV0YmF5Z3FnR1VWX081SVR2VXFHMlBtS0Rkak5vQ0tSQjRzZDZGTTVWSUptcFRBakN2RDAtYlJFUmFKTmV4RDg0aThvbUtmRFhOcVI4d2h5aFE5QWctT2tYTm12bUJCQ050Ni1xVFQzaDV2WjRFaHpRWWQycnAyVXFFUTlzakVNemxFSzhuNnJnQWRzSGhhV0N4aHVvUkZya3R1a0pwWlprY084T0NlbjV2Qnhwbm1hNkt0Sk5nbjVRbFdMdG1Scmd4WlNhRUFGSU1aNkctSG5SNGZpR05BcGxocUE5dlZCTlRQeFBQX1hQdndXUVRsN2Z3eVJfUT09fInnMlpdevRzCMISKdVDgFbo8HPn74QeYK7vtyObJINM]
	//Origin:[http://localhost:8080]
	//Referer:[http://localhost:8080/view]
	//Sec-Ch-Ua:[" Not A;Brand";v="99", "Chromium";v="90", "Google Chrome";v="90"]
	//Sec-Ch-Ua-Mobile:[?0]
	//Sec-Fetch-Dest:[document]
	//Sec-Fetch-Mode:[navigate]
	//Sec-Fetch-Site:[same-origin]
	//Sec-Fetch-User:[?1] Upgrade-Insecure-Requests:[1]
	//User-Agent:[Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.93 Safari/537.36]]
	//info is {name:"testuser15"  password:"$2a$04$aPudUZsMVhsCLgbfHMMmUOQqRboge8F.342hIlOm/NgXTHcdQq71e"  photoUrl:"assets/placeholder.jpeg" name:"testuser15"  password:"testpass3030"  photoUrl:"assets/placeholder.jpeg" {  }     assets/placeholder.jpeg}


	//t.Logf("header is %v", request.Header)
	request.AddCookie(&http.Cookie{Name: "session", Value: cookieString})
	t.Logf("header2 is %v", request.Header)
	//server_helpers.UpdateCookie(cookieString, httptest.NewRecorder(), request)
	resp, err := client.Do(request)
	if err != nil {
		t.Errorf("Error delete, %s", err.Error())
	}
	resp.Body.Close()
}