package goyt

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Insert TODO:
func (y YourTime) Insert(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	t, err := y.getParameters(r)
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}

	_, err = y.DB.Exec("INSERT INTO timemarks VALUES (DEFAULT, $1, $2, $3, $4, DEFAULT, $5, DEFAULT, DEFAULT, DEFAULT)", t.VideoID, t.IP, t.Timemark, t.Content, t.Author)
	if err != nil {
		fmt.Fprintf(w, sCError)
		fmt.Printf("%s", err)
		return
	}
	fmt.Fprintf(w, sCOK)
}

func (y YourTime) getParameters(r *http.Request) (Timemark, error) {
	t := Timemark{}

	videoID, err := getVideoID(r)
	if err != nil {
		return t, err
	}

	timemark, err := getTimemark(r)
	if err != nil {
		return t, err
	}

	content := getContent(r)

	author, err := y.getAuthor(r)
	if err != nil {
		return t, err
	}

	t = Timemark{
		VideoID:  videoID,
		IP:       strings.Split(r.RemoteAddr, ":")[0],
		Timemark: timemark,
		Content:  content,
		Author:   author,
	}

	return t, nil
}

func getVideoID(r *http.Request) (string, error) {
	if len(r.Form["videoid"]) > 0 {
		return r.Form["videoid"][0], nil
	}
	return "", errors.New("There was no videoID parameter")
}

func getTimemark(r *http.Request) (int64, error) {
	timemark := -1
	if len(r.Form["timemark"]) > 0 {
		timemark, err := strconv.Atoi(r.Form["timemark"][0])
		if err != nil {
			return int64(timemark), err
		}
		return int64(timemark), nil
	}
	return int64(timemark), errors.New("No timemark supplied")
}

func getContent(r *http.Request) string {
	if len(r.Form["content"]) > 0 {
		return r.Form["content"][0]
	}
	return ""
}

func (y YourTime) getAuthor(r *http.Request) (int64, error) {
	// User 1 is anonymous
	id := int64(1)
	token := getTokenFromCookies(r)
	if token == "" {
		return id, nil
	}
	rows, err := y.DB.Query("SELECT id FROM users WHERE token=$1", token)
	if err != nil {
		return id, err
	}
	if rows.Next() {
		err = rows.Scan(&id)
		return id, nil
	}
	if rows.Err() != nil {
		return id, rows.Err()
	}
	return id, nil
}

func (t timemarksDB) userExistsByToken(token token) (bool, error) {
	result := false
	rows, err := t.Query("SELECT exists(SELECT 1 FROM users WHERE token=$1)", token)
	if err != nil {
		return false, err
	}
	rows.Next()
	err = rows.Scan(&result)
	if err != nil {
		return false, err
	}
	return result, nil
}
