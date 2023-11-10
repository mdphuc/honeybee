RUN apt install h
FROM ubuntu:latest

RUN apt update 2>/dev/null
RUN apt upgrade 2>/dev/null
