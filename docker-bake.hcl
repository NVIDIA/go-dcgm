target "default" {
  name = "go-dcgm-${replace(distro, ".", "-")}-${replace(go, ".", "-")}-${replace(cuda, ".", "-")}"
  tags = ["go-dcgm:${distro}-go${go}-cuda${cuda}-dcgm${dcgm}"]
  platforms = ["linux/amd64"]
  matrix = {
    go = ["1.26.7"]
    distro = ["ubuntu26.04", "ubuntu24.04", "ubuntu22.04"]
    cuda = ["13.3.1"]
    dcgm = ["4.6.1-1"]
  }
  args = {
    GO_VERSION = go
    DISTRO_FLAVOR = distro
    CUDA_VERSION = cuda
    DCGM_VERSION = dcgm
  }
}
