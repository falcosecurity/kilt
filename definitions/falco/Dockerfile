FROM falcosecurity/falco:latest-slim AS falco
FROM docker.pkg.github.com/falcosecurity/kilt/utilities:latest AS kilt-utils
FROM alpine:latest AS pdig-build

# Build Pdig
RUN apk add g++ gcc cmake cmake make libtool elfutils-dev libelf-static linux-headers git
RUN mkdir /source
RUN git clone https://github.com/falcosecurity/pdig /source/pdig
RUN git clone https://github.com/draios/sysdig /source/sysdig
RUN mkdir /source/pdig/build
RUN cd /source/pdig/build && cmake -DMUSL_OPTIMIZED_BUILD=True ..
RUN cd /source/pdig/build && make

FROM scratch
COPY --from=falco /usr/bin/falco /vendor/falco/bin/falco
COPY --from=falco /usr/share/falco/lua /vendor/falco/share/lua
COPY --from=falco /etc /vendor/falco/etc
COPY --from=pdig-build /source/pdig/build/pdig /vendor/falco/bin/pdig
COPY --from=kilt-utils /kilt/waitforever /vendor/falco/bin/waitforever
ADD falco.yaml /vendor/falco/falco.yaml

VOLUME ["/vendor/falco"]

CMD ["/vendor/falco/bin/waitforever"]