# Building the Test Environment Image with Packer

This guide explains how to build a custom image pre-provisioned with the project's Ansible playbooks using Packer.
This image is configured/used for running GitHub actions,stress tests,benchmark tests etc.
## Prerequisites

*   **Packer**: Installed on your local machine.
*   **QEMU/KVM**: Installed for virtualization.
*   **Ansible**: Installed locally (for `ansible-local` provisioner dependency checks, though it runs inside the VM).

## Configuration

The Packer configuration is located in `packer/fedora.pkr.hcl`.

*   **Base Image**: Defaults to Fedora 43 Cloud Base.
*   **SSH Credentials**: The temporary build user is `fedora` with password `packer`.
*   **Ansible**: The build process uploads the `ansible/` directory and runs `site.yml`.

### Secrets Handling

**Important**: The `ansible/secrets.yml` file is empty,you should fill it up. 
To build the image, you must either:
1.  Provide the vault password (if supported by your workflow).
2.  Temporarily unencrypt `ansible/secrets.yml` ,and encrypt again after the workflow again.

## Building the Image

1. Initialize Packer (download plugins):
    ```bash
    packer init fedora.pkr.hcl
    ```

2. Build the image:
    ```bash
    packer build fedora.pkr.hcl
    ```

    *If the output directory already exists, use `packer build -force fedora.pkr.hcl`or simply remove the previous output.*

## Output

The build artifact will be saved to:
`packer/output-fedora-nocloud/fedora-nocloud.qcow2`

You can import this QCOW2 image into your virtualization platform (Proxmox, OpenStack, KVM/Virt-Manager) or directly use on bare metal with :
   ````bash
    qemu-img  convert fedora-nocloud.qcow2 image.raw
    ``` 