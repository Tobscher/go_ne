package core

type Facts map[string]string

func (f Facts) OS() string {
	return f["OS"]
}

func (f Facts) Arch() string {
	arch, ok := f["MACH"]
	if !ok {
		return ""
	}

	if arch == "x86_64" {
		return "amd64"
	}

	if arch == "i686" {
		return "386"
	}

	return arch
}

func (f Facts) Home() string {
	return f["HOME"]
}
