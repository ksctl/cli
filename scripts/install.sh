#!/usr/bin/env python3

import os
import platform
import subprocess
import urllib.request
import hashlib
import tarfile
import shutil
import json

class Colors:
    Red = '\033[1;31m'
    Green = '\033[1;32m'
    Blue = '\033[1;34m'
    Yellow = '\033[1;33m'
    NoColor = '\033[0m'

def print_colored(message, color):
    print(f"{color}{message}{Colors.NoColor}")

def get_latest_release():
    url = "https://api.github.com/repos/ksctl/cli/releases/latest"
    with urllib.request.urlopen(url) as response:
        data = response.read()
        release_info = json.loads(data)
        return release_info["tag_name"]

def download_file(url, local_filename):
    print_colored(f"Downloading {url} to {local_filename}", Colors.Blue)
    try:
        with urllib.request.urlopen(url) as response:
            with open(local_filename, 'wb') as out_file:
                shutil.copyfileobj(response, out_file)
    except urllib.error.HTTPError as e:
        print_colored(f"HTTP Error: {e.code} - {e.reason}", Colors.Red)
        exit(1)
    except urllib.error.URLError as e:
        print_colored(f"URL Error: {e.reason}", Colors.Red)
        exit(1)

def verify_checksum(file_path, expected_checksum):
    sha256_hash = hashlib.sha256()
    with open(file_path, "rb") as f:
        for byte_block in iter(lambda: f.read(4096), b""):
            sha256_hash.update(byte_block)
    calculated_checksum = sha256_hash.hexdigest()
    return calculated_checksum == expected_checksum

def main():
    print_colored("All necessary dependencies are present.", Colors.Green)

    ksctl_version = os.getenv("KSCTL_VERSION")
    if not ksctl_version:
        print_colored("Fetching latest release version...", Colors.Blue)
        ksctl_version = get_latest_release()
    print_colored(f"Using version: {ksctl_version}", Colors.Green)

    os_name = platform.system()
    arch = platform.machine()

    if arch == "x86_64":
        arch = "amd64"
    elif arch == "arm64":
        arch = "arm64"
    else:
        print_colored(f"Unsupported architecture: {arch}", Colors.Red)
        exit(1)

    if os_name not in ["Linux", "Darwin"]:
        print_colored(f"Unsupported OS: {os_name}", Colors.Red)
        exit(1)

    os_name = os_name.lower()
    print_colored(f"Detected OS: {os_name}, Architecture: {arch}", Colors.Green)

    print_colored("Downloading files...", Colors.Blue)
    download_url_base = f"https://github.com/ksctl/cli/releases/download/{ksctl_version}"
    tar_file = f"ksctl-cli_{ksctl_version[1:]}_{os_name}_{arch}.tar.gz"
    checksum_file = f"ksctl-cli_{ksctl_version[1:]}_checksums.txt"
    download_file(f"{download_url_base}/{tar_file}", tar_file)
    download_file(f"{download_url_base}/{checksum_file}", checksum_file)

    print_colored("Verifying checksum...", Colors.Blue)
    with open(checksum_file, 'r') as f:
        checksums = f.readlines()
    expected_checksum = [line.split()[0] for line in checksums if tar_file in line][0]

    if not verify_checksum(tar_file, expected_checksum):
        print_colored("Checksum verification failed!", Colors.Red)
        exit(1)
    print_colored("Checksum verification passed.", Colors.Green)

    print_colored("Installing ksctl...", Colors.Blue)

    temp_dir = "/tmp/ksctl"
    os.makedirs(temp_dir, exist_ok=True)

    try:
        with tarfile.open(tar_file, "r:gz") as tar:
            tar.extractall(temp_dir)

        ksctl_binary = os.path.join(temp_dir, "ksctl")
        if not os.path.isfile(ksctl_binary):
            print_colored(f"ksctl binary not found in the tarball", Colors.Red)
            exit(1)

        print_colored("Moving ksctl to /usr/local/bin (requires sudo)...", Colors.Blue)
        subprocess.run(["sudo", "mv", "-v", ksctl_binary, "/usr/local/bin/ksctl"])

        if shutil.which("ksctl"):
            print_colored("INSTALL COMPLETE!", Colors.Green)
        else:
            print_colored("Installation failed. Please check the output for errors.", Colors.Red)
            exit(1)
    finally:
        print_colored("Cleaning up temporary files...", Colors.Blue)
        shutil.rmtree(temp_dir)
        os.remove(tar_file)
        os.remove(checksum_file)

if __name__ == "__main__":
    main()
