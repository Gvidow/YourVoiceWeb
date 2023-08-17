FROM golang:alpine3.18 AS build_stage

COPY . /go/src/YourVoiceWeb

WORKDIR /go/src/YourVoiceWeb

RUN go install cmd/main.go cmd/config.go
RUN mv /go/bin/main /go/bin/YourVoiceWeb

RUN go install cmd/buildYAMLConfig/build.go


FROM alpine AS run_stage

WORKDIR /app_binary
COPY --from=build_stage /go/bin/YourVoiceWeb /app_binary/
RUN chmod +x ./YourVoiceWeb
COPY --from=build_stage /go/bin/build /app_binary/
RUN chmod +x ./build
COPY --from=build_stage /go/src/YourVoiceWeb/static /app_binary/static
COPY --from=build_stage /go/src/YourVoiceWeb/index.html /app_binary/

ENV TOKEN_GPT your_token
ENV YANDEX_CLOUD_OAuth your_OAuth
ENV YANDEX_CLOUD_FOLDER_ID your_folder_id
ENV PORT 8080
ENV HOST ""


EXPOSE 8080/tcp


CMD ./build && ./YourVoiceWeb
