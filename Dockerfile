FROM mcr.microsoft.com/dotnet/sdk:10.0-alpine-aot AS build
RUN apk add --no-cache go nodejs npm build-base rustup
RUN rustup-init -y
ENV PATH="/root/.cargo/bin:${PATH}"
ADD . /src
RUN cd /src && ./utils/install.sh linux-musl-x64
RUN cd /src && ./utils/build.sh -v -tags csharp

FROM alpine:latest
WORKDIR /app/cmd/bot
COPY --from=build /src/web/static /app/web/static
COPY --from=build /src/web/views /app/web/views
COPY --from=build /src/cmd/bot/bot /app/cmd/bot/bot
COPY --from=build /src/migrations /app/migrations/
COPY --from=build /src/cmd/bot/MessageHeightTwitch /app/cmd/bot/MessageHeightTwitch
CMD ["./bot"]
