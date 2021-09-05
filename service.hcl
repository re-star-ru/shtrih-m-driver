job "kkt-server" {
  datacenters = [
    "restar"
  ]
  type = "service"

  group "default" {
    network {
      port "kkt_server" {
        static = 19300
      }
    }

    task "kkt_server" {
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
        ADDR = "${NOMAD_ADDR_kkt_server}"
      }

    }
  }
}


