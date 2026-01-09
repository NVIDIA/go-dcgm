target "default" {
  name = "go-dcgm-${replace(distro, ".", "-")}-${replace(go, ".", "-")}-${replace(cuda, ".", "-")}"
  tags = ["go-dcgm:${distro}-go${go}-cuda${cuda}-dcgm${dcgm}"]
  platforms = ["linux/amd64"]
  matrix = {
    go = ["1.25.5"]
    distro = ["ubuntu24.04", "ubuntu22.04"]
    cuda = ["12.9.1", "13.1.0"]
    dcgm = ["4.5.0-1"]
  }
  args = {
    GO_VERSION = go
    DISTRO_FLAVOR = distro
    CUDA_VERSION = cuda
    DCGM_VERSION = dcgm
  }
}
