
artifact "googlecompute" "ubuntu-1804-lts-vault" {
    source = "artifact.googlecompute.ubuntu-1804-lts"
    // this sets the source of this artifact;

    provisioner "shell" {
        script = [
            "./setup-vault-0.11.3.sh" // relative paths are relative to hcl file
        ]
    }


    artifact "compress" "ubuntu-1804-lts-vault.gz" {
        // same as writting: source = "artifact.googlecompute.ubuntu-1804-lts-vault"

        paths = [
            "gs://mybucket1/vault-0.11.3.gz.tar.gz"
        ]
        keep_input_artifact = true
    }
}
