// Copyright (c) 2014 Kohei YOSHIDA. All rights reserved.
// This software is licensed under the 3-Clause BSD License
// that can be found in LICENSE file.

package arib

type DescriptorTag uint8

const (
	ConditionalAccessDescriptor DescriptorTag = 0x09
	CopyrightDescriptor         DescriptorTag = 0x0D

	CarouselIdentifierDescriptor      DescriptorTag = 0x13
	AssociationTagDescriptor          DescriptorTag = 0x14
	DeferredAssociationTagsDescriptor DescriptorTag = 0x15

	AVCVideoDescriptor        DescriptorTag = 0x28
	AVCTimingAndHRDDescriptor DescriptorTag = 0x2A

	NetworkNameDescriptor               DescriptorTag = 0x40
	ServiceListDescriptor               DescriptorTag = 0x41
	StuffingDescriptor                  DescriptorTag = 0x42
	SatelliteDeliverySystemDescriptor   DescriptorTag = 0x43
	TerrestrialDeliverySystemDescriptor DescriptorTag = 0x44
	BouquetNameDescriptor               DescriptorTag = 0x47
	ServiceDescriptor                   DescriptorTag = 0x48
	CountryAvailabilityDescriptor       DescriptorTag = 0x49
	LinkageDescriptor                   DescriptorTag = 0x4A
	NVODReferenceDescriptor             DescriptorTag = 0x4B
	TimeShiftedServiceDescriptor        DescriptorTag = 0x4C
	ShortEventDescriptor                DescriptorTag = 0x4D
	ExtendedEventDescriptor             DescriptorTag = 0x4E
	TimeShiftedEventDescriptor          DescriptorTag = 0x4F

	ComponentDescriptor        DescriptorTag = 0x50
	MosaicDescriptor           DescriptorTag = 0x51
	StreamIdentifierDescriptor DescriptorTag = 0x52
	CAIdentifierDescriptor     DescriptorTag = 0x53
	ContentDescriptor          DescriptorTag = 0x54
	ParentalRatingDescriptor   DescriptorTag = 0x55
	LocalTimeOffsetDescriptor  DescriptorTag = 0x58

	PartialTransportStreamDescriptor DescriptorTag = 0x63

	HierarchicalTransmissionDescriptor   DescriptorTag = 0xC0
	DigitalCopyControlDescriptor         DescriptorTag = 0xC1
	NetworkIdentificationDescriptor      DescriptorTag = 0xC2
	PartialTransportStreamTimeDescriptor DescriptorTag = 0xC3
	AudioComponentDescriptor             DescriptorTag = 0xC4
	HyperlinkDescriptor                  DescriptorTag = 0xC5
	TargetRegionDescriptor               DescriptorTag = 0xC6
	DataContentDescriptor                DescriptorTag = 0xC7
	VideoDecodeControlDescriptor         DescriptorTag = 0xC8
	DownloadContentDescriptor            DescriptorTag = 0xC9
	CA_EMM_TSDescriptor                  DescriptorTag = 0xCA
	CAContractInformationDescriptor      DescriptorTag = 0xCB
	CAServiceDescriptor                  DescriptorTag = 0xCC
	TSInformationDescriptor              DescriptorTag = 0xCD
	ExtendedBroadcasterDescriptor        DescriptorTag = 0xCE
	LogoTransmissionDescriptor           DescriptorTag = 0xCF

	BasicLocalEventDesciptor        DescriptorTag = 0xD0
	ReferenceDescriptor             DescriptorTag = 0xD1
	NodeRelationDescriptor          DescriptorTag = 0xD2
	ShortNodeInformationDescriptor  DescriptorTag = 0xD3
	STCReferenceDescriptor          DescriptorTag = 0xD4
	SeriesDescriptor                DescriptorTag = 0xD5
	EventGroupDescriptor            DescriptorTag = 0xD6
	SIParameterDescriptor           DescriptorTag = 0xD7
	BroadcasterNameDescriptor       DescriptorTag = 0xD8
	ComponentGroupDescriptor        DescriptorTag = 0xD9
	SIPrimeTSDescriptor             DescriptorTag = 0xDA
	BoardInformationDescriptor      DescriptorTag = 0xDB
	LDTLinkageDescriptor            DescriptorTag = 0xDC
	ConnectedTransmissionDescriptor DescriptorTag = 0xDD
	ContentAvailabilityDescriptor   DescriptorTag = 0xDE

	ServiceGroupDescriptor DescriptorTag = 0xE0

	CarouselCompatibleCompositeDescriptor DescriptorTag = 0xF7
	ConditionalPlaybackDescriptor         DescriptorTag = 0xF8
	PartialReceptionDescriptor            DescriptorTag = 0xFB
	EmergencyInformationDescriptor        DescriptorTag = 0xFC
	SystemManagementDescriptor            DescriptorTag = 0xFE
)
