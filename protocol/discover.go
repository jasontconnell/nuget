package protocol

func GetClient(url string) Client {
	v2client := NewV2Client(url)
	if v2client.IsValid() {
		return v2client
	}

	v3client := NewV3Client(url)
	if v3client.IsValid() {
		return v3client
	}

	return nil
}