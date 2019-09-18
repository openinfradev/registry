# Taco-Registry Builder
FROM scratch
LABEL maintainer="linus lee <linus@exntu.com>"

WORKDIR /
COPY ./builder .

EXPOSE 4000

ENTRYPOINT ["./builder"]