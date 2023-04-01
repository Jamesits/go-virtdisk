package ffi

import "golang.org/x/sys/windows"

const MaxPath = 255

type DeviceType uint32

const (
	FileDeviceBeep               DeviceType = 0x00000001
	FileDeviceCdRom              DeviceType = 0x00000002
	FileDeviceCdRomFileSystem    DeviceType = 0x00000003
	FileDeviceController         DeviceType = 0x00000004
	FileDeviceDatalink           DeviceType = 0x00000005
	FileDeviceDfs                DeviceType = 0x00000006
	FileDeviceDisk               DeviceType = 0x00000007
	FileDeviceDiskFileSystem     DeviceType = 0x00000008
	FileDeviceFileSystem         DeviceType = 0x00000009
	FileDeviceInportPort         DeviceType = 0x0000000a
	FileDeviceKeyboard           DeviceType = 0x0000000b
	FileDeviceMailslot           DeviceType = 0x0000000c
	FileDeviceMidiIn             DeviceType = 0x0000000d
	FileDeviceMidiOut            DeviceType = 0x0000000e
	FileDeviceMouse              DeviceType = 0x0000000f
	FileDeviceMultiUncProvider   DeviceType = 0x00000010
	FileDeviceNamedPipe          DeviceType = 0x00000011
	FileDeviceNetwork            DeviceType = 0x00000012
	FileDeviceNetworkBrowser     DeviceType = 0x00000013
	FileDeviceNetworkFileSystem  DeviceType = 0x00000014
	FileDeviceNull               DeviceType = 0x00000015
	FileDeviceParallelPort       DeviceType = 0x00000016
	FileDevicePhysicalNetcard    DeviceType = 0x00000017
	FileDevicePrinter            DeviceType = 0x00000018
	FileDeviceScanner            DeviceType = 0x00000019
	FileDeviceSerialMousePort    DeviceType = 0x0000001a
	FileDeviceSerialPort         DeviceType = 0x0000001b
	FileDeviceScreen             DeviceType = 0x0000001c
	FileDeviceSound              DeviceType = 0x0000001d
	FileDeviceStreams            DeviceType = 0x0000001e
	FileDeviceTape               DeviceType = 0x0000001f
	FileDeviceTapeFileSystem     DeviceType = 0x00000020
	FileDeviceTransport          DeviceType = 0x00000021
	FileDeviceUnknown            DeviceType = 0x00000022
	FileDeviceVideo              DeviceType = 0x00000023
	FileDeviceVirtualDisk        DeviceType = 0x00000024
	FileDeviceWaveIn             DeviceType = 0x00000025
	FileDeviceWaveOut            DeviceType = 0x00000026
	FileDevice8042Port           DeviceType = 0x00000027
	FileDeviceNetworkRedirector  DeviceType = 0x00000028
	FileDeviceBattery            DeviceType = 0x00000029
	FileDeviceBusExtender        DeviceType = 0x0000002a
	FileDeviceModem              DeviceType = 0x0000002b
	FileDeviceVdm                DeviceType = 0x0000002c
	FileDeviceMassStorage        DeviceType = 0x0000002d
	FileDeviceSmb                DeviceType = 0x0000002e
	FileDeviceKs                 DeviceType = 0x0000002f
	FileDeviceChanger            DeviceType = 0x00000030
	FileDeviceSmartcard          DeviceType = 0x00000031
	FileDeviceAcpi               DeviceType = 0x00000032
	FileDeviceDvd                DeviceType = 0x00000033
	FileDeviceFullscreenVideo    DeviceType = 0x00000034
	FileDeviceDfsFileSystem      DeviceType = 0x00000035
	FileDeviceDfsVolume          DeviceType = 0x00000036
	FileDeviceSerenum            DeviceType = 0x00000037
	FileDeviceTermsrv            DeviceType = 0x00000038
	FileDeviceKsec               DeviceType = 0x00000039
	FileDeviceFips               DeviceType = 0x0000003A
	FileDeviceInfiniband         DeviceType = 0x0000003B
	FileDeviceVmbus              DeviceType = 0x0000003E
	FileDeviceCryptProvider      DeviceType = 0x0000003F
	FileDeviceWpd                DeviceType = 0x00000040
	FileDeviceBluetooth          DeviceType = 0x00000041
	FileDeviceMtComposite        DeviceType = 0x00000042
	FileDeviceMtTransport        DeviceType = 0x00000043
	FileDeviceBiometric          DeviceType = 0x00000044
	FileDevicePmi                DeviceType = 0x00000045
	FileDeviceEhstor             DeviceType = 0x00000046
	FileDeviceDevapi             DeviceType = 0x00000047
	FileDeviceGpio               DeviceType = 0x00000048
	FileDeviceUsbex              DeviceType = 0x00000049
	FileDeviceConsole            DeviceType = 0x00000050
	FileDeviceNfp                DeviceType = 0x00000051
	FileDeviceSysenv             DeviceType = 0x00000052
	FileDeviceVirtualBlock       DeviceType = 0x00000053
	FileDevicePointOfService     DeviceType = 0x00000054
	FileDeviceStorageReplication DeviceType = 0x00000055
	FileDeviceTrustEnv           DeviceType = 0x00000056
	FileDeviceUcm                DeviceType = 0x00000057
	FileDeviceUcmtcpci           DeviceType = 0x00000058
	FileDevicePersistentMemory   DeviceType = 0x00000059
	FileDeviceNvdimm             DeviceType = 0x0000005a
	FileDeviceHolographic        DeviceType = 0x0000005b
	FileDeviceSdfxhci            DeviceType = 0x0000005c
)

