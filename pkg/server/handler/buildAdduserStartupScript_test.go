package handler

import (
	"fmt"
	"testing"
)

func TestVmHandler_BuildAdduserStartupScript(t *testing.T) {
	script := buildAdduserStartupScript("root", "Passw0rd@_")

	fmt.Println(script)

}
