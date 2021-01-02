FROM public.ecr.aws/a1p3q7r0/alpine:3.12.0 as ffmpeg-builder

RUN cd /usr/local/bin && \
    mkdir ffmpeg && \
    cd ffmpeg/ && \
    wget https://johnvansickle.com/ffmpeg/builds/ffmpeg-git-amd64-static.tar.xz && \
    tar xvf *.tar.xz && \
    rm -f *.tar.xz && \
    mv ffmpeg-git-*-amd64-static/ffmpeg .

FROM public.ecr.aws/bitnami/golang:1.15.5 as function-builder

WORKDIR /build
COPY . .
RUN go build -o media-converter main.go

FROM public.ecr.aws/lambda/go:latest
COPY --from=ffmpeg-builder /usr/local/bin/ffmpeg/ffmpeg /usr/local/bin/ffmpeg/ffmpeg
RUN ln -s /usr/local/bin/ffmpeg/ffmpeg /usr/bin/ffmpeg

COPY --from=function-builder /build/media-converter /var/task/media-converter

CMD [ "media-converter" ]