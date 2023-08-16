package ethtool

import (
	"errors"
	"fmt"
	"unsafe"
)

type EthtoolLinkSettings struct {
	Cmd                    uint32
	Speed                  uint32
	Duplex                 uint8
	Port                   uint8
	Phy_address            uint8
	Autoneg                uint8
	Mdio_support           uint8
	Eth_tp_mdix            uint8
	Eth_tp_mdix_ctrl       uint8
	Link_mode_masks_nwords int8
	Transceiver            uint8
	Reserved1              [3]uint8
	Reserved               [7]uint32
	Link_mode_masks        [ETHTOOL_LINK_MODE_MASK_MAX_KERNEL_NU32 * 3]uint32
}

type EthtoolLinkNegotiations struct {
	Supported     uint64
	Advertising   uint64
	LpAdvertising uint64
}

func (ecmd *EthtoolLinkSettings) CmdGet(intf string) error {
	e, err := NewEthtool()
	if err != nil {
		return err
	}
	defer e.Close()
	return e.CmdGetLinkSetting(ecmd, intf)
}

func (ecmd *EthtoolLinkSettings) ParseNegotiations() (*EthtoolLinkNegotiations, error) {
	if ecmd.Cmd != ETHTOOL_GLINKSETTINGS {
		return nil, errors.New("must CmdGet before parsing")
	}

	if ecmd.Link_mode_masks_nwords <= 0 {
		return nil, errors.New("nwords not set")
	}

	negotiations := &EthtoolLinkNegotiations{}

	multiplicator := int(ecmd.Link_mode_masks_nwords)

	offset := 0
	nego := uint64(ecmd.Link_mode_masks[offset+1])
	negotiations.Supported = (nego << 32) | uint64(ecmd.Link_mode_masks[offset])

	offset += multiplicator
	nego = uint64(ecmd.Link_mode_masks[offset+1])
	negotiations.Advertising = (nego << 32) | uint64(ecmd.Link_mode_masks[offset])

	offset += multiplicator
	nego = uint64(ecmd.Link_mode_masks[offset+1])
	negotiations.LpAdvertising = (nego << 32) | uint64(ecmd.Link_mode_masks[offset])

	return negotiations, nil
}

// CmdGetLinkSetting returns the interface settings in the receiver struct
// and returns speed
func (e *Ethtool) CmdGetLinkSetting(ecmd *EthtoolLinkSettings, intf string) error {
	ecmd.Cmd = ETHTOOL_GLINKSETTINGS
	err := e.ioctl(intf, uintptr(unsafe.Pointer(ecmd)))

	if err != nil {
		return fmt.Errorf("getKernelNwords: %s", err)
	}

	ecmd.Link_mode_masks_nwords *= -1

	err = e.ioctl(intf, uintptr(unsafe.Pointer(ecmd)))
	if err != nil {
		return fmt.Errorf("getEthtoolLinkSettings: %s", err)
	}

	return nil
}

// CmdGetLinkSetting returns the interface settings in the receiver struct
// and returns speed
func (e *Ethtool) CmdSetLinkSetting(ecmd *EthtoolLinkSettings, intf string) error {
	if ecmd.Cmd != ETHTOOL_GLINKSETTINGS {
		return errors.New("load data before saving")
	}

	if ecmd.Link_mode_masks_nwords <= 0 {
		return errors.New("nwords must be filled with CmdGetLinkSetting")
	}

	ecmd.Cmd = ETHTOOL_SLINKSETTINGS
	err := e.ioctl(intf, uintptr(unsafe.Pointer(ecmd)))

	if err != nil {
		return fmt.Errorf("set link settings: %s", err)
	}

	return nil
}
