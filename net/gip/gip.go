package gip

import (
	"encoding/binary"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/basicfu/gf/text/gregex"
)

func RealIP(r *http.Request) string {
	if contextIp := r.Context().Value("remote_addr"); contextIp != nil {
		return contextIp.(string)
	}
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		return realIP
	}
	if fw := r.Header.Get("X-Forwarded-For"); fw != "" {
		str := fw
		if i := strings.IndexByte(str, ','); i >= 0 {
			str = str[:i]
		}
		if i := strings.LastIndexByte(str, ':'); i >= 0 {
			str = str[i+1:] //截取本地ip情况::ffff:10.0.0.2
		}
		return str
	}
	host, _, _ := net.SplitHostPort(r.RemoteAddr)
	return host
}

func HasLocalIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip.IsLoopback() {
		return true
	}
	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}
	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}
func LocalIp() []string {
	ips := []string{}
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, address := range addr {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips = append(ips, ipnet.IP.String())
			}
		}
	}
	return ips
}
func Ip2Long[T ~string | net.IP](ip T) int64 {
	var netIp net.IP
	switch v := any(ip).(type) {
	case string:
		netIp = net.ParseIP(v)
	case net.IP:
		netIp = v
	}
	if netIp == nil {
		return 0
	}
	return int64(binary.BigEndian.Uint32(netIp.To4()))
}
func Long2ip(long int64) net.IP {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, uint32(long))
	return net.IP(ipByte)
}
func Long2ipStr(long int64) string {
	ipByte := make([]byte, 4)
	binary.BigEndian.PutUint32(ipByte, uint32(long))
	return net.IP(ipByte).String()
}
func Copy(ip net.IP) net.IP {
	return append(net.IP(nil), ip.To4()...)
}
func Validate(ip string) bool {
	return gregex.IsMatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$`, ip)
}

type Cidr struct {
	CidrIpRange string
	Min         string
	Max         string
	Netmask     string
	Count       string
}

func NewCidr(ipRange string) *Cidr {
	tmp := Cidr{CidrIpRange: ipRange}
	return &tmp
}

func (c *Cidr) GetCidrIpRange() *Cidr {
	ip := strings.Split(c.CidrIpRange, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen := c.GetMaskLen()
	seg3MinIp, seg3MaxIp := c.GetIpSeg3Range(ipSegs, maskLen)
	seg4MinIp, seg4MaxIp := c.GetIpSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."

	c.Min = ipPrefix + strconv.Itoa(seg3MinIp) + "." + strconv.Itoa(seg4MinIp)
	c.Max = ipPrefix + strconv.Itoa(seg3MaxIp) + "." + strconv.Itoa(seg4MaxIp)
	return c
}

// 计算得到CIDR地址范围内可拥有的主机数量
func (c *Cidr) GetCidrHostNum() *Cidr {
	cidrIpNum := uint(0)
	var i uint = uint(32 - c.GetMaskLen() - 1)
	for ; i >= 1; i-- {
		cidrIpNum += 1 << i
	}
	c.Count = fmt.Sprintf("%d", cidrIpNum)
	return c
}

func (c *Cidr) GetMaskLen() int {
	maskLen, _ := strconv.Atoi(strings.Split(c.CidrIpRange, "/")[1])
	return maskLen
}

// 获取Cidr的掩码
func (c *Cidr) GetCidrIpMask() *Cidr {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-c.GetMaskLen())
	fmt.Println(fmt.Sprintf("%b \n", cidrMask))
	//计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))

	c.Netmask = fmt.Sprint(cidrMaskSeg1) + "." + fmt.Sprint(cidrMaskSeg2) + "." + fmt.Sprint(cidrMaskSeg3) + "." + fmt.Sprint(cidrMaskSeg4)
	return c
}

// 得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func (c *Cidr) GetIpSeg3Range(ipSegs []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIp, _ := strconv.Atoi(ipSegs[2])
		return segIp, segIp
	}
	ipSeg, _ := strconv.Atoi(ipSegs[2])
	return c.GetIpSegRange(uint8(ipSeg), uint8(24-maskLen))
}

// 得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func (c *Cidr) GetIpSeg4Range(ipSegs []string, maskLen int) (int, int) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	segMinIp, segMaxIp := c.GetIpSegRange(uint8(ipSeg), uint8(32-maskLen))
	return segMinIp + 1, segMaxIp
}

// 根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func (c *Cidr) GetIpSegRange(userSegIp, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIp := ipSegMax << offset
	segMinIp := netSegIp & userSegIp
	segMaxIp := userSegIp&(255<<offset) | ^(255 << offset)
	return int(segMinIp), int(segMaxIp)
}
