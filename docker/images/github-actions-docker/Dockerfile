# base
FROM ubuntu:24.10

ARG DEBIAN_FRONTEND=noninteractive
ARG RUNNER_VERSION=2.325.0

RUN apt-get update -y && apt-get upgrade -y \
  && useradd -m docker \
  && usermod -aG root docker

RUN  apt-get install -y --no-install-recommends \
    curl jq ca-certificates


RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc \
    && chmod a+r /etc/apt/keyrings/docker.asc \
    && echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
    $(. /etc/os-release && echo "${UBUNTU_CODENAME:-$VERSION_CODENAME}") stable" | \
    tee /etc/apt/sources.list.d/docker.list > /dev/null \
    && apt-get update \
    && apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin


WORKDIR /home/docker

RUN mkdir -p ./actions-runner && cd ./actions-runner \
    && curl -o actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz -L https://github.com/actions/runner/releases/download/v${RUNNER_VERSION}/actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz \
    && tar xzf ./actions-runner-linux-x64-${RUNNER_VERSION}.tar.gz

RUN chown -R docker ~docker \
    && sed -i 's/apt_get_with_fallbacks\(.*libssl.*\)/apt_get_with_fallbacks libssl3t64 \1/g' actions-runner/bin/installdependencies.sh \
    && sed -i 's/apt_get_with_fallbacks\(.*libicu72.*\)/apt_get_with_fallbacks libicu74 \1/g' actions-runner/bin/installdependencies.sh \
    && ./actions-runner/bin/installdependencies.sh


COPY docker-entry.sh docker-entry.sh
COPY start.sh start.sh
RUN chmod +x start.sh docker-entry.sh

# USER docker

CMD ["./docker-entry.sh", "su", "docker", "-c", "./start.sh"]

