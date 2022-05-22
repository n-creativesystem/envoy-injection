FROM envoyproxy/envoy:v1.19-latest as envoy

FROM ncreativesystem/golang-developments:dev

RUN mkdir -p /etc/envoy \
  && curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl" \
  && chmod +x kubectl \
  && mv ./kubectl /usr/local/bin/kubectl \
  && curl "https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3" | bash \
  && apt-get update \
  && apt-get install -y \
  sudo \
  xdg-utils
COPY --from=envoy /usr/local/bin/envoy /usr/local/bin/envoy
COPY --from=envoy /etc/envoy/envoy.yaml /etc/envoy/envoy.yaml

ENV USERNAME=user \
  GROUPNAME=user
ARG UID=1001 \
  GID=1002 \
  PASSWORD=user

RUN groupadd -g $GID $GROUPNAME \
  && useradd -m -s /bin/bash -u $UID -g $GID -G sudo $USERNAME \
  && echo $USERNAME:$PASSWORD | chpasswd \
  && echo "$USERNAME   ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers \
  && chown -R "${USERNAME}:${GROUPNAME}" /go

USER $USERNAME
ENV HOME="/home/${USERNAME}"
WORKDIR /home/$USERNAME/

CMD [ "bash" ]
