FROM mcr.microsoft.com/dotnet/sdk:5.0-alpine3.13 AS build
RUN apk add --no-cache go nodejs npm build-base
ADD . /src
RUN cd /src && ./utils/install.sh
RUN cd /src && ./utils/build.sh -ldflags="-s -w" -v -tags csharp

FROM mcr.microsoft.com/dotnet/runtime:5.0.3-alpine3.13
WORKDIR /app/cmd/bot
RUN apk add --no-cache icu
ENV LIBCOREFOLDER /usr/share/dotnet/shared/Microsoft.NETCore.App/5.0.3
COPY --from=build /src/web/static /app/web/static
COPY --from=build /src/web/views /app/web/views
COPY --from=build /src/cmd/bot/bot /app/cmd/bot/bot
COPY --from=build /src/migrations /app/migrations/
COPY --from=build /src/cmd/bot/*.dll /app/cmd/bot/
COPY --from=build /src/cmd/bot/charmap.bin.gz /app/cmd/bot/
RUN chmod 777 /app/cmd/bot/charmap.bin.gz
CMD ["./bot"]
