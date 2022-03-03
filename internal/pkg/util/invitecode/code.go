package invitecode

const (
	base    = "E8S2DZX9WYLTN6BQF7CP5IK3MJUAR4HV"
	decimal = 32
	pad     = "G"
	length  = 6
)

func Gen(uid int) string {
	id := uid
	mod := 0
	res := ""
	for id != 0 {
		mod = id % decimal
		id = id / decimal
		res += string(base[mod])
	}
	resLen := len(res)
	if resLen < length {
		res += pad
		for i := 0; i < length-resLen-1; i++ {
			res += string(base[(uid+i)%decimal])
		}
	}
	return res
}
