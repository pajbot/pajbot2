FROM mcr.microsoft.com/dotnet/sdk:6.0.301-1-alpine3.16 AS build
RUN apk add --no-cache go nodejs npm build-base
ADD . /src
RUN cd /src && ./utils/install.sh
RUN cd /src && ./utils/build.sh -v -tags csharp

FROM mcr.microsoft.com/dotnet/runtime:6.0.6-alpine3.16
WORKDIR /app/cmd/bot
RUN apk add --no-cache icu
ENV LIBCOREFOLDER /usr/share/dotnet/shared/Microsoft.NETCore.App/6.0.6
COPY --from=build /src/web/static /app/web/static
COPY --from=build /src/web/views /app/web/views
COPY --from=build /src/cmd/bot/bot /app/cmd/bot/bot
COPY --from=build /src/migrations /app/migrations/
COPY --from=build /src/cmd/bot/*.dll /app/cmd/bot/
COPY --from=build /src/cmd/bot/charmap.bin.gz /app/cmd/bot/
RUN chmod 777 /app/cmd/bot/charmap.bin.gz
CMD ["./bot"]
