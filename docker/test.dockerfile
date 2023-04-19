FROM steebchen/go-prisma:go_v1.13-prisma_v1.34.10 as builder

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

# required since renaming https://www.prisma.io/blog/prisma-2-beta-b7bcl0gd8d8e#renaming-the-prisma2-cli
RUN ln -s /usr/bin/prisma /usr/bin/prisma1

RUN PRISMA_HOST=x APP_NAMESPACE=x go generate

ENV TZ=UTC

RUN go test ./...
