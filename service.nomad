job "kkt-server" {
  datacenters = ["dc1"]
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
        cpu    = 200
        memory = 128
      }

      config {
        image        = "ghcr.io/${REPO}:${TAG}"
        network_mode = "host"
      }

      // config {
      //   image        = "ghcr.io/[[.repo]]:[[.tag]]"
      //   network_mode = "host"
      // }

      env {
        ADDR = "${NOMAD_ADDR_kkt_server}"
      }

    }
  }
}


