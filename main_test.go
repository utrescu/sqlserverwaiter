package main

import "testing"

func TestConnectWithSqlServer(t *testing.T) {
	expected := ""
	resultat := ""

	if resultat != expected {
		t.Errorf("La connectivitat no funciona, rebut:'%s' esperat:'%s'", resultat, expected)
	}
}
