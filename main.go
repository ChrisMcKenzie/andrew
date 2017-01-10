package main

import (
	"encoding/xml"
	"fmt"

	"github.com/huin/goupnp"
	"github.com/lnguyen/wemo"
)

func main() {
	devices, err := Switches()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%+v", devices)
	for _, sw := range devices {
		sw.Off()
	}
}

func Switches() ([]wemo.Switch, error) {
	var switches []wemo.Switch
	devices, err := goupnp.DiscoverDevices("urn:Belkin:device:*")
	if err != nil {
		return switches, err
	}

	for _, device := range devices {
		var wemoSwitch wemo.Switch
		var setup wemo.Setup
		setupXML := wemo.Get(device.Root.URLBaseStr)
		xml.Unmarshal(setupXML, &setup)
		wemoSwitch.Host = device.Root.URLBase.Host
		wemoSwitch.Name = setup.FriendlyName
		switches = append(switches, wemoSwitch)
	}
	return switches, nil
}
