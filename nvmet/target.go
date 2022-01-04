package nvmet

import (
	"fmt"
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

// Create creates subsystem in nvme target
func (subsystem *NvmfSubsystem) Create() error {
	path := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystem.Name)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create subsystem %v, %v", subsystem, err)
	}

	return err
}

// WriteAllowAnyHost configures subsystem to allow any host to connect
func (subsystem *NvmfSubsystem) WriteAllowAnyHost() error {
	return subsystem.WriteAttr("attr_allow_any_host", "1")
}

// WriteAttr writes subsystems attirube
func (subsystem *NvmfSubsystem) WriteAttr(attr string, value string) error {
	path := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystem.Name, attr)
	err := os.WriteFile(path, []byte(value), os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not write subsystem attr %v, subsystem %v, %v", attr, subsystem.Name, err)
	}

	return err
}

// Create creates namespace in specified subsystem
func (namespace *NvmfNamespace) Create(subsystem NvmfSubsystem) error {
	path := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystem.Name, "namespaces", strconv.Itoa(namespace.Name))
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create namespace %v, subsystem %v, %v", namespace.Name, subsystem.Name, err)
	}

	return err
}

// WriteAttr writes namespace attirube in specified subsystem
func (namespace *NvmfNamespace) WriteAttr(subsystem NvmfSubsystem, attr string, value string) error {
	path := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystem.Name, "namespaces", strconv.Itoa(namespace.Name), attr)
	err := os.WriteFile(path, []byte(value), os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not write namespace attr %v, subsystem %v, namespace %v, %v", attr, subsystem.Name, namespace.Name, err)
	}

	return err
}

// Create creates port
func (port *NvmfPort) Create() error {
	path := filepath.Join("/sys/kernel/config/nvmet/ports", port.Name)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create port %v, %v", port.Name, err)
	}

	return err
}

// WriteAttr write port attribute
func (port *NvmfPort) WriteAttr(attr string, value string) error {
	path := filepath.Join("/sys/kernel/config/nvmet/ports", port.Name, attr)
	err := os.WriteFile(path, []byte(value), os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not write port attr %v, port %v, %v", attr, port.Name, err)
	}

	return err
}

// Bind binds specified subsystem to port
func (port *NvmfPort) Bind(subsystemName string) error {
	pathSymlink := filepath.Join("/sys/kernel/config/nvmet/subsystems", subsystemName)
	pathTarget := filepath.Join("/sys/kernel/config/nvmet/ports", port.Name, "subsystems", subsystemName)
	err := os.Symlink(pathSymlink, pathTarget)
	if err != nil {
		return fmt.Errorf("could not bind subsystem %v to port %v", subsystemName, port.Name)
	}

	return err
}

// ReadFile reads yaml file into struct
func (nvmf *Nvmf) ReadFile(file string) error {
	dat, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("could not read from file, %v", err)
	}
	if err = yaml.Unmarshal([]byte(dat), &nvmf); err != nil {
		return fmt.Errorf("could not unmarshal data, %v", err)
	}

	return err
}

// Apply runs through struct to apply it to target
func (nvmf *Nvmf) Apply() error {
	var err error

	for _, subsystem := range nvmf.Subsystems {
		if err := subsystem.Create(); err != nil {
			return err
		}

		if err := subsystem.WriteAllowAnyHost(); err != nil {
			return err
		}

		for _, namespace := range subsystem.Namespaces {
			if err := namespace.Create(subsystem); err != nil {
				return err
			}

			if err := namespace.WriteAttr(subsystem, "device_path", namespace.DevicePath); err != nil {
				return err
			}

			if err := namespace.WriteAttr(subsystem, "device_uuid", namespace.DeviceUUID); err != nil {
				return err
			}

			if err := namespace.WriteAttr(subsystem, "enable", strconv.Itoa(namespace.Enable)); err != nil {
				return err
			}
		}
	}
	for _, port := range nvmf.Ports {
		if err := port.Create(); err != nil {
			return err
		}

		if err := port.WriteAttr("addr_adrfam", port.AddrAdrfam); err != nil {
			return err
		}

		if err := port.WriteAttr("addr_traddr", port.AddrTraddr); err != nil {
			return err
		}

		if err := port.WriteAttr("addr_trsvcid", strconv.Itoa(port.AddrTrsvcid)); err != nil {
			return err
		}

		if err := port.WriteAttr("addr_trtype", port.AddrTrtype); err != nil {
			return err
		}

		for _, subsystem := range port.Subsystems {
			if err := port.Bind(subsystem.Name); err != nil {
				return err
			}
		}
	}

	return err
}
