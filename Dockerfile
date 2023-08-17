FROM golang:1.16-alpine

COPY . /go/src/YourVoice

WORKDIR /go/src/YourVoice

RUN go install cmd/main.go cmd/config.go
RUN mv /go/bin/main /go/bin/YourVoice

RUN go install cmd/buildYAMLConfig/build.go

ENV TOKEN_GPT your_token
ENV YANDEX_CLOUD_OAuth your_OAuth
ENV YANDEX_CLOUD_FOLDER_ID your_folder_id
ENV PORT 8080
ENV HOST ""


EXPOSE 8080/tcp

CMD [ "build", "YourVoice" ]