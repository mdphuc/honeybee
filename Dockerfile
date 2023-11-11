FROM fedora:latest

RUN dnf update -y 2>/dev/null
RUN dnf upgrade -y 2>/dev/null

RUN rpm --import https://packages.microsoft.com/keys/microsoft.asc 2>/dev/null
RUN dnf install curl -y 2>/dev/null
RUN curl https://packages.microsoft.com/config/rhel/7/prod.repo | tee /etc/yum.repos.d/microsoft.repo 2>/dev/null
RUN dnf makecache 2>/dev/null
RUN dnf install powershell -y 2>/dev/null
