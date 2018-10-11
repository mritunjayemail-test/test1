
artifact "googlecompute" "ubuntu-1804-lts-consul" {
    source = "artifact.googlecompute.ubuntu-1804-lts"
    // this sets the source of this artifact;

    provisioner "shell" {
        script = [
            "./setup-consul-1.2.3.sh" // relative paths are relative to hcl file
        ]
    }
}


artifact "compress" "ubuntu-1804-lts-consul.gz" {
    source = "artifact.googlecompute.ubuntu-1804-lts-consul"

}
