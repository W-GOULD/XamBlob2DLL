// pkg/assemblystore/assemblystore.go

package assemblystore

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/W-GOULD/XamBlob2DLL/pkg/compression"
	"github.com/W-GOULD/XamBlob2DLL/pkg/manifest"
	"github.com/pierrec/lz4"
)

type AssemblyStore struct {
	Raw             []byte
	FileName        string
	ManifestEntries manifest.ManifestList
	HdrMagic        []byte
	HdrVersion      uint32
	HdrLec          uint32
	HdrGec          uint32
	HdrStoreID      uint32
	AssembliesList  []AssemblyStoreAssembly
	GlobalHash32    []HashEntry
	GlobalHash64    []HashEntry
}

type AssemblyStoreAssembly struct {
	DataOffset       uint32
	DataSize         uint32
	DebugDataOffset  uint32
	DebugDataSize    uint32
	ConfigDataOffset uint32
	ConfigDataSize   uint32
}

type HashEntry struct {
	HashVal         string
	MappingIndex    uint32
	LocalStoreIndex uint32
	StoreID         uint32
}

func NewAssemblyStore(fileName string, manifestEntries manifest.ManifestList) (*AssemblyStore, error) {
	as := &AssemblyStore{
		FileName:        filepath.Base(fileName),
		ManifestEntries: manifestEntries,
		AssembliesList:  make([]AssemblyStoreAssembly, 0),
		GlobalHash32:    make([]HashEntry, 0),
		GlobalHash64:    make([]HashEntry, 0),
	}

	err := as.ParseAndReadStore(fileName)
	if err != nil {
		return nil, err
	}

	return as, nil
}

// Method to parse and read store
func (as *AssemblyStore) ParseAndReadStore(fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// Safely read the entire file
	as.Raw, err = ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Parsing the header section
	as.HdrMagic = as.Raw[:4]
	if string(as.HdrMagic) != string(constants.AssemblyStoreMagic) {
		return errors.New("invalid magic value")
	}

	as.HdrVersion = binary.LittleEndian.Uint32(as.Raw[4:8])
	if as.HdrVersion > constants.AssemblyStoreFormatVersion {
		return errors.New("version is higher than expected")
	}

	as.HdrLec = binary.LittleEndian.Uint32(as.Raw[8:12])
	as.HdrGec = binary.LittleEndian.Uint32(as.Raw[12:16])
	as.HdrStoreID = binary.LittleEndian.Uint32(as.Raw[16:20])

	// Parsing assembly data
	offset := 20 // Starting position after the header
	for i := uint32(0); i < as.HdrLec; i++ {
		assembly := AssemblyStoreAssembly{
			DataOffset:       binary.LittleEndian.Uint32(as.Raw[offset : offset+4]),
			DataSize:         binary.LittleEndian.Uint32(as.Raw[offset+4 : offset+8]),
			DebugDataOffset:  binary.LittleEndian.Uint32(as.Raw[offset+8 : offset+12]),
			DebugDataSize:    binary.LittleEndian.Uint32(as.Raw[offset+12 : offset+16]),
			ConfigDataOffset: binary.LittleEndian.Uint32(as.Raw[offset+16 : offset+20]),
			ConfigDataSize:   binary.LittleEndian.Uint32(as.Raw[offset+20 : offset+24]),
		}
		as.AssembliesList = append(as.AssembliesList, assembly)
		offset += 24 // Move to the next assembly entry
	}

	// Parsing 32-bit hashes
	for i := uint32(0); i < as.HdrLec; i++ {
		hashEntry := HashEntry{
			HashVal:         fmt.Sprintf("0x%08x", binary.LittleEndian.Uint32(as.Raw[offset:offset+4])),
			MappingIndex:    binary.LittleEndian.Uint32(as.Raw[offset+8 : offset+12]),
			LocalStoreIndex: binary.LittleEndian.Uint32(as.Raw[offset+12 : offset+16]),
			StoreID:         binary.LittleEndian.Uint32(as.Raw[offset+16 : offset+20]),
		}
		as.GlobalHash32 = append(as.GlobalHash32, hashEntry)
		offset += 20 // Move to the next hash32 entry
	}

	// Parsing 64-bit hashes
	for i := uint32(0); i < as.HdrLec; i++ {
		hashEntry := HashEntry{
			HashVal:         fmt.Sprintf("0x%016x", binary.LittleEndian.Uint64(as.Raw[offset:offset+8])),
			MappingIndex:    binary.LittleEndian.Uint32(as.Raw[offset+8 : offset+12]),
			LocalStoreIndex: binary.LittleEndian.Uint32(as.Raw[offset+12 : offset+16]),
			StoreID:         binary.LittleEndian.Uint32(as.Raw[offset+16 : offset+20]),
		}
		as.GlobalHash64 = append(as.GlobalHash64, hashEntry)
		offset += 20 // Move to the next hash64 entry
	}

	return nil
}

