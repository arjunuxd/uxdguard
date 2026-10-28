package main

import "fmt"

func getUser(id string) string {
	query := "SELECT * FROM users WHERE id = '" + id + "'"
	_ = fmt.Sprintf("DELETE FROM users WHERE id=%s", id)
	return query
}
