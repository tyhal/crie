FROM node:24-alpine3.22

RUN adduser -D standards

# hadolint ignore=DL3059
RUN apk add --no-cache npm && npm install -g jsonlint2 remark-cli remark-preset-lint-recommended standard

# hadolint ignore=DL3059
RUN mkdir /.standard-cache \
    && chmod -R o+rw /home /.standard-cache