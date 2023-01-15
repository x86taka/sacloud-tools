data "sakuracloud_archive" "ubuntu_server_22041_lts_64bit" {
  filter {
    names = ["Ubuntu Server 22.04.1 LTS 64bit"]
  }
}
