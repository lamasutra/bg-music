package client

import "hash/crc32"

func MapIdentity(host string) (uint32, error) {
	img, err := GetDataFromUrl(host + "map.img?gen=3")
	if err != nil {
		return 0, err
	}

	return crc32.ChecksumIEEE(img), nil
}
