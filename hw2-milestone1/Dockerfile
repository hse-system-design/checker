FROM mongo:4.4

# Install stuff needed
RUN apt-get update
RUN apt-get install software-properties-common -y
RUN add-apt-repository ppa:longsleep/golang-backports
RUN apt-get update
RUN apt-get install golang-go iproute2 -y