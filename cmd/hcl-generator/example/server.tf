resource "sakuracloud_server" "test_ubuntu_01" {
  name = "test-ubuntu-01"
  tags =[ "test", "ubuntu",]

  core   = 1
  memory = 1

  disks = ["sakuracloud_disk.test_ubuntu_01.id"]

  network_interface {
    upstream = "shared"
  }

  timeouts {
    create = "1h"
    delete = "1h"
  }
}
resource "sakuracloud_server" "test_ubuntu_02" {
  name = "test-ubuntu-02"
  tags =[ "test", "ubuntu",]

  core   = 1
  memory = 1

  disks = ["sakuracloud_disk.test_ubuntu_02.id"]

  network_interface {
    upstream = "shared"
  }

  timeouts {
    create = "1h"
    delete = "1h"
  }
}
