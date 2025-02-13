package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/NVIDIA/go-dcgm/pkg/dcgm"
)

var (
	process = flag.Uint("pid", 0, "Provide pid to get this process information.")
)

type GPUDcgm struct {
	deviceGroupName   string
	deviceGroupHandle dcgm.GroupHandle
	cleanup           func()
}

func (d *GPUDcgm) Init() error {
	cleanup, err := dcgm.Init(dcgm.Embedded)
	if err != nil {
		if cleanup != nil {
			cleanup()
		}
		return fmt.Errorf("not able to connect to DCGM: %s", err)
	}
	d.cleanup = cleanup
	dcgm.FieldsInit()

	if err := d.createDeviceGroup(); err != nil {
		log.Printf("failed to create device group: %v\n", err)
		d.Shutdown()
		return err
	}

	if err := d.addDevicesToGroup(); err != nil {
		log.Printf("failed to add device group: %v\n", err)
		d.Shutdown()
		return err
	}

	if err := d.setupWatcher(); err != nil {
		log.Printf("failed to set up watcher: %v\n", err)
		d.Shutdown()
		return err
	}
	log.Println("DCGM initialized successfully")
	return nil
}
func (d *GPUDcgm) Shutdown() bool {
	log.Println("Shutting down DCGM")
	dcgm.FieldsTerm()
	if d.deviceGroupName != "" {
		dcgm.DestroyGroup(d.deviceGroupHandle)
	}
	if d.cleanup != nil {
		d.cleanup()
	}
	return true
}

func (d *GPUDcgm) createDeviceGroup() error {
	deviceGroupName := "dg-" + time.Now().Format("20060102150405")
	deviceGroup, err := dcgm.CreateGroup(deviceGroupName)
	if err != nil {
		return fmt.Errorf("failed to create group %q: %v", deviceGroupName, err)
	}
	d.deviceGroupName = deviceGroupName
	d.deviceGroupHandle = deviceGroup
	log.Printf("Created device group %q\n", deviceGroupName)
	return nil
}

func (d *GPUDcgm) addDevicesToGroup() error {
	supportedDeviceIndices, err := dcgm.GetSupportedDevices()
	if err != nil {
		return fmt.Errorf("failed to find supported devices: %v", err)
	}
	log.Printf("found %d supported devices\n", len(supportedDeviceIndices))
	for _, gpuIndex := range supportedDeviceIndices {
		log.Printf("Adding device %d to group %q\n", gpuIndex, d.deviceGroupName)
		err = dcgm.AddToGroup(d.deviceGroupHandle, gpuIndex)
		if err != nil {
			log.Printf("failed to add device %d to group %q: %v\n", gpuIndex, d.deviceGroupName, err)
		}
	}
	// add entity to the group
	hierarchy, err := dcgm.GetGpuInstanceHierarchy()
	if err != nil {
		d.Shutdown()
		return fmt.Errorf("failed to get gpu hierachy: %v", err)
	}

	if hierarchy.Count > 0 {
		// MIG is enabled
		for i := uint(0); i < hierarchy.Count; i++ {
			if hierarchy.EntityList[i].Parent.EntityGroupId == dcgm.FE_GPU {
				// add a GPU instance
				info := hierarchy.EntityList[i].Info
				entityId := hierarchy.EntityList[i].Entity.EntityId
				gpuId := hierarchy.EntityList[i].Parent.EntityId
				err = dcgm.AddEntityToGroup(d.deviceGroupHandle, dcgm.FE_GPU_I, entityId)
				log.Printf("Adding GPU ID[%v] MIG Device[%v] GPU Device Index[%v] Instance ID[%v]: err %v\n", gpuId, entityId, info.NvmlGpuIndex, info.NvmlInstanceId, err)
			}
		}
	}

	return nil
}

func (d *GPUDcgm) setupWatcher() error {
	err := dcgm.WatchPidFieldsWithGroup(d.deviceGroupHandle)
	if err != nil {
		return fmt.Errorf("failed to set up pid watcher, err %v", err)
	}
	log.Println("DCGM pid watcher set up successfully")
	return nil
}

func main() {
	dg := &GPUDcgm{}
	if err := dg.Init(); err != nil {
		log.Panicln(err)
	}
	defer dg.Shutdown()
	flag.Parse()
	// create a tick to watch the process every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// get the process info
		pidInfo, err := dcgm.GetProcessInfo(dg.deviceGroupHandle, *process)
		if err != nil {
			log.Printf("failed to get process info: %v\n", err)
		} else {
			for _, gpu := range pidInfo {
				log.Printf("gpu %d, process info: %+v\n", gpu.GPU, gpu)
			}
		}
	}
}
