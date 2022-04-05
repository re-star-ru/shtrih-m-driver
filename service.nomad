job "kkt-server" {
  datacenters = ["dc1"]
  type = "service"

  group "default" {
    network {
      port "kkt_server" {
        static = 19300
      }
    }

    service {
      name = "kkt"
      port = "kkt_server"

      # The "check" stanza instructs Nomad to create a Consul health check for
      # this service. A sample check is provided here for your convenience;
      # uncomment it to enable it. The "check" stanza is documented in the
      # "service" stanza documentation.

      # check {
      #   name     = "alive"
      #   type     = "tcp"
      #   interval = "10s"
      #   timeout  = "2s"
      # }
    }

    task "kkt_server" {
      driver = "docker"

      resources {
        cpu    = 200
        memory = 128
      }

      config {
        image        = "ghcr.io/${IMAGE_NAME}:${TAG}"
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


