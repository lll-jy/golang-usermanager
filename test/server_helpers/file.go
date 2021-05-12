package server_helpers

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git.garena.com/jiayu.li/entry-task/cmd/handlers"
	"git.garena.com/jiayu.li/entry-task/cmd/paths"
	"git.garena.com/jiayu.li/entry-task/cmd/protocol"
)

// https://github.com/gobuffalo/httptest/blob/master/file.go
func Upload(t *testing.T, db *sql.DB, i int) string {
	// https://www.programmersought.com/article/6833575288/
	filename := fmt.Sprintf("data/original/sample%d.jpeg", i%3+1)
	fieldname := "photo_file"
	bb := &bytes.Buffer{}
	writer := multipart.NewWriter(bb)
	defer writer.Close()
	part, err := writer.CreateFormFile(fieldname, filename)
	if err != nil {
		t.Errorf("The file cannot be created as form file.")
		writer.Close()
	}
	file, err := os.Open(filename)
	if err != nil {
		t.Errorf("File %s not found.", filename)
	}
	io.Copy(part, file)
	contentType := writer.FormDataContentType()
	file.Close()
	writer.Close()
	request, err := http.NewRequest(http.MethodPost, "/upload", bb)
	if err != nil {
		t.Errorf("Cannot make request.")
	} else {
		request.Header.Set("Content-Type", contentType)
	}
	response := httptest.NewRecorder()
	name := fmt.Sprintf("user%d", i)
	pass := fmt.Sprintf("pass%d%d", i*2, i*2)
	nickname := fmt.Sprintf("nick%d", i)
	cookieString := handlers.SetSessionInfo(
		&protocol.User{
			Name:     name,
			Password: pass,
			Nickname: nickname,
		},
		&protocol.User{
			Name:     name,
			Password: pass,
			Nickname: nickname,
		},
		handlers.InfoErr{},
		paths.PlaceholderPath,
	)
	UpdateCookie(cookieString, response, request)
	http.HandlerFunc(makeHandler(db, handlers.UploadHandler)).ServeHTTP(response, request)
	header := response.Header()
	photo := header["Tempphoto"][0]
	flag, err := areIdenticalFiles(filename, photo)
	if err != nil {
		t.Errorf("The files %s are invalid. %s.", photo, err.Error())
	} else if !flag {
		t.Errorf("The file copied is wrong.")
	}
	return header["Photo"][0]
}
