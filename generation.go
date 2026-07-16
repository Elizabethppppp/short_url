package main

func toBase62(num uint64) string {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

	if num == 0 {
		return "0000000"
	}

	var result []byte
	for num > 0 {
		s := num % 62
		result = append(result, alphabet[s])
		num = num / 62
	}

	if len(result) < 7 {
		result = append(result, alphabet[0])
	}

	for i := 0; i < len(result)/2; i++ {
		result[i], result[len(result)-1-i] = result[len(result)-1-i], result[i]
	}

	return string(result)
}

func (u *URLstore) generateShortURL() string {
	u.counter++

	return toBase62(u.counter)
}
