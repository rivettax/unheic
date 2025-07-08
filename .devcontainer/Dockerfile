FROM mcr.microsoft.com/devcontainers/go:1-1.24-bookworm

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive

# Ensure pipx path is in the main PATH (usually handled by pipx install, but explicitly adding doesn't hurt)
# RUN pipx ensurepath
# Add pipx bin dir for vscode user to PATH
ENV PATH="/home/vscode/.local/bin:${PATH}" 