// Method to extract all assemblies
func (as *AssemblyStore) ExtractAll(jsonConfig map[string]interface{}, outputPath string) error {
	storeJson := make(map[string]interface{})
	header := map[string]uint32{
		"version":  as.HdrVersion,
		"lec":      as.HdrLec,
		"gec":      as.HdrGec,
		"store_id": as.HdrStoreID,
	}
	storeJson[as.FileName] = map[string]interface{}{"header": header}
	jsonConfig["stores"] = append(jsonConfig["stores"].([]interface{}), storeJson)

	for i, assembly := range as.AssembliesList {
		assemblyDict := make(map[string]interface{})
		entry := as.ManifestEntries.GetIdx(as.HdrStoreID, i)

		if entry == nil {
			continue // or handle the missing entry appropriately
		}

		assemblyDict["name"] = entry.Name
		assemblyDict["store_id"] = entry.BlobID
		assemblyDict["blob_idx"] = entry.BlobIdx
		assemblyDict["hash32"] = entry.Hash32
		assemblyDict["hash64"] = entry.Hash64

		outFileName := fmt.Sprintf("%s/%s.dll", outputPath, entry.Name)
		assemblyDict["file"] = outFileName

		if err := os.MkdirAll(filepath.Dir(outFileName), 0755); err != nil {
			return err
		}

		fileData, err := as.extractAssemblyData(assembly)
		if err != nil {
			return err
		}

		if err := os.WriteFile(outFileName, fileData, 0644); err != nil {
			return err
		}

		// Update JSON config
		jsonConfig["assemblies"] = append(jsonConfig["assemblies"].([]interface{}), assemblyDict)
	}

	return nil
}

func (as *AssemblyStore) extractAssemblyData(assembly AssemblyStoreAssembly) ([]byte, error) {
	// Check if assembly data is valid
	if assembly.DataSize == 0 || int(assembly.DataOffset+assembly.DataSize) > len(as.Raw) {
		return nil, errors.New("invalid assembly data or offset out of bounds")
	}

	data := as.Raw[assembly.DataOffset : assembly.DataOffset+assembly.DataSize]

	// Check if data is compressed
	if isCompressed(data) {
		decompressedData, err := compression.DecompressLZ4(data)
		if err != nil {
			return nil, err
		}
		return decompressedData, nil
	}

	return data, nil
}

func isCompressed(data []byte) bool {
	// Check for the compressed data magic number at the beginning of the data
	magic := constants.CompressedDataMagic
	return len(data) >= len(magic) && string(data[:len(magic)]) == magic
}

func DecompressLZ4(compressedData []byte) ([]byte, error) {
	// Assuming the compressed data has a format that the lz4 package can handle directly
	var out bytes.Buffer
	reader := lz4.NewReader(bytes.NewReader(compressedData))

	_, err := out.ReadFrom(reader)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
