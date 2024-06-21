package common

var Archs = []string{"x86_64", "aarch64", "i386", "amd64", "arm64", "sw_64", "mips64", "loong64", "loongarch64", "sw_64"}

const (
	TotalSheet    = "总览"                                // excel首页
	WorkDir       = "/tmp/app_check_work_dir/"          // 默认执行工作目录
	BaseLineDir   = "/etc/deepin-app-analyze/baseline"  // 系统基线配置目录
	MkdirError    = "mkdir work directory error"        // 创建目录错误提示
	CopyError     = "copy file to work directory error" // 复制文件错误提示
	FileInfoError = "get file information error"        // 获取app文件信息错误提示
	ArchError     = "Unknown system architecture"       // 未知的系统架构
)

// 软件包信息
type PackageInfo struct {
	Name    string `json:"name"`    // 软件包名称
	Arch    string `json:"arch"`    // 软件包架构
	Version string `json:"version"` // 软件包版本
	Size    int64  `json:"size"`    // 软件包大小
}

// 配置有无变动
type Changed struct {
	Change   []string `json:"change"`   // 有变化分类
	NoChange []string `json:"nochange"` // 无变化分类
}

// 配置CPU架构
type Config struct {
	I386        Changed `json:"i386"`
	AMD64       Changed `json:"amd64"`
	ARM64       Changed `json:"arm64"`
	X86_64      Changed `json:"x86_64"`
	AARCH64     Changed `json:"aarch64"`
	X86         Changed `json:"x86"`
	ARM         Changed `json:"arm"`
	SW          Changed `json:"sw"`
	SW_64       Changed `json:"sw_64"`
	LOONG64     Changed `json:"loong64"`
	LOONGARCH64 Changed `json:"loongarch64"`
	MIPS        Changed `json:"mips"`
	MIPS64      Changed `json:"mips64"`
	MIPSEL      Changed `json:"mipsel"`
	MIPSEL64    Changed `json:"mipsel64"`
	RISCV64     Changed `json:"riscv64"`
	RISCV32     Changed `json:"riscv32"`
}
