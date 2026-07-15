package main

func toBase62(num uint64) string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if num == 0 {
		return "0000000"
	}

	var result []byte
	for num > 0 {
		s := num % 62
		result = append([]byte{alphabet[s]}, result...)
		num = num / 62
	}

	if len(result) < 7 {
		result = append([]byte{0}, result...)
	}

	return string(result)
}

func (u *URLstore) generateShortURL() string {
	u.counter++

	return toBase62(u.counter)
}
