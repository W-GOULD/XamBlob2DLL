// pkg/constants/constants.go

package constants

const (
	// AssemblyStore Constants
	AssemblyStoreMagic         string = "XABA"
	AssemblyStoreFormatVersion int    = 1

	CompressedDataMagic string = "XALZ"

	// Assemblies related
	FileAssembliesBlob      string = "assemblies.blob"
	FileAssembliesBlobARM   string = "assemblies.armeabi_v7a.blob"
	FileAssembliesBlobARM64 string = "assemblies.arm64_v8a.blob"
	FileAssembliesBlobX86   string = "assemblies.x86.blob"
	FileAssembliesBlobX8664 string = "assemblies.x86_64.blob"

	FileAssembliesManifest string = "assemblies.manifest"

	// Output / Internal
	FileAssembliesJSON string = "assemblies.json"
)

// GetArchitectureMap returns a map of architectures to their corresponding file names
func GetArchitectureMap() map[string]string {
	return map[string]string{
		"arm":    FileAssembliesBlobARM,
		"arm64":  FileAssembliesBlobARM64,
		"x86":    FileAssembliesBlobX86,
		"x86_64": FileAssembliesBlobX8664,
	}
}
