package icc

import "time"

type Header struct {
	ProfileSize            uint32
	PreferredCMM           Signature
	Version                Version
	DeviceClass            DeviceClass
	DataColorSpace         ColorSpace
	ProfileConnectionSpace ColorSpace
	CreatedAt              time.Time
	PrimaryPlatform        PrimaryPlatform
	Embedded               bool
	DependsOnEmbeddedData  bool
	DeviceManufacturer     Signature
	DeviceModel            Signature
	DeviceAttributes       uint64
	RenderingIntent        RenderingIntent
	PCSIlluminant          [3]uint32
	ProfileCreator         Signature
	ProfileID              [16]byte
}
