package type_convert

import (
	"fmt"
	"strconv"
	"strings"
)

//将字符串形式的形如192.168.1.1的IP转为数字形式
func StringToIPv4(ipstr string) (u32 uint32, err error) {
	segments := strings.SplitN(ipstr, ".", 4)
	if len(segments) != 4 {
		err = ErrBadFormat
		return
	}
	var u81, u82, u83, u84 uint64
	u81, err = strconv.ParseUint(segments[0], 10, 32)
	if err != nil {
		return
	}
	u82, err = strconv.ParseUint(segments[1], 10, 32)
	if err != nil {
		return
	}
	u83, err = strconv.ParseUint(segments[2], 10, 32)
	if err != nil {
		return
	}
	u84, err = strconv.ParseUint(segments[3], 10, 32)
	if err != nil {
		return
	}
	if u81 > 255 || u82 > 255 || u83 > 255 || u84 > 255 {
		err = ErrBadFormat
		return
	}
	u32 |= uint32(u81) << 24
	u32 |= (uint32(u82) << 16) & 0x00ff0000
	u32 |= (uint32(u83) << 8) & 0x0000ff00
	u32 |= uint32(u84) & 0x000000ff
	return u32, nil
}

/**
将32位IPv4的数字形式IP转为字符串表示
*/
func IPv4ToString(u32 uint32) string {
	var u81, u82, u83, u84 uint8
	u81 = uint8(u32 >> 24)
	u82 = uint8((u32 & 0x00ff0000) >> 16)
	u83 = uint8((u32 & 0x0000ff00) >> 8)
	u84 = uint8(u32 & 0x000000ff)
	return fmt.Sprintf("%d.%d.%d.%d", u81, u82, u83, u84)
}
