FROM node:22-alpine

RUN adduser -D standards

# hadolint ignore=DL3059
RUN apk add --no-cache npm && npm install -g jsonlint2 remark-cli remark-preset-lint-recommended standard

# hadolint ignore=DL3059
RUN mkdir /.standard-v14-cache \
    && chmod -R o+rw /home /.standard-v14-cache