data "cdo_ftd_device" "ftd" {
  name = var.ftd_name
}

resource "cdo_ftd_device_version" "ftd" {
  ftd_uid = data.cdo_ftd_device.ftd.id
  software_version   = "7.2.5-208"
}