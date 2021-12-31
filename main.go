package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v2"
)

// Nvmf struct
type Nvmf struct {
	Ports      []NvmfPort      `yaml:"ports"`
	Subsystems []NvmfSubsystem `yaml:"subsystems"`
}

// NvmfPort struct
type NvmfPort struct {
	Name        string                `yaml:"name"`
	AddrAdrfam  string                `yaml:"addr_adrfam"`
	AddrTraddr  string                `yaml:"addr_traddr"`
	AddrTrsvcid int                   `yaml:"addr_trsvcid"`
	AddrTrtype  string                `yaml:"addr_trtype"`
	Subsystems  []NvmfPortsSubsystems `yaml:"subsystems"`
}

// NvmfPortsSubsystems struct
type NvmfPortsSubsystems struct {
	Name string `yaml:"name"`
}

// NvmfSubsystem struct
type NvmfSubsystem struct {
	Name             string          `yaml:"name"`
	AttrAllowAnyHost int             `yaml:"attr_allow_any_host"`
	Namespaces       []NvmfNamespace `yaml:"namespaces"`
}

// NvmfNamespace struct
type NvmfNamespace struct {
	Name       int    `yaml:"name"`
	Enable     int    `yaml:"enable" default:"1"`
	DevicePath string `yaml:"device_path"`
	DeviceUUID string `yaml:"device_uuid"`
}

func (subsystem NvmfSubsystem) createSubsystem() error {
	path := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystem.Name)
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

func (subsystem NvmfSubsystem) writeAllowAnyHost() error {
	err := subsystem.writeSubsystemAttr("attr_allow_any_host", "1")
	return err
}

func (subsystem NvmfSubsystem) writeSubsystemAttr(attr string, value string) error {
	path := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystem.Name, attr)
	err := os.WriteFile(path, []byte(value), os.ModePerm)
	return err
}

func (namespace NvmfNamespace) createNamespace(subsystem NvmfSubsystem) error {
	path := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystem.Name, "namespaces", strconv.Itoa(namespace.Name))
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

func (namespace NvmfNamespace) writeNamespaceAttr(subsystem NvmfSubsystem, attr string, value string) error {
	path := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystem.Name, "namespaces", strconv.Itoa(namespace.Name), attr)
	err := os.WriteFile(path, []byte(value), os.ModePerm)
	return err
}

func (port NvmfPort) createPort() error {
	path := filepath.Join("/sys/kernel/config/nvmet/ports", port.Name)
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

func (port NvmfPort) writePortAttr(attr string, value string) error {
	path := filepath.Join("/sys/kernel/config/nvmet/ports", port.Name, attr)
	err := os.WriteFile(path, []byte(value), os.ModePerm)
	return err
}

func (port NvmfPort) bindSubsystem(subsystemName string) error {
	pathSymlink := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystemName)
	pathTarget := filepath.Join("/sys/kernel/config/nvmet/ports", port.Name, "subsystems", subsystemName)
	err := os.Symlink(pathSymlink, pathTarget)
	return err
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// ImportFromFile reads nvmet state from file and set it to system
func ImportFromFile(file string) {
	dat, err := os.ReadFile(file)
	check(err)

	nvmf := Nvmf{}
	err = yaml.Unmarshal([]byte(dat), &nvmf)
	check(err)

	for _, subsystem := range nvmf.Subsystems {

		err = subsystem.createSubsystem()
		check(err)

		err = subsystem.writeAllowAnyHost()
		check(err)

		for _, namespace := range subsystem.Namespaces {
			err := namespace.createNamespace(subsystem)
			check(err)

			err = namespace.writeNamespaceAttr(subsystem, "device_path", namespace.DevicePath)
			check(err)

			err = namespace.writeNamespaceAttr(subsystem, "device_uuid", namespace.DeviceUUID)
			check(err)

			err = namespace.writeNamespaceAttr(subsystem, "enable", strconv.Itoa(namespace.Enable))
			check(err)
		}
	}

	for _, port := range nvmf.Ports {
		err = port.createPort()
		check(err)

		err = port.writePortAttr("addr_adrfam", port.AddrAdrfam)
		check(err)

		err = port.writePortAttr("addr_traddr", port.AddrTraddr)
		check(err)

		err = port.writePortAttr("addr_trsvcid", strconv.Itoa(port.AddrTrsvcid))
		check(err)

		err = port.writePortAttr("addr_trtype", port.AddrTrtype)
		check(err)

		for _, subsystem := range port.Subsystems {
			err = port.bindSubsystem(subsystem.Name)
			check(err)
		}
	}
}

func main() {
	_, err := os.Stat("/sys/kernel/config/nvmet")
	check(err)

	args := os.Args[1:]
	switch args[0] {
	case "import":
		ImportFromFile(args[1])
	default:
		fmt.Println("invalid command")
		os.Exit(1)
	}
}
