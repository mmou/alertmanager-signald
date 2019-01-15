FROM golang:1.11 AS build

WORKDIR /alertmanager-signald
COPY . ./
RUN GOARCH=arm GOARM=5 go build

FROM registry.cutelab.house/signald:arm

COPY --from=build /alertmanager-signald/alertmanager-signald /bin/alertmanager-signald
EXPOSE 8080
ENV SIGNALD_BIND_ADDR 0.0.0.0:8888
CMD ["/bin/alertmanager-signald"]
