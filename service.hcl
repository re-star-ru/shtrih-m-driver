job "shtrih-m-driver" {
  datacenters = [
    "restar"
  ]

  type = "service"

  group "default" {
    network {
      port "feziv" {
        host_network = "private"
      }
    }

    task "shtrih-m-driver" {
      //      service {
      //        port = "feziv"
      //        tags = [
      //          "reproxy.enabled=1",
      //          "reproxy.server=feziv.com,www.feziv.com"
      //        ]
      //      }
      // serve static files for feziv.com

      resources {
        memory = 64
      }

      driver = "docker"

      env {
        LISTEN = "${NOMAD_ADDR_feziv}"
      }

      config {
        image = "ghcr.io/fess932/shtrih-m-driver:latest"
        network_mode = "host"
      }

    }
  }
}


