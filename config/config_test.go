package config

import (
	"testing"
)

func TestNewEmpApp(t *testing.T)  {
	app := NewEmpApp()
	if app == nil {
		t.Fatal("Couldn't create app!")
	}
	t.Log("TestNewEmpApp")
}

//func TestExecute(t *testing.T)  {
//	Execute()
//}