type StoragePropertyId uint32

// https://learn.microsoft.com/en-us/windows/win32/api/winioctl/ne-winioctl-storage_property_id
const (
	StorageDeviceProperty                                  StoragePropertyId = 0
	StorageAdapterProperty                                 StoragePropertyId = 1
	StorageDeviceIdProperty                                StoragePropertyId = 2
	StorageDeviceUniqueIdProperty                          StoragePropertyId = 3
	StorageDeviceWriteCacheProperty                        StoragePropertyId = 4
	StorageMiniportProperty                                StoragePropertyId = 5
	StorageAccessAlignmentProperty                         StoragePropertyId = 6
	StorageDeviceSeekPenaltyProperty                       StoragePropertyId = 7
	StorageDeviceTrimProperty                              StoragePropertyId = 8
	StorageDeviceWriteAggregationPropertyStoragePropertyId                   = 9
	StorageDeviceDeviceTelemetryProperty                   StoragePropertyId = 10 // 0xA
	StorageDeviceLBProvisioningProperty                    StoragePropertyId = 11 // 0xB
	StorageDevicePowerProperty                             StoragePropertyId = 12 // 0xC
	StorageDeviceCopyOffloadProperty                       StoragePropertyId = 13 // 0xD
	StorageDeviceResiliencyProperty                        StoragePropertyId = 14 // 0xE
	StorageDeviceMediumProductType                         StoragePropertyId = 15
	StorageAdapterRpmbProperty                             StoragePropertyId = 16
	StorageAdapterCryptoProperty                           StoragePropertyId = 17
	StorageDeviceIoCapabilityProperty                      StoragePropertyId = 48
	StorageAdapterProtocolSpecificProperty                                   = 49
	StorageDeviceProtocolSpecificPropertyStoragePropertyId                   = 50
	StorageAdapterTemperatureProperty                      StoragePropertyId = 51
	StorageDeviceTemperatureProperty                       StoragePropertyId = 52
	StorageAdapterPhysicalTopologyProperty                                   = 53
	StorageDevicePhysicalTopologyPropertyStoragePropertyId                   = 54
	StorageDeviceAttributesProperty                        StoragePropertyId = 55
	StorageDeviceManagementStatus                          StoragePropertyId = 56
	StorageAdapterSerialNumberProperty                     StoragePropertyId = 57
	StorageDeviceLocationProperty                          StoragePropertyId = 58
	StorageDeviceNumaProperty                              StoragePropertyId = 59
	StorageDeviceZonedDeviceProperty                       StoragePropertyId = 60
	StorageDeviceUnsafeShutdownCount                       StoragePropertyId = 61
	StorageDeviceEnduranceProperty                         StoragePropertyId = 62
	StorageDeviceSelfEncryptionProperty                    StoragePropertyId = 64
	StorageFruIdProperty                                   StoragePropertyId = 65
)

type StorageQueryType uint32

// https://learn.microsoft.com/en-us/windows/win32/api/winioctl/ne-winioctl-storage_query_type
const (
	PropertyStandardQuery StorageQueryType = 0
	PropertyExistsQuery   StorageQueryType = 1
	PropertyMaskQuery     StorageQueryType = 2
)

var GuidDevInterfaceDisk = &windows.GUID{
	Data1: 0x53F56307,
	Data2: 0xB6BF,
	Data3: 0x11D0,
	Data4: [8]byte{0x94, 0xF2, 0x00, 0xA0, 0xC9, 0x1E, 0xFB, 0x8B},
}
