package main

import "context"

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

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

func (u *URLstore) generateShortURL(ctx context.Context) (string, uint64, error) {
	var currentCount uint64
	err := u.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(last_counter), 100000000000) FROM url_schema.url").Scan(&currentCount)

	if err != nil {
		return "", 0, err
	}

	newCounter := currentCount + 1

	return toBase62(newCounter), newCounter, nil
}
