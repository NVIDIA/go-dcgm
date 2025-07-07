target "default" {
  name = "go-dcgm-${replace(distro, ".", "-")}-${replace(go, ".", "-")}-${replace(cuda, ".", "-")}"
  tags = ["go-dcgm:${distro}-go${go}-cuda${cuda}-dcgm${dcgm}"]
  platforms = ["linux/amd64"]
  matrix = {
    go = ["1.24.4"]
    distro = ["ubuntu24.04", "ubuntu22.04", "ubuntu20.04"]
    cuda = ["12.9.1", "12.5.1"]
    dcgm = ["4.2.3-2"]
  }
  args = {
    GO_VERSION = go
    DISTRO_FLAVOR = distro
    CUDA_VERSION = cuda
    DCGM_VERSION = dcgm
  }
}
