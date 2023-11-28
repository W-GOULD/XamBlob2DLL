// pkg/manifest/manifest.go

package manifest

import (
    "bufio"
    "os"
    "strings"
)

type ManifestEntry struct {
    Hash32   string
    Hash64   string
    BlobID   int
    BlobIdx  int
    Name     string
}

type ManifestList []ManifestEntry

func (ml ManifestList) GetIdx(blobID, blobIdx int) *ManifestEntry {
    for _, entry := range ml {
        if entry.BlobIdx == blobIdx && entry.BlobID == blobID {
            return &entry
        }
    }
    return nil
}

func ReadManifest(inManifest string) (ManifestList, error) {
    file, err := os.Open(inManifest)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var manifestList ManifestList
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        if line == "" || strings.HasPrefix(line, "Hash") {
            continue
        }

        splitLine := strings.Fields(line)
        // Add error checking for splitLine elements
        manifestList = append(manifestList, ManifestEntry{
            Hash32:  splitLine[0],
            Hash64:  splitLine[1],
            BlobID:  //Convert splitLine[2] to int,
            BlobIdx: //Convert splitLine[3] to int,
            Name:    splitLine[4],
        })
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return manifestList, nil
}
