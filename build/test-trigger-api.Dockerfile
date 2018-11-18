FROM krostar/ci:build-go.v0.5.0-dev AS builder

COPY . /app
RUN /app-build/build.sh trigger-api

#####
# puppetr10k is a stage that contains puppet and r10k
#####
FROM puppet/puppet-agent-alpine:6.0.4 AS test-trigger-api

RUN apk --no-cache add git~=2.18 && \
    gem install --no-rdoc --no-ri r10k

RUN adduser -h /etc/puppetlabs/r10k/ -s /bin/nologin -S -D r10k && \
    rm -rf /etc/puppetlabs/code/environment/* && \
    chown -R r10k: /etc/puppetlabs/code/environment && \
    chmod 755 /etc/puppetlabs/code/environment

USER r10k
WORKDIR /etc/puppetlabs/r10k/

COPY test/docker/r10k-local.yaml r10k.yaml
COPY test/docker/puppetsources /opt/puppetsources/
COPY configs/trigger-api/docker.yaml api-trigger.yaml
COPY --from=builder /app/build/bin/trigger-api /usr/local/bin/trigger-api
COPY scripts/trigger-api-entrypoint.sh /entrypoint.sh

EXPOSE 8080
ENTRYPOINT ["/entrypoint.sh"]
CMD [ "-config", "./api-trigger.yaml"]
