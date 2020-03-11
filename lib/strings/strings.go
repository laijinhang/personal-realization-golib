package strings

func ToUpper(s string) string {
	buf := []byte(s)
	for i := 0;i < len(buf);i++ {
		if buf[i] >= 'a' && buf[i] <= 'z' {
			buf[i] = byte('a' - 'A')
		}
	}
	return string(buf)
}

func ToLower(s string) string {
	buf := []byte(s)
	for i := 0;i < len(s);i++ {
		if buf[i] >= 'A' && buf[i] <= 'Z' {
			buf[i] += 'a' - 'A'
		}
	}
	return string(buf)
}

func ToTitle(s string) string {
	return ToUpper(s)
}

