package pe_details

type PEOffset int64
type PEFile struct {
	ExecutableType ExecutableType
	HasStub bool
	Offset PEOffset

	CoffHeader COFFHeader
	OptionalHeader OptionalHeader
	SectionTables []SectionTable

	CoffRelocations []CoffRelocation
	CoffSymbolTable CoffSymbolTable
}

type ExecutableType int
const (
	Unknown ExecutableType = iota + 1
	Image
	Object
)

type COFFHeader struct {
	Machine uint16
	NumberOfSections int16
	TimeStamp int32
	SymbolTablePtr int32
	NumberOfSymbols int32
	OptionalHeaderSz int16
	Characteristics int16
}

type OptionalHeader struct {
	Is64Bit bool
	StandardFields OHStandardFields
	WindowsFields OHWindowsFieldsPrintable
	DataDirectories OHDataDirectories
}

type OHStandardFields struct {
	Magic int16
	MajorLinkerVersion byte
	MinorLinkerVersion byte
	SizeOfCode int32
	SizeOfInitializedData int32
	SizeOfUninitializedData int32
	AddressOfEntryPoint int32
	BaseOfCode int32
	BaseOfData int32
}

type OHWindowsFieldsPrintable interface {
	Print(tabCont int) (ret string)
}

type OHWindowsFields32 struct {
	ImageBase int32
	SectionAlignment int32
	FileAlignment int32
	MajorOSVersion int16
	MinorOSVersion int16
	MajorImageVersion int16
	MinorImageVersion int16
	MajorSubsystemVersion int16
	MinorSubsystemVersion int16
	Win32SessionVal int32
	SizeOfImage int32
	SizeOfHeaders int32
	CheckSum int32
	Subsystem int16
	DllCharacteristics uint16
	SizeOfStackReserve int32
	SizeOfStackCommit int32
	SizeOfHeapReserve int32
	SizeOfHeapCommit int32
	LoaderFlags int32
	NumberOfRvaAndSizes int32
}

type OHWindowsFields64 struct {
	ImageBase int64
	SectionAlignment int32
	FileAlignment int32
	MajorOSVersion int16
	MinorOSVersion int16
	MajorImageVersion int16
	MinorImageVersion int16
	MajorSubsystemVersion int16
	MinorSubsystemVersion int16
	Win32SessionVal int32
	SizeOfImage int32
	SizeOfHeaders int32
	CheckSum int32
	Subsystem int16
	DllCharacteristics int16
	SizeOfStackReserve int64
	SizeOfStackCommit int64
	SizeOfHeapReserve int64
	SizeOfHeapCommit int64
	LoaderFlags int32
	NumberOfRvaAndSizes int32
}

type OHDataDirectories struct {
	ImageDataDirectories []ImageDataDirectory
}

type ImageDataDirectory struct {
	VirtualAddress int32
	Size int32
}

type SectionTable struct {
	Name                 [8]byte
	VirtualSize          int32
	VirtualAddress       int32
	SizeOfRawData        int32
	PointerToRawData     int32
	PointerToRelocations int32
	PointerToLineNumbers int32
	NumberOfRelocations  int16
	NumberOfLineNumbers  int16
	Characteristics      AvailableCharacteristics
}

type AvailableCharacteristics uint32
const (
	IMAGE_SCN_TYPE_NO_PAD AvailableCharacteristics = 0x00000008
	IMAGE_SCN_CNT_CODE AvailableCharacteristics = 0x00000020
	IMAGE_SCN_CNT_INITIALIZED_DATA AvailableCharacteristics = 0x00000040
	IMAGE_SCN_CNT_UNINITIALIZED_DATA AvailableCharacteristics = 0x00000080
	IMAGE_SCN_LNK_OTHER AvailableCharacteristics = 0x00000100
	IMAGE_SCN_LNK_INFO AvailableCharacteristics = 0x00000200
	IMAGE_SCN_LNK_REMOVE AvailableCharacteristics = 0x00000800
	IMAGE_SCN_LNK_COMDAT AvailableCharacteristics = 0x00001000
	IMAGE_SCN_GPREL AvailableCharacteristics = 0x00008000
	IMAGE_SCN_MEM_PURGEABLE AvailableCharacteristics = 0x00020000
	IMAGE_SCN_MEM_16BIT AvailableCharacteristics = 0x00020000
	IMAGE_SCN_MEM_LOCKED AvailableCharacteristics = 0x00040000
	IMAGE_SCN_MEM_PRELOAD AvailableCharacteristics = 0x00080000
	IMAGE_SCN_ALIGN_1BYTES AvailableCharacteristics = 0x00100000
	IMAGE_SCN_ALIGN_2BYTES AvailableCharacteristics = 0x00200000
	IMAGE_SCN_ALIGN_4BYTES AvailableCharacteristics = 0x00300000
	IMAGE_SCN_ALIGN_8BYTES AvailableCharacteristics = 0x00400000
	IMAGE_SCN_ALIGN_16BYTES AvailableCharacteristics = 0x00500000
	IMAGE_SCN_ALIGN_32BYTES AvailableCharacteristics = 0x00600000
	IMAGE_SCN_ALIGN_64BYTES AvailableCharacteristics = 0x00700000
	IMAGE_SCN_ALIGN_128BYTES AvailableCharacteristics = 0x00800000
	IMAGE_SCN_ALIGN_256BYTES AvailableCharacteristics = 0x00900000
	IMAGE_SCN_ALIGN_512BYTES AvailableCharacteristics = 0x00A00000
	IMAGE_SCN_ALIGN_1024BYTES AvailableCharacteristics = 0x00B00000
	IMAGE_SCN_ALIGN_2048BYTES AvailableCharacteristics = 0x00C00000
	IMAGE_SCN_ALIGN_4096BYTES AvailableCharacteristics = 0x00D00000
	IMAGE_SCN_ALIGN_8192BYTES AvailableCharacteristics = 0x00E00000
	IMAGE_SCN_LNK_NRELOC_OVFL AvailableCharacteristics = 0x01000000
	IMAGE_SCN_MEM_DISCARDABLE AvailableCharacteristics = 0x02000000
	IMAGE_SCN_MEM_NOT_CACHED AvailableCharacteristics = 0x04000000
	IMAGE_SCN_MEM_NOT_PAGED AvailableCharacteristics = 0x08000000
	IMAGE_SCN_MEM_SHARED AvailableCharacteristics = 0x10000000
	IMAGE_SCN_MEM_EXECUTE AvailableCharacteristics = 0x20000000
	IMAGE_SCN_MEM_READ AvailableCharacteristics = 0x40000000
	IMAGE_SCN_MEM_WRITE AvailableCharacteristics = 0x80000000
)

type CoffRelocation struct {
	SectionIdx int32
	VirtualAddress int32
	SymbolTableIndex int32
	Type int16
}

type CoffSymbolTable struct {
	Symbols []CoffSymbol
}

type CoffSymbol struct {
	Name [8]byte
	Value int32
	SectionNumber int16
	Type int16
	StorageClass byte
	NumberOfAuxSymbols byte
	AuxSymbols []AuxSymbol
}

type AuxSymbol struct {
	TagIndex int32
	TotalSize int32
	PointerToLinenumber int32
	PointerToNextFunction int32
	Unused int16
}