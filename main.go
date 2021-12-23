package main

import (
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Nvmf yaml config struct
type Nvmf struct {
	Ports []struct {
		Name        string `yaml:"name"`
		AddrAdrfam  string `yaml:"addr_adrfam"`
		AddrTraddr  string `yaml:"addr_traddr"`
		AddrTrsvcid int    `yaml:"addr_trsvcid"`
		AddrTrtype  string `yaml:"addr_trtype"`
		Subsystems  []struct {
			Name string `yaml:"name"`
		} `yaml:"subsystems"`
	} `yaml:"ports"`
	Subsystems []struct {
		Name             string `yaml:"name"`
		AttrAllowAnyHost int    `yaml:"attr_allow_any_host"`
		Namespaces       []struct {
			Name       int    `yaml:"name"`
			Enable     int    `yaml:"enable" default:"1"`
			DevicePath string `yaml:"device_path"`
			DeviceUUID string `yaml:"device_uuid"`
		} `yaml:"namespaces"`
	} `yaml:"subsystems"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	_, err := os.Stat("/sys/kernel/config/nvmet")
	check(err)

	args := os.Args[1:]
	for _, arg := range args {
		dat, err := os.ReadFile(arg)
		check(err)

		nvmf := Nvmf{}
		err = yaml.Unmarshal([]byte(dat), &nvmf)
		check(err)

		subsystemPath := "/sys/kernel/config/nvmet/subsystems"

		for _, subsystem := range nvmf.Subsystems {
			path := filepath.Join(subsystemPath, subsystem.Name)
			err = os.MkdirAll(path, os.ModePerm)
			check(err)

			path = filepath.Join(subsystemPath, subsystem.Name, "attr_allow_any_host")
			err := os.WriteFile(path, []byte(strconv.Itoa(subsystem.AttrAllowAnyHost)), os.ModePerm)
			check(err)

			for _, namespaces := range subsystem.Namespaces {
				path = filepath.Join(subsystemPath, subsystem.Name, "namespaces", strconv.Itoa(namespaces.Name))
				err = os.MkdirAll(path, os.ModePerm)
				check(err)

				if namespaces.Enable == 0 {
					path = filepath.Join(subsystemPath, subsystem.Name, "namespaces", strconv.Itoa(namespaces.Name), "device_path")
					err := os.WriteFile(path, []byte(namespaces.DevicePath), os.ModePerm)
					check(err)

					path = filepath.Join(subsystemPath, subsystem.Name, "namespaces", strconv.Itoa(namespaces.Name), "device_uuid")
					err = os.WriteFile(path, []byte(namespaces.DeviceUUID), os.ModePerm)
					check(err)
				}

				path = filepath.Join(subsystemPath, subsystem.Name, "namespaces", strconv.Itoa(namespaces.Name), "enable")
				err = os.WriteFile(path, []byte(strconv.Itoa(namespaces.Enable)), os.ModePerm)
				check(err)
			}
		}

		portPath := "/sys/kernel/config/nvmet/ports"

		for _, port := range nvmf.Ports {
			path := filepath.Join(portPath, port.Name)
			err := os.MkdirAll(path, os.ModePerm)
			check(err)

			path = filepath.Join(portPath, port.Name, "addr_adrfam")
			err = os.WriteFile(path, []byte(port.AddrAdrfam), os.ModePerm)
			check(err)

			path = filepath.Join(portPath, port.Name, "addr_traddr")
			err = os.WriteFile(path, []byte(port.AddrTraddr), os.ModePerm)
			check(err)

			path = filepath.Join(portPath, port.Name, "addr_trsvcid")
			err = os.WriteFile(path, []byte(strconv.Itoa(port.AddrTrsvcid)), os.ModePerm)
			check(err)

			path = filepath.Join(portPath, port.Name, "addr_trtype")
			err = os.WriteFile(path, []byte(port.AddrTrtype), os.ModePerm)
			check(err)

			for _, subsystem := range port.Subsystems {
				pathSymlink := filepath.Join(subsystemPath, subsystem.Name)
				pathTarget := filepath.Join(portPath, port.Name, "subsystems", subsystem.Name)
				err = os.Symlink(pathSymlink, pathTarget)
				check(err)
			}
		}
	}

}
