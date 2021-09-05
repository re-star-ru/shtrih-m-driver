job "shtrih-m-driver" {
  datacenters = [
    "restar"
  ]
  type = "service"

  group "default" {
    network {
      port "shtrih-m" {
        static = 19300
      }
    }

    task "shtrih-m-driver" {
      driver = "docker"

      resources {
        cpu = 100
        memory = 64
      }

      config {
        image = "ghcr.io/[[.repo]]:[[.tag]]"
        network_mode = "host"
      }

      env {
        ADDR = "${NOMAD_ADDR_feziv}"
      }

    }
  }
}


