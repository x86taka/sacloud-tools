resource "sakuracloud_disk" "test_ubuntu_01" {
  name = "test-ubuntu-01"
  tags =[]

  source_archive_id = data.sakuracloud_archive.ubuntu_server_22041_lts_64bit.id
  size              = 20

  timeouts {
    create = "1h"
    delete = "1h"
  }
}
resource "sakuracloud_disk" "test_ubuntu_02" {
  name = "test-ubuntu-02"
  tags =[]

  source_archive_id = data.sakuracloud_archive.ubuntu_server_22041_lts_64bit.id
  size              = 20

  timeouts {
    create = "1h"
    delete = "1h"
  }
}
