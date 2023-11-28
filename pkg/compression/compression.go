package compression

import (
	"bytes"
	"errors"

	"github.com/pierrec/lz4/v4"
)

// CompressData compresses the given data using LZ4 algorithm.
func CompressData(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	lz4Writer := lz4.NewWriter(&buffer)

	// You can set various options on lz4Writer if needed
	_, err := lz4Writer.Write(data)
	if err != nil {
		return nil, err
	}
	err = lz4Writer.Close()
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

// DecompressData decompresses LZ4 compressed data.
func DecompressData(compressedData []byte, originalSize int) ([]byte, error) {
	if originalSize <= 0 {
		return nil, errors.New("original size must be positive")
	}

	decompressedData := make([]byte, originalSize)
	lz4Reader := lz4.NewReader(bytes.NewReader(compressedData))

	n, err := lz4Reader.Read(decompressedData)
	if err != nil {
		return nil, err
	}

	if n != originalSize {
		return nil, errors.New("decompressed data size does not match original size")
	}

	return decompressedData, nil
}
