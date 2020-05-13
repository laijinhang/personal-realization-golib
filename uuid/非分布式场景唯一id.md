func Uid(data []byte) string {
	h := md5.New()
	return fmt.Sprintf("%d%d%s", time.Now().Unix(), rand.Intn(100000), hex.EncodeToString(h.Sum(data)))
}
