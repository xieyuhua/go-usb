package relay

import (
	"fmt"
	"github.com/karalabe/hid"
	"runtime"
)

const (
	OFF = iota
	ON
)

type IoStatus int

const (
	C1 = iota + 1
	C2
	C3
	C4
	C5
	C6
	C7
	C8
	ALL
)

type ChannelNumber int

type ChannelStatus struct {
	Channel_1 IoStatus
	Channel_2 IoStatus
	Channel_3 IoStatus
	Channel_4 IoStatus
	Channel_5 IoStatus
	Channel_6 IoStatus
	Channel_7 IoStatus
	Channel_8 IoStatus
}

type Relay struct {
	info *hid.DeviceInfo
	dev  *hid.Device
}

func List() []*Relay {
	list := make([]*Relay, 0, 5)
	// 枚举所有VID为0x16C0，PID为0x05DF的HID设备  
	// relayInfos := hid.Enumerate(0x16C0, 0x05DF)
	relayInfos := hid.Enumerate(0x0, 0x0)
	if len(relayInfos) <= 0 {
		return list
	}
    // 遍历设备  
	for i, info := range relayInfos {
		//遍历设备  
		fmt.Printf("Vendor ID: 0x%04X, Product ID: 0x%04X, Path: %s, Serial：%s, Interface:%d \n", info.VendorID, info.ProductID, info.Path, info.Serial, info.Interface)  
        // 你可以根据需要进一步处理每个设备，比如打开它以读取或写入数据  
		relay, err := info.Open()
		if err == nil {
			list = append(list, &Relay{info: &relayInfos[i]})
			relay.Close()
		}
	}
	return list
}

func (this *Relay) Open() error {
	dev, err := this.info.Open()
	if err == nil {
		this.dev = dev
	}
	return err
}

func (this *Relay) Close() error {
	return this.dev.Close()
}

func (this *Relay) setIO(s IoStatus, no ChannelNumber) error {
	cmd := make([]byte, 9)
	cmd[0] = 0x0
	if no < C1 && no > ALL {
		return fmt.Errorf("channel number (%d) is illegal", no)
	}

	if no == ALL {
		if s == ON {
			cmd[1] = 0xFE
		} else {
			cmd[1] = 0xFC
		}
	} else {
		if s == ON {
			cmd[1] = 0xFF
		} else {
			cmd[1] = 0xFD
		}
		cmd[2] = byte(no)
	}

	_, err := this.dev.Write(cmd)
	return err
}

func (this *Relay) GetStatus() (*ChannelStatus, error) {
	cmd := make([]byte, 9)
	_, err := this.dev.Read(cmd)
	if err != nil {
		return nil, err
	}

	// Remove HID report ID on Windows, others OSes don't need it.
	if runtime.GOOS == "windows" {
		cmd = cmd[1:]
	}

	status := &ChannelStatus{}
	status.Channel_1 = IoStatus(cmd[7] >> 0 & 0x01)
	status.Channel_2 = IoStatus(cmd[7] >> 1 & 0x01)
	status.Channel_3 = IoStatus(cmd[7] >> 2 & 0x01)
	status.Channel_4 = IoStatus(cmd[7] >> 3 & 0x01)
	status.Channel_5 = IoStatus(cmd[7] >> 4 & 0x01)
	status.Channel_6 = IoStatus(cmd[7] >> 5 & 0x01)
	status.Channel_7 = IoStatus(cmd[7] >> 6 & 0x01)
	status.Channel_8 = IoStatus(cmd[7] >> 7 & 0x01)
	return status, err
}

func (this *Relay) TurnOn(num ChannelNumber) error {
	return this.setIO(ON, num)
}

func (this *Relay) TurnOff(num ChannelNumber) error {
	return this.setIO(OFF, num)
}

func (this *Relay) TurnAllOn() error {
	return this.setIO(ON, ALL)
}

func (this *Relay) TurnAllOff() error {
	return this.setIO(OFF, ALL)
}

func (this *Relay) SetSN(sn string) error {
	if len(sn) > 5 {
		return fmt.Errorf("The length of `%s` is large than 5 bytes.", sn)
	}
	cmd := make([]byte, 9)
	cmd[0] = 0x00
	cmd[1] = 0xFA
	copy(cmd[2:], sn)
	_, err := this.dev.Write(cmd)
	if err != nil {
		return err
	}
	return err
}

func (this *Relay) GetSN() (string, error) {
	cmd := make([]byte, 9)
	_, err := this.dev.Read(cmd)
	var sn string
	if err != nil {
		sn = ""
	} else {
		// Remove HID report ID on Windows, others OSes don't need it.
		if runtime.GOOS == "windows" {
			cmd = cmd[1:]
		}
		sn = string(cmd[:5])
	}
	return sn, err
}
