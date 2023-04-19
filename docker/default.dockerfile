ARG GO=v1.13
ARG PRISMA=v1.34.10
ARG TAG=go_${GO}-prisma_${PRISMA}
ARG IMAGE=steebchen/go-prisma:go_${GO}-prisma_${PRISMA}

FROM $IMAGE as builder

ENV GO111MODULE=on

COPY go.mod go.sum ./

RUN go mod download

COPY . ./

# required since renaming https://www.prisma.io/blog/prisma-2-beta-b7bcl0gd8d8e#renaming-the-prisma2-cli
RUN ln -s /usr/bin/prisma /usr/bin/prisma1

RUN PRISMA_HOST=x APP_NAMESPACE=x go generate

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o /main .

# the prisma cli in the image is required to run prisma deploy in kubernetes
FROM $IMAGE

COPY --from=builder /main /main

# copy files for prisma deploy
COPY --from=builder /app /app

ENV TZ=UTC

CMD ["/main"]
