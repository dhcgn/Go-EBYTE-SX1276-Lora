package commands

var (
	ReadingOperatingParameters = []byte{0xC1, 0xC1, 0xC1}
	ReadingVersionNumber       = []byte{0xC3, 0xC3, 0xC3}

	MessagerBroadcastHeader = []byte{0xFF, 0xFF, 0x06}
	MessagerHeader          = []byte{0x00, 0x00, 0x06}
)
