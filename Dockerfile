FROM golang:1.19 as builder

RUN mkdir /app
WORKDIR /app
COPY [ "go.*", "/app/" ]
COPY ["vendor", "/app/vendor"]
RUN go mod download

COPY . ./
RUN make build

# FROM alpine:3 as tgswitch

# RUN apk update --no-cache && apk add git curl bash
# RUN curl -L https://raw.githubusercontent.com/warrensbox/tgswitch/release/install.sh | bash

FROM alpine:3

RUN apk update --no-cache && apk add git curl bash libc6-compat
RUN curl -L https://raw.githubusercontent.com/warrensbox/tgswitch/release/install.sh | bash
RUN tgswitch 0.45.0
RUN git clone --depth=1 https://github.com/tfutils/tfenv.git $HOME/.tfenv && \
    ln -s $HOME/.tfenv/bin/terraform /usr/local/bin/terraform && \
    ln -s $HOME/.tfenv/bin/tfenv /usr/local/bin/tfenv
RUN tfenv use 1.4.4
COPY --from=builder /app/dist/app ./

ENTRYPOINT [ "/app" ]