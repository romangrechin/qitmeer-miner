package common

import (
	"github.com/Qitmeer/go-opencl/cl"
)

var DevicesTypesForGPUMining = cl.DeviceTypeGPU
var DevicesTypesForCPUMining = cl.DeviceTypeCPU
func GetDevices(t cl.DeviceType) []*cl.Device {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		MinerLoger.Error("Get Graphics card platforms error,please check!","error",err.Error())
		return nil
	}
	clDevices := make([]*cl.Device, 0)
	i := 0
	for _, platform := range platforms {
		platormDevices, err := cl.GetDevices(platform, t)
		if err != nil {
			MinerLoger.Error("Get Devices Error","platform",platform.Name(),"Error",err.Error())
			continue
		}
		for _, device := range platormDevices {
			clDevices = append(clDevices, device)
			MinerLoger.Info("Found Device","platform",platform.Name(),"minerID",i,"deviceName",device.Name(),
				"MaxWorkGroupSize(GB)",device.MaxWorkGroupSize(),"MaxMemAllocSize(GB)",float64(device.MaxMemAllocSize())/1024.00/1024.00,float64(device.GlobalMemSize())/1024.00/1024.00 )

			i++
		}
	}
	if len(clDevices) < 1{
		MinerLoger.Error("Don't had devices to mining,please check your PC!")
		return nil
	}
	return clDevices
}