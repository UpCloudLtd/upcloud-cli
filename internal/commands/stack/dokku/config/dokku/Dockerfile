FROM dokku/dokku:latest

COPY bootstrap/* /etc/my_init.d/

RUN echo 'export DOKKU_APP_PROXY_TYPE="nginx-ingress"' >> /skel/home/dokku/.bashrc
RUN echo 'export DOCKER_HOST="tcp://localhost:2375"' >> /skel/home/dokku/.bashrc
RUN echo 'DOKKU_APP_PROXY_TYPE="nginx-ingress"' >> /skel/home/dokku/ENV
RUN echo 'DOCKER_HOST="tcp://localhost:2375"' >> /skel/home/dokku/ENV
RUN echo 'DOCKER_HOST="tcp://localhost:2375"' >> /etc/environment
RUN echo 'DOKKU_APP_PROXY_TYPE="nginx-ingress"' >> /etc/environment
RUN mkdir -p /skel/var/lib/dokku/config/scheduler/--global && echo -n "kubernetes" > /skel/var/lib/dokku/config/scheduler/--global/selected
RUN mkdir -p /skel/var/lib/dokku/config/registry/--global && echo -n "true" > /skel/var/lib/dokku/config/registry/--global/push-on-release
RUN chown dokku:dokku /skel/home/dokku/ENV
RUN chown dokku:dokku /skel/home/dokku/.bashrc
RUN chown dokku:dokku -R /skel/var/lib/dokku/config

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
RUN sudo install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl