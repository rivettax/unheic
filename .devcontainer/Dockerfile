FROM mcr.microsoft.com/devcontainers/go:2-1.26-trixie

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends \
        libheif-dev \
        pkg-config \
    && apt-get autoremove -y && apt-get clean -y && rm -rf /var/lib/apt/lists/*

# Ensure pipx path is in the main PATH (usually handled by pipx install, but explicitly adding doesn't hurt)
# RUN pipx ensurepath
# Add pipx bin dir for vscode user to PATH
ENV PATH="/home/vscode/.local/bin:${PATH}" 