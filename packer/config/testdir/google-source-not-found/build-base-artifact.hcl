
artifact "googlecompute" "ubuntu-1804-lts" {

    source = "artifact.foo.bar"

    // default name = "ubuntu-1804-lts"
    // source defaults to "ubuntu-1804-lts"

    provisioner "shell" {
        inline = [
            "my name is ${self.name}",
            "my type is ${self.type}"
        ]
    }

    // this generates an artifact we could use later;
    // this artifact can be based upon to build more images
}
