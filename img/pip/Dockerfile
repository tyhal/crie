FROM python:3-alpine3.22

# [ Pips ]
ENV CRYPTOGRAPHY_DONT_BUILD_RUST=1
COPY requirements.txt /requirements.txt
ENV BUILD_LIBS="python3-dev build-base libffi-dev libressl-dev"
RUN apk add --no-cache python3 py3-pip $BUILD_LIBS \
    && pip3 --no-cache-dir install -r requirements.txt \
    && apk del --no-cache $BUILD_LIBS

# hadolint ignore=DL3059
RUN mkdir /.ansible \
    && chmod -R o+rw /home /.ansible

# hadolint ignore=DL3059
RUN adduser -D standards
USER standards

ENV XDG_CONFIG_DIR=/etc/.config

ENV YAMLLINT_CONFIG_FILE=$XDG_CONFIG_DIR/yamllint/config
COPY --chown=standards:standards conf/yaml/.yamllint.yaml $YAMLLINT_CONFIG_FILE
ENV CMAKELINT_CONFIG_FILE=$XDG_CONFIG_DIR/cmakelintrc
COPY --chown=standards:standards conf/cmake/.cmakelintrc $CMAKELINT_CONFIG_FILE